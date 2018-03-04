package manager

import (
	"k8s.io/client-go/kubernetes"
)

type KubernetesImpl struct {
	clientset kubernetes.Interface
}

func (s *KubernetesImpl) ManageExecutingTask(taskID string) {}
