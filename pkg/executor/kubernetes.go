package executor

import (
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"github.com/wayofthepie/task-executor/pkg/model/k8s"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)

type KubernetesImpl struct {
	clientSet kubernetes.Interface
}

func NewKubernetesClientImpl(clientSet kubernetes.Interface) *KubernetesImpl {
	return &KubernetesImpl{clientSet: clientSet}
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

	createdJob, err := batch.Jobs(v12.NamespaceDefault).Create(k8sJob)
	if err != nil {
		return nil, err
	}

	return &task.TaskInfo{Id: createdJob.Name, Metadata: createdJob}, nil
}

func (s *KubernetesImpl) GetExecutingTaskInfo(taskID string) (*task.TaskInfo, error) {
	job, err := s.clientSet.BatchV1().Jobs(v12.NamespaceDefault).Get(taskID, v12.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &task.TaskInfo{Id: job.Name, Metadata: job}, nil
}
