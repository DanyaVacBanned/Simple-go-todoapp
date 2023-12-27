package updateTask

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
	TaskId          int64  `json:"task_id" validate:"required"`
	TaskName        string `json:"task_name,omitempty"`
	TaskDescription string `json:"task_description,omitempty" `
	TaskDone        int    `json:"task_done,omitempty"`
}

type Response struct {
	Resp   response.BaseResponse `json:"response"`
	TaskId int64                 `json:"task_id"`
}

type TaskUpdater interface {
	UpdateTask(taskId int64, taskName string, taskDescription string, taskDone int) (int64, error)
}

func New(logger *slog.Logger, updater TaskUpdater) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.updateTask.New"
		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("Failed to decode request body", slogLib.Err(err))
			render.JSON(writer, r, response.Error("Failed to decode request body"))
			return
		}
		logger.Info("Request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			logger.Error("Invalid request data", slogLib.Err(err))
			render.JSON(writer, r, response.Error("Invalid request data"))
			return
		}
		taskId, err := updater.UpdateTask(req.TaskId, req.TaskName, req.TaskDescription, req.TaskDone)
		render.JSON(writer, r, Response{
			Resp:   response.Ok(),
			TaskId: taskId,
		})
	}
}
