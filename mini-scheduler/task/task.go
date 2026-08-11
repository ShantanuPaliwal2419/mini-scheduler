package task

const (
	PENDING = Status("PENDING")
	READY   = Status("READY")
	RUNNING = Status("RUNNING")
	SUCCESS = Status("SUCCESS")
	FAILED  = Status("FAILED")
)

type Status string
type Task struct {
	id string
}
