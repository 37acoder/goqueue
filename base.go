package goqueue

import "context"

type Queue interface {
	Push(ctx context.Context, task Task) error
	Pop(ctx context.Context) (Task, error)
}

type Task interface {
	Execute(ctx context.Context) error
}
