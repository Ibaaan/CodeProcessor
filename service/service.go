package service

import (
	"codeproccessor/domain"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateTask(ctx context.Context) (uuid.UUID, error)
	FindTask(ctx context.Context, id uuid.UUID) (domain.Task, error)
	UpdateTask(ctx context.Context, id uuid.UUID, status domain.TaskStatus) bool
}

type Service struct {
	r Repository
}

func NewTaskService(r Repository) Service {
	return Service{r: r}
}

func (s *Service) CreateTask(ctx context.Context) (uuid.UUID, error) {
	id, err := s.r.CreateTask(ctx)
	if err != nil {
		return id, err
	}
	go func() {
		timer := time.NewTimer(30 * time.Second)
		<-timer.C
		s.r.UpdateTask(ctx, id, domain.Ready)
		slog.Info("task updated", "id", id)
	}()
	slog.Info("task created", "id", id)
	return id, nil
}

func (s *Service) FindTask(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	task, err := s.r.FindTask(ctx, id)
	if err != nil {
		slog.Error("couldn't find", "id", id, "err", err)
		return task, err
	}
	slog.Info("successfully found", "task", task)
	return task, nil
}
