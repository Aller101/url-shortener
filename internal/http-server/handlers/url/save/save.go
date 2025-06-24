package save

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLenght = 6

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=URLSaver
type URLSaver interface {
	SaveURL(ctx context.Context, urlToSave, alias string) (int64, error)
}

func New(ctx context.Context, log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		//unmarshal
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		// if err := validator.New().Struct(req); err != nil {
		// 	validateErr := err.(validator.ValidationErrors)

		// 	log.Error("invalid request", sl.Err(err))
		// 	// render.JSON(w, r, resp.Error("invalid request"))
		// 	render.JSON(w, r, resp.ValidationError(validateErr))

		// 	return
		// }

		//TODO: добавить свою валидацию для тестов
		if err := ValidReq(&req); err != nil {
			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, err)

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLenght)
		}

		id, err := urlSaver.SaveURL(ctx, req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}
			log.Info("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url addd", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})

	}
}

func ValidReq(req *Request) error {
	if req.URL == "" {
		return errors.New("ErrInvalidURL")
	}

	if len(req.URL) >= 10 || len(req.URL) <= 5 {
		return errors.New("ErrLengts")
	}

	return nil
}
