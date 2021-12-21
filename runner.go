package hotomata

import (
	"bytes"
)

type TaskAction string

const (
	TaskActionAbort    = "abort"
	TaskActionContinue = "continue"
)

type TaskStatus string

const (
	TaskStatusSuccess = "success"
	TaskStatusWarning = "warning"
	TaskStatusError   = "error"
	TaskStatusSkip    = "skip"
)

type TaskResponse struct {
	Log    *bytes.Buffer
	Action TaskAction `json:"action"`
	Status TaskStatus `json:"status"`
}

func (r TaskResponse) Color() Color {
	switch r.Status {
	case TaskStatusSuccess:
		return ColorGreen
	case TaskStatusError:
		return ColorRed
	case TaskStatusWarning:
		return ColorYellow
	case TaskStatusSkip:
		return ColorBlue

	}
	return ColorNone
}

type Runner interface {
	Run(string) *TaskResponse
}
