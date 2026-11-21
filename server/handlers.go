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

type Handlers struct {
	s *service.Service
}

func NewHandlers(s *service.Service) *Handlers {
	return &Handlers{s: s}
}

func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /task", h.taskCreation)
	mux.HandleFunc("GET /status/{task_id}", h.taskStatus)
	mux.HandleFunc("GET /result/{task_id}", h.taskResult)
}

type taskResponse struct {
	TaskId string `json:"task_id"`
}

func (h *Handlers) taskCreation(w http.ResponseWriter, r *http.Request) {
	id, err := h.s.CreateTask(r.Context())
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

type statusResponse struct {
	Status string `json:"status"`
}

func (h *Handlers) taskStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("task_id"))
	if err != nil {
		writeJSON(w, 400, "wrong id")
		return
	}
	task, err := h.s.FindTask(r.Context(), id)
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

type resultResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func (h *Handlers) taskResult(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("task_id"))
	if err != nil {
		writeJSON(w, 400, "wrong id")
		return
	}
	task, err := h.s.FindTask(r.Context(), id)
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

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encoding response", "error", err)
	}
}
