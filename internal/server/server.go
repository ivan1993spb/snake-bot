package server

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
	"github.com/ivan1993spb/snake-bot/internal/server/handlers"
	"github.com/ivan1993spb/snake-bot/internal/server/middlewares"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type Core interface {
	ApplyState(state map[int]int) (map[int]int, error)
	SetupOne(gameId, botsNumber int) (map[int]int, error)
	ReadState() map[int]int
}

type Secure interface {
	VerifyToken(token string) bool
}

type Server struct {
	ctx context.Context

	server *http.Server

	core Core
	sec  Secure
}

func NewServer(ctx context.Context, cfg config.Server, appInfo string,
	core Core, sec Secure) *Server {
	s := &Server{
		server: &http.Server{
			Addr: cfg.Address,
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
		},

		core: core,
		sec:  sec,
	}

	s.initRoutes(ctx, cfg.Debug, cfg.ForbidCORS, appInfo)

	return s
}

const requestPostBotsThrottleLimit = 1

func (s *Server) initRoutes(ctx context.Context, debug, forbidCORS bool,
	appInfo string) {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middlewares.NewRequestLogger(utils.Log(ctx)))
	r.Use(middleware.SetHeader("Server", appInfo))
	r.Use(middleware.GetHead)
	// By default origins (domain, scheme, or port) don't matter.
	if !forbidCORS {
		r.Use(cors.AllowAll().Handler)
	}

	r.Get("/", handlers.WelcomeHandler)
	r.With(middleware.NoCache).Get("/openapi.yaml", handlers.OpenAPIHandler)
	r.Route("/api/bots", func(r chi.Router) {
		r.Use(middlewares.TokenAuth(s.sec))
		r.With(
			middleware.AllowContentType(
				"application/x-www-form-urlencoded",
				"application/json",
				"text/yaml",
			),
			middleware.Throttle(requestPostBotsThrottleLimit),
		).Post("/", handlers.NewApplyStateHandler(s.core))
		r.Get("/", handlers.NewReadStateHandler(s.core))
	})

	if debug {
		r.Mount("/debug", middleware.Profiler())
	}
	r.Handle("/metrics", promhttp.Handler())

	s.server.Handler = r
}

const serverShutdownTimeout = time.Second

const fieldShutdownTimeout = "shutdown_timeout"

func (s *Server) ListenAndServe(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		log := utils.Log(ctx).WithField(fieldShutdownTimeout, serverShutdownTimeout)
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
