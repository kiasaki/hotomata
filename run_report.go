package hotomata

type RunReport struct {
	Completed       bool
	TaskStatusCount map[TaskStatus]int
}
