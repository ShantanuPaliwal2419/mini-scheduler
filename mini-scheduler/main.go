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
		"A": {Name: "A", Status: task.PENDING},
		"B": {Name: "B", Status: task.PENDING},
		"C": {Name: "C", Status: task.PENDING},
		"D": {Name: "D", Status: task.PENDING},
	}

	// Initially only A should be ready.
	fmt.Println("Ready:", scheduler.GetReadyTasks(g, tasks))

	// A completes.
	tasks["A"].Status = task.SUCCESS

	// Now B and C should be ready.
	fmt.Println("After A:", scheduler.GetReadyTasks(g, tasks))

	// Only B completes.
	tasks["B"].Status = task.SUCCESS

	// D should NOT be ready because C hasn't completed.
	fmt.Println("After B:", scheduler.GetReadyTasks(g, tasks))

	// C completes.
	tasks["C"].Status = task.SUCCESS

	// Now D should be ready.
	fmt.Println("After C:", scheduler.GetReadyTasks(g, tasks))
}
