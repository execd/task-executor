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

// KubernetesImpl : kubernetes implementation of an executor
type KubernetesImpl struct {
	clientSet kubernetes.Interface
	manager   manager.Service
}

// NewKubernetesClientImpl : build a KubernetesImpl
func NewKubernetesClientImpl(clientSet kubernetes.Interface, manager manager.Service) *KubernetesImpl {
	return &KubernetesImpl{clientSet: clientSet, manager: manager}
}

// ExecuteTask : execute and executable
func (s *KubernetesImpl) ExecuteTask(executable Executable) (*task.Info, error) {
	spec := executable.GetTask()
	container := v1.Container{
		Name:    spec.Name,
		Image:   spec.Image,
		Command: []string{spec.Init},
		Args:    spec.InitArgs,
	}

	k8sJob := k8s.Job(fmt.Sprintf("%s-", spec.Name), []v1.Container{container})
	batch := s.clientSet.BatchV1()

	fmt.Printf("Creating job for task %s\n", spec.Name)
	createdJob, err := batch.Jobs(v12.NamespaceDefault).Create(k8sJob)
	if err != nil {
		return nil, err
	}

	initialTaskInfo := task.Info{
		ID:       executable.GetTask().ID,
		Metadata: createdJob,
	}

	info, err := s.manager.ManageExecutingTask(initialTaskInfo, make(chan int))
	if err != nil {
		return nil, err
	}

	return info, nil
}
