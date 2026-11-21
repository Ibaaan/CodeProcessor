package server

import (
	"codeproccessor/domain"
	"codeproccessor/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func SetHandlers(mux *http.ServeMux, service *service.Service) {
	mux.HandleFunc("POST /task", taskCreation(service))
	mux.HandleFunc("GET /status/{task_id}", taskStatus(service))
	mux.HandleFunc("GET /result/{task_id}", taskResult(service))
}

type taskResponse struct {
	TaskId string `json:"task_id"`
}

func taskCreation(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := s.CreateTask(r.Context())
		if errors.Is(err, domain.ErrTaskCreation) {
			slog.Error("cannot create task", "err", err)
			writeJSON(w, 400, "cannot create task")
			return
		} else if err != nil {
			slog.Error("INTERNAL ERROR", "err", err)
			writeJSON(w, 500, "")
			return
		}
		writeJSON(w, 201, taskResponse{TaskId: id.String()})
	}
}

type statusResponse struct {
	Status string `json:"status"`
}

func taskStatus(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("task_id"))
		if err != nil {
			writeJSON(w, 400, "wrong id")
			return
		}
		task, err := s.FindTask(r.Context(), id)
		if errors.Is(err, domain.ErrNotFound) {
			writeJSON(w, 404, "couldn't find task")
			return
		} else if err != nil {
			slog.Error("INTERNAL ERROR", "err", err)
			writeJSON(w, 500, "")
			return
		}
		writeJSON(w, 200, statusResponse{Status: string(task.Status)})
	}
}

type resultResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func taskResult(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("task_id"))
		if err != nil {
			writeJSON(w, 400, "wrong id")
			return
		}
		task, err := s.FindTask(r.Context(), id)
		if errors.Is(err, domain.ErrNotFound) {
			writeJSON(w, 404, "couldn't find task")
			return
		} else if err != nil {
			slog.Error("INTERNAL ERROR", "err", err)
			writeJSON(w, 500, "")
			return
		}
		writeJSON(w, 200, resultResponse{
			Status: string(task.Status),
			Result: task.Result.Data,
		})
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encoding response", "error", err)
	}
}
