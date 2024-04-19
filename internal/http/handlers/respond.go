package handlers

import (
	"encoding/json"
	"mime"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

const (
	acceptTypeJson = "application/json"
	acceptTypeYaml = "text/yaml"
)

func responseType(r *http.Request) string {
	acceptType := r.Header.Get("Accept")

	// If accept type is not specified, then we derive it from
	// content type
	if acceptType == "" || acceptType == "*/*" {
		contentType := r.Header.Get("Content-type")
		mediaType, _, _ := mime.ParseMediaType(contentType)

		if mediaType == mediaTypeJson {
			acceptType = acceptTypeJson
		} else {
			// Default response type is YAML
			acceptType = acceptTypeYaml
		}
	}

	return acceptType
}

type errorResponse struct {
	Code int    `json:"code" yaml:"code"`
	Text string `json:"text" yaml:"text"`
}

func respondError(w http.ResponseWriter, r *http.Request, status int) {
	respond(w, r, status, &errorResponse{
		Code: status,
		Text: http.StatusText(status),
	})
}

func respond(w http.ResponseWriter, r *http.Request, status int, data any) {
	ctx := r.Context()
	log := utils.GetLogger(ctx)

	acceptType := responseType(r)

	log = log.WithField("accept_type", acceptType)
	ctx = utils.WithLogger(ctx, log)
	r = r.WithContext(ctx)

	if acceptType == acceptTypeYaml {
		respondYaml(w, r, status, data)
	} else if acceptType == acceptTypeJson {
		respondJson(w, r, status, data)
	} else {
		log.Error("invalid accept type")

		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func respondJson(w http.ResponseWriter, r *http.Request, status int, data any) {
	ctx := r.Context()
	log := utils.GetLogger(ctx)

	w.Header().Set("Content-Type", acceptTypeJson)
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.WithError(err).Error("write json response")
	}
}

func respondYaml(w http.ResponseWriter, r *http.Request, status int, data any) {
	ctx := r.Context()
	log := utils.GetLogger(ctx)

	w.Header().Set("Content-Type", acceptTypeYaml)
	w.WriteHeader(status)

	enc := yaml.NewEncoder(w)
	defer enc.Close()

	err := enc.Encode(data)
	if err != nil {
		log.WithError(err).Error("write yaml response")
	}
}
