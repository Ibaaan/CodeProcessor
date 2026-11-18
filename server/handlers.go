package server

import (
	"codeproccessor/service"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func SetHandlers(mux *http.ServeMux, service service.Service) {
	mux.HandleFunc("POST /task", taskCreation(service))
	mux.HandleFunc("GET /status/{task_id}", taskStatus(service))
	mux.HandleFunc("GET /result/{task_id}", taskResult(service))
}

func taskCreation(s service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := s.CreateTask(r.Context())
		if err != nil {
			writeJSON(w, 400, "ure stupid")
			return
		}
		writeJSON(w, 200, map[string]string{"uuid": id.String()})
	}
}

func taskStatus(s service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("task_id"))
		if err != nil {
			writeJSON(w, 400, "wrong id")
			return
		}
		task, err := s.FindTask(r.Context(), id)
		if err != nil {
			writeJSON(w, 400, "couldn't find task")
			return
		}
		writeJSON(w, 200, map[string]string{"status": string(task.Status)})
	}
}

func taskResult(s service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("task_id"))
		if err != nil {
			writeJSON(w, 400, "wrong id")
			return
		}
		task, err := s.FindTask(r.Context(), id)
		if err != nil {
			writeJSON(w, 400, "couldn't find task")
			return
		}
		writeJSON(w, 200, map[string]string{"status": string(task.Status)})
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encoding response", "error", err)
	}
}
