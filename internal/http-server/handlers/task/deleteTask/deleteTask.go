package deleteTask

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"todoapp/internal/lib/handlers/response"
	"todoapp/internal/lib/logger/slogLib"
)

type Request struct {
	Id int64 `json:"id" validate:"required"`
}

type Response struct {
	Resp response.BaseResponse `json:"response"`
}

type TaskDeleter interface {
	DeleteTask(taskId int64) error
}

func New(logger *slog.Logger, deleter TaskDeleter) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.deleteTask.New"
		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("Failed to decode request body", slogLib.Err(err))
			render.JSON(writer, r, response.Error("Failed to decode request"))
			return
		}
		logger.Info("Request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validationErr := err.(validator.ValidationErrors)
			logger.Error("Invalid request", slogLib.Err(validationErr))
			render.JSON(writer, r, response.ValidationError(validationErr))
			return
		}
		err = deleter.DeleteTask(req.Id)
		if err != nil {
			logger.Error("Failed to deleteTask task", slogLib.Err(err))
			render.JSON(writer, r, "Failed to deleteTask Task")
			return
		}
		render.JSON(writer, r, Response{
			Resp: response.Ok(),
		})
	}
}
