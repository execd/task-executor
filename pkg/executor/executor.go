package executor

import "github.com/wayofthepie/task-executor/pkg/model/task"

type Service interface{
	ExecuteTask(spec *task.TaskSpec) (*task.TaskInfo, error)
	GetExecutingTaskInfo(taskID string) (*task.TaskInfo, error)
}
