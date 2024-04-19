package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ivan1993spb/snake-bot/internal/config"
	"github.com/ivan1993spb/snake-bot/internal/http/handlers"
	"github.com/ivan1993spb/snake-bot/internal/http/middlewares"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type Core interface {
	handlers.AppGetState
	handlers.AppSetState
}

type Secure interface {
	middlewares.Secure
}

type ServerParams struct {
	Config  config.Server
	AppInfo string
	Core    Core
	Secure  Secure
}

type Server struct {
	server *http.Server
	params ServerParams
}

func NewServer(params ServerParams) *Server {
	s := &Server{
		server: &http.Server{
			Addr: params.Config.Address,
		},
		params: params,
	}

	s.server.Handler = s.initRoutes()

	return s
}

const requestPostBotsThrottleLimit = 1

func (s *Server) initRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middlewares.RequestID)
	r.Use(middlewares.NewRequestLogger())
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Server", s.params.AppInfo))
	r.Use(middleware.GetHead)

	// By default origins (domain, scheme, or port) don't matter.
	if !s.params.Config.ForbidCORS {
		r.Use(cors.AllowAll().Handler)
	}

	r.Get("/", handlers.WelcomeHandler)
	r.With(middleware.NoCache).Get("/openapi.yaml", handlers.OpenAPIHandler)

	r.Route("/api/bots", func(r chi.Router) {
		r.Use(middlewares.JwtTokenAuth(s.params.Secure))
		r.With(
			middleware.AllowContentType(
				"application/x-www-form-urlencoded",
				"application/json",
				"text/yaml",
			),
			middleware.Throttle(requestPostBotsThrottleLimit),
		).Method("POST", "/", handlers.NewSetStateHandler(s.params.Core))
		r.Method("GET", "/", handlers.NewGetStateHandler(s.params.Core))
	})

	if s.params.Config.Debug {
		r.Mount("/debug", middleware.Profiler())
	}

	r.Handle("/metrics", promhttp.Handler())

	return r
}

const serverShutdownTimeout = time.Second

const fieldShutdownTimeout = "shutdown_timeout"

func (s *Server) ListenAndServe(ctx context.Context) error {
	log := utils.GetLogger(ctx)
	log.WithField("address", s.server.Addr).Info("starting server")

	s.server.BaseContext = func(net.Listener) context.Context {
		return utils.WithModule(ctx, "handler")
	}

	go func() {
		<-ctx.Done()

		log := log.WithField(fieldShutdownTimeout, serverShutdownTimeout)
		log.Info("shutting down")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			serverShutdownTimeout,
		)

		defer cancel()

		go func() {
			<-shutdownCtx.Done()

			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Error("graceful shutdown timed out")
				log.Fatal("forcing exit")
			}
		}()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.WithError(err).Error("server shutdown fail")
		}
	}()

	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return errors.Wrap(err, "listen and serve")
}
