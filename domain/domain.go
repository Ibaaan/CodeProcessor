package domain

import "github.com/google/uuid"

type TaskStatus string

const (
	InProgress TaskStatus = "in_progress"
	Ready      TaskStatus = "ready"
)

type Task struct {
	Id          uuid.UUID
	Status      TaskStatus
	CodeRequest CodeRequest
	Result      Result
}

type Result struct {
	Data string
}

type CodeRequest struct {
	Trans string `json:"translator" validate:"required" minLength:"1"`
	Code  string `json:"code" validate:"required" minLength:"1"`
}
