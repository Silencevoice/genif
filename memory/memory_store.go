package memory

import (
	"context"
	"errors"
	"sync"
)

type MemStore[T any] struct {
	sync.RWMutex
	data map[string]T
}

func NewMemStore[T any]() *MemStore[T] {
	return &MemStore[T]{
		data: make(map[string]T),
	}
}

func (m *MemStore[T]) GetByID(ctx context.Context, id string) (*T, error) {
	m.RLock()
	defer m.RUnlock()

	ent, ok := m.data[id]
	if !ok {
		return nil, errors.New("entity not found")
	}

	return &ent, nil
}

func (m *MemStore[T]) GetMultipleByID(ctx context.Context, ids []string) ([]*T, error) {
	m.RLock()
	defer m.RUnlock()

	ents := make([]*T, len(ids))
	for idx, id := range ids {
		ent, ok := m.data[id]
		if !ok {
			return nil, errors.New("entity not found")
		}
		ents[idx] = &ent
	}

	return ents, nil
}

func (m *MemStore[T]) GetAll(ctx context.Context) ([]*T, error) {
	m.RLock()
	defer m.RUnlock()

	ents := []*T{}
	for _, value := range m.data {
		ents = append(ents, &value)
	}

	return ents, nil
}

func (m *MemStore[T]) Insert(ctx context.Context, id string, entity *T) (*T, error) {
	m.Lock()
	defer m.Unlock()

	_, ok := m.data[id]
	if ok {
		return nil, errors.New("already existing key")
	}

	m.data[id] = *entity

	return entity, nil
}

func (m *MemStore[T]) Delete(ctx context.Context, id string) error {
	m.Lock()
	defer m.Unlock()

	_, ok := m.data[id]
	if !ok {
		return errors.New("entity not found")
	}

	delete(m.data, id)

	return nil
}

func (m *MemStore[T]) Update(ctx context.Context, id string, entity *T) error {
	m.Lock()
	defer m.Unlock()

	_, ok := m.data[id]
	if !ok {
		return errors.New("entity not found")
	}

	m.data[id] = *entity
	return nil
}

func (m *MemStore[T]) ExecuteQuery(ctx context.Context, f func(ctx context.Context, data map[string]T) ([]*T, error)) ([]*T, error) {
	m.RLock()
	defer m.RUnlock()
	return f(ctx, m.data)
}

func (m *MemStore[T]) ExecuteUpdate(ctx context.Context, f func(ctx context.Context, data map[string]T) (int, error)) (int, error) {
	m.Lock()
	defer m.Unlock()
	return f(ctx, m.data)
}
