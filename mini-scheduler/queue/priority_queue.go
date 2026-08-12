package queue

import "mini-scheduler/task"

type PriorityQueue struct {
	tasks []*task.Task
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		tasks: make([]*task.Task, 0),
	}
}

func (q *PriorityQueue) Push(t *task.Task) {
	q.tasks = append(q.tasks, t)
}

func (q *PriorityQueue) Pop() *task.Task {
	if len(q.tasks) == 0 {
		return nil
	}

	highestIndex := 0

	for i := 1; i < len(q.tasks); i++ {
		if q.tasks[i].Priority > q.tasks[highestIndex].Priority {
			highestIndex = i
		}
	}

	task := q.tasks[highestIndex]

	q.tasks = append(
		q.tasks[:highestIndex],
		q.tasks[highestIndex+1:]...,
	)

	return task
}

func (q *PriorityQueue) Len() int {
	return len(q.tasks)
}
