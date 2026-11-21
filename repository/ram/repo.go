package ram

import (
	"codeproccessor/domain"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type RamRepo struct {
	rwm sync.RWMutex
	m   map[uuid.UUID]domain.Task
}

func NewRamRepository() *RamRepo {
	ma := make(map[uuid.UUID]domain.Task)
	return &RamRepo{m: ma, rwm: sync.RWMutex{}}
}

func (r *RamRepo) CreateTask(ctx context.Context) (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w:%w", domain.ErrTaskCreation, err)
	}

	r.rwm.Lock()
	defer r.rwm.Unlock()

	r.m[id] = domain.Task{
		Id:     id,
		Status: domain.InProgress,
	}

	return id, nil
}

func (r *RamRepo) FindTask(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	r.rwm.RLock()
	defer r.rwm.RUnlock()

	task, ok := r.m[id]
	if !ok {
		return task, domain.ErrNotFound
	}
	return task, nil
}

func (r *RamRepo) UpdateTask(ctx context.Context, task domain.Task) (bool, error) {
	r.rwm.Lock()
	defer r.rwm.Unlock()

	_, ok := r.m[task.Id]
	if !ok {
		return false, domain.ErrNotFound
	}

	r.m[task.Id] = task

	return true, nil
}
