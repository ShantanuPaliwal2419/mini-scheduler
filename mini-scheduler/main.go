package main

import (
	"fmt"

	"mini-scheduler/queue"
	"mini-scheduler/task"
)

func main() {
	q := queue.NewPriorityQueue()

	a := &task.Task{
		Name:     "A",
		Priority: 5,
	}

	b := &task.Task{
		Name:     "B",
		Priority: 10,
	}

	c := &task.Task{
		Name:     "C",
		Priority: 3,
	}

	q.Push(a)
	q.Push(b)
	q.Push(c)

	fmt.Println("Queue length:", q.Len())

	for q.Len() > 0 {
		t := q.Pop()
		fmt.Printf("Executing: %s | Priority: %d\n", t.Name, t.Priority)
	}

	fmt.Println("Queue length after Pop:", q.Len())

	// Test empty queue
	t := q.Pop()

	if t == nil {
		fmt.Println("Empty queue: Pop returned nil")
	}
}
