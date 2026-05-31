// Package pool provides a generic wrapper around sync.Pool.
//
// Objects stored in the pool must implement the Reset method.
// The pool automatically resets objects before returning them back
// for reuse.
package pool

import "sync"

type Resetter interface {
	Reset()
}

type Pool[T Resetter] struct {
	pool sync.Pool
}

func New[T Resetter](newFn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return newFn()
			},
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(v T) {
	v.Reset()
	p.pool.Put(v)
}
