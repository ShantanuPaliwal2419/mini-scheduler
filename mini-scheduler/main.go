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

	fmt.Println("Starting scheduler...")

	scheduler.Run(g, tasks)

	fmt.Println("\nFinal task states:")

	for name, t := range tasks {
		fmt.Printf("%s: %v\n", name, t.Status)
	}
}
