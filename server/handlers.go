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

// taskCreation creates a task and starts background processing.
// @Summary Create a task
// @Description No request body is required. Any supplied body is ignored.
// @Description Returns task_id as a canonical 36-character UUIDv7 string (128 bits, not a 64-bit number).
// @Tags tasks
// @ID createTask
// @Produce json
// @Success 201 {object} taskResponse "Created task with its UUIDv7 task_id"
// @Failure 400 {string} string "cannot create task"
// @Failure 500 {string} string "Internal error; response body is an empty JSON string"
// @Router /task [post]
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

// @Summary Get task status
// @Description Status is in_progress or ready. Generated task IDs are UUIDv7; other valid UUID versions are also accepted for lookup.
// @Tags tasks
// @ID getTaskStatus
// @Produce json
// @Param task_id path string true "Task UUID returned by POST /task" Format(uuid)
// @Success 200 {object} statusResponse "Task status: in_progress or ready"
// @Failure 400 {string} string "wrong id"
// @Failure 404 {string} string "couldn't find task"
// @Failure 500 {string} string "Internal error; response body is an empty JSON string"
// @Router /status/{task_id} [get]
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

// @Summary Get task result
// @Description Returns 200 with status in_progress and an empty result before completion; status ready and the placeholder result nothing after completion.
// @Description Generated task IDs are UUIDv7; other valid UUID versions are also accepted for lookup.
// @Tags tasks
// @ID getTaskResult
// @Produce json
// @Param task_id path string true "Task UUID returned by POST /task" Format(uuid)
// @Success 200 {object} resultResponse "Task status (in_progress or ready) and result"
// @Failure 400 {string} string "wrong id"
// @Failure 404 {string} string "couldn't find task"
// @Failure 500 {string} string "Internal error; response body is an empty JSON string"
// @Router /result/{task_id} [get]
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
