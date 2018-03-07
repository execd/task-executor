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

// NewKubernetesImpl : build a new kubernetes manager impl
func NewKubernetesImpl(clientSet kubernetes.Interface) *KubernetesImpl {
	return &KubernetesImpl{clientSet: clientSet}
}

// KubernetesImpl : kubernetes implementation of a task manager
type KubernetesImpl struct {
	clientSet kubernetes.Interface
}

// ManageExecutingTask : manages the task with the given task id
func (s *KubernetesImpl) ManageExecutingTask(taskInfo task.Info, quit chan int) (*task.Info, error) {
	jobID := taskInfo.Metadata.(*v12.Job).Name
	fmt.Printf("Watching for events on task %s \n", jobID)
	opts := v1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", jobID),
	}
	jobs := s.clientSet.BatchV1().Jobs(v1.NamespaceDefault)
	w, err := jobs.Watch(opts)
	if err != nil {
		return nil, err
	}

	events := w.ResultChan()
	return s.handleEvent(&taskInfo, events, jobs)
}

func (s *KubernetesImpl) handleEvent(taskInfo *task.Info, events <-chan watch.Event, jobs v14.JobInterface) (*task.Info, error) {
	jobID := taskInfo.Metadata.(*v12.Job).Name
	for event := range events {
		switch event.Type {
		case watch.Deleted:
			return nil, errors.New("job has been deleted before execution completed")
		default:
			job := event.Object.(*v12.Job)
			if job.Status.Failed != 0 {
				fmt.Printf("Task %s failed.\n", job.Name)

				stats := s.findFailureCause(job)

				err := jobs.Delete(jobID, &v1.DeleteOptions{})
				if err != nil {
					return nil, fmt.Errorf("cleanup for failed task %s failed : %s", jobID, err.Error())
				}

				return buildResultInfo(taskInfo, job, false, stats), nil
			}
			if job.Status.Succeeded >= 1 {
				fmt.Printf("Task %s succeeded.\n", job.Name)

				err := jobs.Delete(jobID, &v1.DeleteOptions{})
				if err != nil {
					return nil, fmt.Errorf("cleanup for successful task %s failed : %s", jobID, err.Error())
				}

				return buildResultInfo(taskInfo, job, true, nil), nil
			}
		}
	}
	return nil, fmt.Errorf("an error occurred managing task %s", jobID)
}

func (s *KubernetesImpl) findFailureCause(failedJob *v12.Job) *task.FailureStatus {
	name := failedJob.Name
	core := s.clientSet.CoreV1()
	opts := v1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", name),
	}
	pods, err := core.Pods(v1.NamespaceDefault).List(opts)
	if err != nil {
		stats := &task.FailureStatus{
			Type:   "CouldNotFindCause",
			Reason: fmt.Sprintf("Failed to list pods related to job %s", name),
		}
		return stats
	}

	var parentStats *task.FailureStatus
	for _, pod := range pods.Items {
		var childStats task.FailureStatus
		parentStats = &task.FailureStatus{
			Type: "pod",
			Name: fmt.Sprintf("job (%s) pod (%s) failure", name, pod.Name),
		}

		for _, containerStat := range pod.Status.ContainerStatuses {
			childStats = task.FailureStatus{
				Type:    "container",
				Name:    fmt.Sprintf("job (%s) pod (%s) container (%s) failure", name, pod.Name, containerStat.Name),
				Message: containerStat.State.Terminated.Message,
				Reason:  containerStat.State.Terminated.Reason,
			}
		}
		parentStats.ChildStatus = append(parentStats.ChildStatus, childStats)
	}
	return parentStats
}

func buildResultInfo(taskInfo *task.Info, job *v12.Job, success bool, status *task.FailureStatus) *task.Info {
	taskInfo.Metadata = job
	taskInfo.Succeeded = success
	taskInfo.FailureStats = status
	return taskInfo
}
