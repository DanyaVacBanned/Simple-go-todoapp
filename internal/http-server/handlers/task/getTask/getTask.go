package getTask

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"todoapp/internal/lib/handlers/response"
	"todoapp/internal/lib/logger/slogLib"
	"todoapp/internal/storage/sqlite"
)

type Response struct {
	Resp response.BaseResponse `json:"response"`
	Task []sqlite.Task         `json:"task"`
}

type TaskGetter interface {
	GetTask(taskId string) ([]sqlite.Task, error)
	GetTasks() ([]sqlite.Task, error)
}

func New(logger *slog.Logger, getter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.getTask.New"
		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		taskId := r.URL.Query().Get("task_id")

		fmt.Println("task_id = ", taskId)
		if taskId != "" {
			task, err := getter.GetTask(taskId)
			if err != nil {
				logger.Error("No such task", slogLib.Err(err))
				render.JSON(w, r, response.Error("No such task"))
				return
			}
			render.JSON(w, r, Response{
				Resp: response.Ok(),
				Task: task,
			})
			return
		}
		tasks, err := getter.GetTasks()
		if err != nil {
			logger.Error("Cant get tasks", slogLib.Err(err))
			render.JSON(w, r, response.Error("Cant get tasks"))
			return
		}
		render.JSON(w, r, Response{
			Resp: response.Ok(),
			Task: tasks,
		})
	}
}
