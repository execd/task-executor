package main

import (
	"github.com/wayofthepie/task-executor/pkg/event"
	"github.com/wayofthepie/task-executor/pkg/executor"
	"github.com/wayofthepie/task-executor/pkg/k8s"
	"github.com/wayofthepie/task-executor/pkg/manager"
)

func main() {
	clientSet := k8s.InitializeClientSet()
	kubeManager := manager.NewKubernetesImpl(clientSet)
	kubeExec := executor.NewKubernetesClientImpl(clientSet, kubeManager)
	rabbit, _ := event.NewRabbitServiceImpl("amqp://guest:guest@localhost:5672/")
	exec, _ := event.NewServiceImpl(rabbit, kubeExec)
	exec.ListenForTasks()

	forever := make(chan bool)
	<-forever
}
