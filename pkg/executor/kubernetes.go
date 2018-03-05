package executor

import (
	"github.com/wayofthepie/task-executor/pkg/model/k8s"
	"github.com/wayofthepie/task-executor/pkg/model/task"

	"fmt"
	"github.com/wayofthepie/task-executor/pkg/manager"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesImpl struct {
	clientSet kubernetes.Interface
	manager   manager.Service
}

func NewKubernetesClientImpl(clientSet kubernetes.Interface, manager manager.Service) *KubernetesImpl {
	return &KubernetesImpl{clientSet: clientSet, manager: manager}
}

func (s *KubernetesImpl) ExecuteTask(spec *task.TaskSpec) (*task.TaskInfo, error) {
	container := v1.Container{
		Name:    spec.Name,
		Image:   spec.Image,
		Command: []string{spec.Init},
		Args:    spec.InitArgs,
	}

	k8sJob := k8s.Job(fmt.Sprintf("%s-", spec.Name), []v1.Container{container})
	batch := s.clientSet.BatchV1()

	fmt.Printf("Creating job for task %s", spec.Name)
	createdJob, err := batch.Jobs(v12.NamespaceDefault).Create(k8sJob)
	if err != nil {
		return nil, err
	}

	info, err := s.manager.ManageExecutingTask(createdJob.Name, make(chan int))
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *KubernetesImpl) GetExecutingTaskInfo(taskID string) (*task.TaskInfo, error) {
	job, err := s.clientSet.BatchV1().Jobs(v12.NamespaceDefault).Get(taskID, v12.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &task.TaskInfo{Id: job.Name, Metadata: job}, nil
}
