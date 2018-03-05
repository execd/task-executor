// Package manager exposes implementations of task managers for
// specific platforms
package manager

import "github.com/wayofthepie/task-executor/pkg/model/task"

// Service : defined methods necessary for a manager
type Service interface {
	ManageExecutingTask(taskID string, quit chan int) (*task.Info, error)
}
