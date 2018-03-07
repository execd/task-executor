package event

import (
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"github.com/wayofthepie/task-executor/pkg/executor"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"log"
	"time"
)

// Service : interface for execution of tasks
type Service interface {
	ListenForTasks() error
}

// ServiceImpl : default implementation of Service
type ServiceImpl struct {
	rabbit   RabbitService
	executor executor.Service
}

// NewServiceImpl : build a ServiceImpl
func NewServiceImpl(rabbit RabbitService, executor executor.Service) (*ServiceImpl, error) {
	return &ServiceImpl{rabbit: rabbit, executor: executor}, nil
}

// ListenForTasks : listen for and execute tasks
func (s *ServiceImpl) ListenForTasks() error {
	fmt.Println("Listening for tasks")
	go func() {
		for {
			for msg := range s.rabbit.GetWorkQueueChan() {
				go s.handleMsg(msg)
			}
			fmt.Println("Stopped listening for messages, waiting 2 seconds for reconnect...")
			time.Sleep(time.Second * 2)
		}
	}()
	return nil
}

func (s *ServiceImpl) handleMsg(msg wabbit.Delivery) {
	fmt.Println("received msg")
	taskSpec := new(task.Spec)
	err := taskSpec.UnmarshalBinary(msg.Body())
	if err == nil {
		executable := executor.NewExecutableImpl(taskSpec)
		info, err := s.executor.ExecuteTask(executable)
		if err != nil {
			fmt.Printf("Failed to execute task: %s\n", err.Error())
		}

		err = s.rabbit.PublishTaskStatus(info)
		if err != nil {
			fmt.Printf("Failed to publish task status: %s\n", err.Error())
		}
	} else {
		log.Printf("error unmarshalling msg %s : %s", string(msg.Body()[:]), err.Error())
		unmarshalFailureInfo := &task.Info{
			Metadata: string(msg.Body()),
			FailureStats: &task.FailureStatus{
				Name:    "unmarshal failure",
				Type:    "UnmarshalFailure",
				Reason:  "failed to unmarshal provided task spec",
				Message: err.Error(),
			},
		}
		log.Println("Publishing failure")

		err = s.rabbit.PublishTaskStatus(unmarshalFailureInfo)
		if err != nil {
			fmt.Printf("Failed to publish task unmarshal failure status: %s\n", err.Error())
		}
	}
	msg.Ack(false)
}
