package goqueue

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrorQueueFull  = errors.New("queue was full")
	ErrorQueueEmpty = errors.New("queue was empty")
)

type Config struct {
	PushBlocking bool
	PopBlocking  bool
	MaxBuffer    int64
}

type InMemoryQueue struct {
	Config
	receiver chan Task
}

func NewInMemoryQueue(config Config) *InMemoryQueue {
	return &InMemoryQueue{
		Config:   config,
		receiver: make(chan Task, config.MaxBuffer),
	}
}

func (i *InMemoryQueue) Push(ctx context.Context, task Task) error {
	if i.PushBlocking {
		return i.BlockingPush(ctx, task)
	} else { // non-blocking
		return i.NonBlockingPush(ctx, task)
	}
}

func (i *InMemoryQueue) BlockingPush(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case i.receiver <- task:
		return nil
	}
}

func (i *InMemoryQueue) NonBlockingPush(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case i.receiver <- task:
		return nil
	default:
		return fmt.Errorf("push task failed, %w", ErrorQueueFull)
	}
}

func (i *InMemoryQueue) Pop(ctx context.Context) (Task, error) {
	if i.PopBlocking {
		return i.BlockingPop(ctx)
	} else {
		return i.NonBlockingPop(ctx)
	}
}

func (i *InMemoryQueue) BlockingPop(ctx context.Context) (Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case task := <-i.receiver:
		return task, nil
	}
}

func (i *InMemoryQueue) NonBlockingPop(ctx context.Context) (Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case task := <-i.receiver:
		return task, nil
	default:
		return nil, fmt.Errorf("pop task failed, %w", ErrorQueueEmpty)
	}
}
