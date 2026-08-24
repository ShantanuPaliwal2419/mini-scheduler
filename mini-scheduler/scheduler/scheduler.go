package scheduler

import (
	executor "mini-scheduler/executer"
	"mini-scheduler/graph"
	"mini-scheduler/queue"
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
func Run(
	g *graph.Graph,
	tasks map[string]*task.Task,
) {
	for {
		readyTasks := GetReadyTasks(g, tasks)
		if len(readyTasks) == 0 {
			break
		}
		// pushing ready tasks to the priority queue
		pq := queue.NewPriorityQueue()
		for _, taskName := range readyTasks {
			pq.Push(tasks[taskName])
		}
		var readytasks []*task.Task
		// Process the ready tasks (implementation for task execution would go here)
		for pq.Len() != 0 {
			t := pq.Pop()
			readytasks = append(readytasks, t)
		}
		executor.Execute(readytasks)

	}
}
