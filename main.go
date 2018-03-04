package main

import (
	"fmt"
	"github.com/wayofthepie/task-executor/pkg/k8s"
	"github.com/wayofthepie/task-executor/pkg/executor"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"time"
	"log"
	"k8s.io/api/batch/v1"
)

func main() {
	clientSet := k8s.InitializeClientSet()
	exec := executor.NewKubernetesClientImpl(clientSet)
	spec := &task.TaskSpec{Name:"test-execution", Image:"alpine", Init:"sleep", InitArgs:[]string{"100"}}
	job, err := exec.ExecuteTask(spec)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(job.Metadata.(*v1.Job).Name)
	fmt.Println("Sleeping...")
	time.Sleep(time.Second * 10)
	info, err := exec.GetExecutingTaskInfo(job.Metadata.(*v1.Job).Name)
	if err!= nil {
		log.Fatal(err.Error())
	}

	fmt.Println(info.Metadata.(*v1.Job).Name)
	fmt.Println(info.Metadata.(*v1.Job).Spec.Template.Spec.Containers[0].Command)

	fmt.Println("Executing!")
}
