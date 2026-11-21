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
	UpdateTask(ctx context.Context, id uuid.UUID, status domain.TaskStatus) (bool, error)
}

type Service struct {
	r      Repository
	appCtx context.Context
}

func NewTaskService(r Repository, ctx context.Context) *Service {
	return &Service{r: r, appCtx: ctx}
}

func (s *Service) CreateTask(ctx context.Context) (uuid.UUID, error) {
	id, err := s.r.CreateTask(ctx)
	if err != nil {
		slog.Error("failed to create task", "err", err)
		return uuid.Nil, err
	}
	go s.simulatingWork(id)
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

func (s *Service) simulatingWork(id uuid.UUID) {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		if ok, err := s.r.UpdateTask(s.appCtx, id, domain.Ready); !ok {
			slog.Error("cannot update task", "id", id, "err", err)
			return
		}
		slog.Info("task updated", "id", id)
		return
	case <-s.appCtx.Done():
		slog.Error("context canceled: cannot update task", "id", id)
		return
	}
}
