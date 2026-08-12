package scheduler

import (
	"mini-scheduler/graph"
	"mini-scheduler/task"
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
