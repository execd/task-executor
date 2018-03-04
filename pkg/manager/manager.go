package manager

import "github.com/wayofthepie/task-executor/pkg/model/task"

type Service interface {
	ManageExecutingTask(taskID string, quit chan int) (*task.TaskInfo, error)
}
