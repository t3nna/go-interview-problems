package main

import (
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	data []int
	size int
	mu   sync.RWMutex
}

func NewQueue(size int) *Queue {
	return &Queue{
		data: make([]int, 0),
		size: size,
	}
}

func (q *Queue) Push(val int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.data) >= q.size {
		return ErrQueueFull
	}
	q.data = append(q.data, val)

	return nil

}

func (q *Queue) Pop() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) <= 0 {
		return -1
	}
	res := q.data[0]
	q.data = q.data[1:]
	return res
}

func (q *Queue) Peek() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.data) <= 0 {
		return -1
	}

	return q.data[0]
}
