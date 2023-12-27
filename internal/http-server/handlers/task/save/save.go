package save

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
	TaskName        string `json:"task_name" validate:"required"`
	TaskDescription string `json:"task_description"`
}

type Response struct {
	Resp     response.BaseResponse `json:"response"`
	TaskName string                `json:"taskName,omitempty"`
}

type TaskCreator interface {
	CreateTask(taskName string, taskDescription string) (int64, error)
}

func New(logger *slog.Logger, taskCreator TaskCreator) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.save.New"
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
			logger.Error("invalid request", slogLib.Err(err))
			render.JSON(writer, r, response.ValidationError(validationErr))
			return
		}
		_, err = taskCreator.CreateTask(req.TaskName, req.TaskDescription)
		if err != nil {
			logger.Error("Failed to create task", slogLib.Err(err))
			render.JSON(writer, r, response.Error("Task name is required"))
			return
		}
		render.JSON(writer, r, Response{
			Resp:     response.Ok(),
			TaskName: req.TaskName,
		})
	}
}
