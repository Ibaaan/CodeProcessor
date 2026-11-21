package domain

import "github.com/google/uuid"

type TaskStatus string

const (
	InProgress TaskStatus = "in_progress"
	Ready      TaskStatus = "ready"
)

type Task struct {
	Id     uuid.UUID
	Status TaskStatus
	Result Result
}

type Result struct {
	Data string
}
