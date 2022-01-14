package queue

import "sync"

type Queue struct {
	queue []string
	mu    sync.Mutex
}

func (q *Queue) Empty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.queue) == 0
}

func (q *Queue) Push(path string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue = append(q.queue, path)
}

func (q *Queue) Pop() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	value := (q.queue)[0]
	q.queue = (q.queue)[1:]
	return value
}
