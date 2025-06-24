package redirect

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	URL string `json:"alias,omitempty"`
}

// интерфейс по месту использования
// сигнатура метода интерфейса должна дублировать сигнатуру метода из storage
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

		//просто возвращает ссылку
		// render.JSON(w, r, Response{Response: response.OK(),
		// 	URL: resURL})

		//в postman не работает)), только через браузер, тк перенеправляет по ссылке
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
