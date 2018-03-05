package main

import (
	"github.com/wayofthepie/task-executor/pkg/event"
	"github.com/wayofthepie/task-executor/pkg/execution"
	"github.com/wayofthepie/task-executor/pkg/executor"
	"github.com/wayofthepie/task-executor/pkg/k8s"
	"github.com/wayofthepie/task-executor/pkg/manager"
	"log"
)

func main() {
	clientSet := k8s.InitializeClientSet()
	kubeManager := manager.NewKubernetesImpl(clientSet)
	kubeExec := executor.NewKubernetesClientImpl(clientSet, kubeManager)
	conn, _ := event.NewRabbitConnection("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	exec, _ := execution.NewServiceImpl(ch, kubeExec)
	err := exec.ListenForTasks()
	if err != nil {
		log.Fatal(err.Error())
	}
	forever := make(chan bool)
	<-forever
}
