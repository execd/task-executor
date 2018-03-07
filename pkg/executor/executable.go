package executor

import (
	"github.com/wayofthepie/task-executor/pkg/model/task"
)

// Executable : interface for executables
type Executable interface {
	GetName() string
	GetTask() *task.Spec
}

// ExecutableImpl : implementation of an executable
type ExecutableImpl struct {
	task *task.Spec
}

// NewExecutableImpl : build a new executable
func NewExecutableImpl(task *task.Spec) *ExecutableImpl {
	return &ExecutableImpl{task: task}
}

// GetName : get the name of this executable
func (e *ExecutableImpl) GetName() string {
	return e.task.Name
}

// GetTask : get the task associated with this executable
func (e *ExecutableImpl) GetTask() *task.Spec {
	return e.task
}
