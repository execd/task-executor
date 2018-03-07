package executor

import (
	"github.com/wayofthepie/task-executor/pkg/model/task"
)

// Service : defines methods for an executor
type Service interface {
	ExecuteTask(executable Executable) (*task.Info, error)
	GetExecutingTaskInfo(taskID string) (*task.Info, error)
}
