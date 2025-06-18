package delete

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"url_shortener/internal/lib/api/response"
	"url_shortener/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(ctx context.Context, alias string) error
}

func New(ctx context.Context, log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log := log.With(
			slog.String("operation: ", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		alias := chi.URLParam(r, "alias")
		fmt.Println(alias)
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}
		err := urlDeleter.DeleteURL(ctx, alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, response.Error("not found"))

			return
		}

		render.JSON(w, r, response.OK())

	}
}
