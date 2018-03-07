package executor

import (
	"github.com/wayofthepie/task-executor/pkg/model/task"
)

type Executable interface {
	GetName() string
	GetTask() *task.Spec
}

type ExecutableImpl struct {
	task *task.Spec
}

func NewExecutableImpl(task *task.Spec) *ExecutableImpl {
	return &ExecutableImpl{task: task}
}

func (e *ExecutableImpl) GetName() string {
	return e.task.Name
}

func (e *ExecutableImpl) GetTask() *task.Spec {
	return e.task
}
