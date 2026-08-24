package executor

import (
	"fmt"
	"sync"
	"time"

	"mini-scheduler/task"
)

func Execute(tasks []*task.Task) {
	var wg sync.WaitGroup

	for _, t := range tasks {
		wg.Add(1)

		go func(t *task.Task) {
			defer wg.Done()
			start := time.Now()

			t.Status = task.RUNNING

			fmt.Printf(
				"Task %s started (priority=%d)\n",
				t.Name,
				t.Priority,
			)

			time.Sleep(300 * time.Millisecond)

			t.Status = task.SUCCESS
			t.DurationMs = time.Since(start).Milliseconds()
			fmt.Printf(
				"Task %s completed\n duration=%dms\n ",
				t.Name,
				t.DurationMs,
			)
		}(t)
	}

	wg.Wait()
}
