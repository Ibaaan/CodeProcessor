package domain

import "github.com/google/uuid"

type TaskStatus string

const (
	Pending TaskStatus = "pending"
	Ready   TaskStatus = "ready"
)

type Task struct {
	Id     uuid.UUID
	Status TaskStatus
}
