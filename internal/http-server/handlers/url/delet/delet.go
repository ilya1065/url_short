package delet

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	resp "url_shortner/internal/lib/api/response"
	"url_shortner/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

type Response struct {
	Response resp.Response
	Alias    string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.DeleteURL"
		log.With(
			slog.String("op", op))
		alias := chi.URLParam(r, "alias")
		err := urlDeleter.DeleteURL(alias)
		var req string
		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		err = urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delet url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delet url"))
			return
		}
		responseOK(w, r, alias)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
