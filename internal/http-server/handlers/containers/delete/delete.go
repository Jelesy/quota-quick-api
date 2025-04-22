package delete

import (
	resp "api.quota-quick/api/internal/lib/api/responce"
	"api.quota-quick/api/internal/lib/logger/sl"
	"api.quota-quick/api/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

type Request struct {
	Id uint64 `json:"id" validate:"required"`
}

type Response struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name ContSaver
//go:generate go run mockery --dir=domain --output=domain/mocks --outpkg=mocks --all

type ContDeleter interface {
	//GetContainers(models.Container) error
	DeleteContainerById(uint64) error
}

func DeleteById(log *slog.Logger, contGetter ContDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.containers.get.DeleteById"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		stringContId := chi.URLParam(r, "id")
		if stringContId == "" {
			log.Info("id is empty")
			render.JSON(w, r, resp.Error("invalid request, id is empty"))
			return
		}

		log.Info("request id checked")

		contId, err := strconv.Atoi(stringContId)
		if err != nil {
			log.Info("id is invalid format")
			render.JSON(w, r, resp.Error("invalid request, id is invalid format"))
			return
		}

		req := Request{Id: uint64(contId)}

		// TODO: validate id

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			//render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		//Работа с остальными полям

		err = contGetter.DeleteContainerById(req.Id)
		if errors.Is(err, storage.ErrContainerNotFound) {
			log.Info("container not found", "cont id", req.Id)
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to save container", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("got container", slog.Any("container", req))

		responceOk(w, r)
	}
}

func responceOk(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
