package get

import (
	resp "api.quota-quick/api/internal/lib/api/responce"
	"api.quota-quick/api/internal/lib/logger/sl"
	"api.quota-quick/api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

//type Request struct {
//	Title       string `json:"title" validate:"required"`
//	Description string `json:"description"`
//	OwnerId     int    `json:"owner_id" validate:"required"`
//}

type Response struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name ContSaver
//go:generate go run mockery --dir=domain --output=domain/mocks --outpkg=mocks --all

type ContGetter interface {
	GetContainers(models.Container) error
	GetContainerById(uint64) (models.Container, error)
}

func GetById(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.containers.get.GetById"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		contId := chi.URLParam(r, "id")
		if contId == "" {
			log.Info("id is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		log.Info("request id checked")

		// TODO: validate id

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			//render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		//Работа с остальными полям

		err = contSaver.SaveContainer(req)
		if err != nil {
			log.Error("failed to save container", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save container"))
			return
		}

		log.Info("container saved", slog.Any("container", req))

		responceOk(w, r)
	}
}

//
//func GetAll(log *slog.Logger) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "handlers.containers.save.New"
//
//		log = log.With(
//			slog.String("op", op),
//			slog.String("request_id", middleware.GetReqID(r.Context())),
//		)
//
//		var req models.Container
//
//		err := render.DecodeJSON(r.Body, &req)
//		if err != nil {
//			log.Error("failed to decode request body", sl.Err(err))
//			render.JSON(w, r, resp.Error("failed to decode request"))
//			return
//		}
//
//		log.Info("request body decoded", slog.Any("request", req))
//
//		if err := validator.New().Struct(req); err != nil {
//			validateErr := err.(validator.ValidationErrors)
//			log.Error("invalid request", sl.Err(err))
//			//render.JSON(w, r, resp.Error("invalid request"))
//			render.JSON(w, r, resp.ValidationError(validateErr))
//			return
//		}
//		//Работа с остальными полям
//
//		err = contSaver.SaveContainer(req)
//		if err != nil {
//			log.Error("failed to save container", sl.Err(err))
//			render.JSON(w, r, resp.Error("failed to save container"))
//			return
//		}
//
//		log.Info("container saved", slog.Any("container", req))
//
//		responceOk(w, r)
//	}
//}
//
//func GetByUserId(log *slog.Logger) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "handlers.containers.save.New"
//
//		log = log.With(
//			slog.String("op", op),
//			slog.String("request_id", middleware.GetReqID(r.Context())),
//		)
//
//		var req models.Container
//
//		err := render.DecodeJSON(r.Body, &req)
//		if err != nil {
//			log.Error("failed to decode request body", sl.Err(err))
//			render.JSON(w, r, resp.Error("failed to decode request"))
//			return
//		}
//
//		log.Info("request body decoded", slog.Any("request", req))
//
//		if err := validator.New().Struct(req); err != nil {
//			validateErr := err.(validator.ValidationErrors)
//			log.Error("invalid request", sl.Err(err))
//			//render.JSON(w, r, resp.Error("invalid request"))
//			render.JSON(w, r, resp.ValidationError(validateErr))
//			return
//		}
//		//Работа с остальными полям
//
//		err = contSaver.SaveContainer(req)
//		if err != nil {
//			log.Error("failed to save container", sl.Err(err))
//			render.JSON(w, r, resp.Error("failed to save container"))
//			return
//		}
//
//		log.Info("container saved", slog.Any("container", req))
//
//		responceOk(w, r)
//	}
//}

func responceOk(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
