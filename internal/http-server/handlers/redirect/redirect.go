package redirect

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"url_shortener/internal/lib/api/response"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	URL string `json:"alias,omitempty"`
}

type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func New(ctx context.Context, log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"
		log := log.With(
			slog.String("op", op),
		)
		alias := chi.URLParam(r, "alias")
		fmt.Println(alias)
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(ctx, alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, response.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		render.JSON(w, r, Response{Response: response.OK(),
			URL: resURL})

		// http.Redirect(w, r, resURL, http.StatusFound)
	}
}
