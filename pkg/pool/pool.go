// Package pool - обертка над sync.Pool для объектов с методом Reset().
package pool

import (
	"errors"
	"sync"
)

type Resetter interface {
	Reset()
}

// Pool хранит объекты типа T и сбрасывает их состояние перед возвратом в пул.
type Pool[T Resetter] struct {
	pool sync.Pool
}

func New[T Resetter](newFn func() T) (*Pool[T], error) {
	if newFn == nil {
		return nil, errors.New("pool.New: newFn must not be nil")
	}

	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return newFn()
			},
		},
	}, nil
}

// Get возвращает объект из пула.
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put сбрасывает состояние объекта и возвращает его в пул.
func (p *Pool[T]) Put(v T) {
	v.Reset()
	p.pool.Put(v)
}
