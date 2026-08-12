package scheduler

import (
	"fmt"
	"mini-scheduler/graph"
	"mini-scheduler/task"
	"time"
)

func GetReadyTasks(
	g *graph.Graph,
	tasks map[string]*task.Task,
) []string {
	var readyTasks []string

	for name, t := range tasks {

		if t.Status != task.PENDING {
			continue
		}

		parents := g.Parents(name)

		if len(parents) == 0 {
			readyTasks = append(readyTasks, name)
			continue
		}
		allParentsDone := true

		for _, parent := range parents {
			if tasks[parent].Status != task.SUCCESS {
				allParentsDone = false
				break
			}
		}

		if allParentsDone {
			readyTasks = append(readyTasks, name)
		}
	}

	return readyTasks
}
func Run(
	g *graph.Graph,
	tasks map[string]*task.Task,
) {
	for {
		readyTasks := GetReadyTasks(g, tasks)
		if len(readyTasks) == 0 {
			break
		}
		// Process the ready tasks (implementation for task execution would go here)
		for _, taskName := range readyTasks {
			tasks[taskName].Status = task.RUNNING
			fmt.Println("task is executing ...")
			time.Sleep(300 * time.Millisecond) // Simulate task execution time
			tasks[taskName].Status = task.SUCCESS
			fmt.Printf("Task %s completed successfully\n", taskName)
		}
	}
}
