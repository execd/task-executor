// Package manager exposes implementations of task managers for
// specific platforms
package manager

import "github.com/wayofthepie/task-executor/pkg/model/task"

// Service : defined methods necessary for a manager
type Service interface {
	ManageExecutingTask(taskInfo task.Info, quit chan int) (*task.Info, error)
}
