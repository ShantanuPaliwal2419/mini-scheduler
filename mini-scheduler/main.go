package main

import (
	"fmt"

	"mini-scheduler/graph"
	"mini-scheduler/scheduler"
	"mini-scheduler/task"
)

func main() {
	g := graph.NewGraph()

	// A → B
	// A → C
	// B → D
	// C → D

	g.Adddependency("A", "B")
	g.Adddependency("A", "C")
	g.Adddependency("B", "D")
	g.Adddependency("C", "D")

	tasks := map[string]*task.Task{
		"A": {
			Name:     "A",
			Status:   task.PENDING,
			Priority: 1,
		},
		"B": {
			Name:     "B",
			Status:   task.PENDING,
			Priority: 5,
		},
		"C": {
			Name:     "C",
			Status:   task.PENDING,
			Priority: 10,
		},
		"D": {
			Name:     "D",
			Status:   task.PENDING,
			Priority: 1,
		},
	}

	fmt.Println("Starting scheduler...")

	scheduler.Run(g, tasks)

	fmt.Println("\nFinal status:")

	for name, t := range tasks {
		fmt.Printf(
			"%s -> status=%d priority=%d\n",
			name,
			t.Status,
			t.Priority,
		)
	}
}
