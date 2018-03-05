package manager

import (
	"errors"
	"fmt"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	v12 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	v14 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

func NewKubernetesImpl(clientSet kubernetes.Interface) *KubernetesImpl {
	return &KubernetesImpl{clientSet: clientSet}
}

type KubernetesImpl struct {
	clientSet kubernetes.Interface
}

func (s *KubernetesImpl) ManageExecutingTask(taskID string, quit chan int) (*task.TaskInfo, error) {
	fmt.Printf("Watching for events on task %s", taskID)
	opts := v1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", taskID),
	}
	jobs := s.clientSet.BatchV1().Jobs(v1.NamespaceDefault)
	w, err := jobs.Watch(opts)
	if err != nil {
		return nil, err
	}

	events := w.ResultChan()
	return handleEvent(taskID, events, jobs)
}

func handleEvent(taskID string, events <-chan watch.Event, jobs v14.JobInterface) (*task.TaskInfo, error) {
	for event := range events {
		switch event.Type {
		case watch.Deleted:
			return nil, errors.New("job has been deleted before execution completed")
		default:
			j := event.Object.(*v12.Job)
			if j.Status.Failed != 0 {
				err := jobs.Delete(taskID, &v1.DeleteOptions{})
				if err != nil {
					return nil, fmt.Errorf("cleanup for failed task %s failed : %s", taskID, err.Error())
				}

				return &task.TaskInfo{Id: taskID, Metadata: j}, nil
			}
			if j.Status.Succeeded > 1 {
				err := jobs.Delete(taskID, &v1.DeleteOptions{})
				if err != nil {
					return nil, fmt.Errorf("cleanup for successful task %s failed : %s", taskID, err.Error())
				}
				return &task.TaskInfo{Id: taskID, Metadata: j}, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("An error occurred managing task %s", taskID))
}

