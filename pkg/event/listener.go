package event

import (
	"fmt"
	"github.com/wayofthepie/task-executor/pkg/executor"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"log"
	"time"
	"github.com/NeowayLabs/wabbit"
)


// Service : interface for execution of tasks
type Service interface {
	ListenForTasks() error
}

// ServiceImpl : default implementation of Service
type ServiceImpl struct {
	rabbit RabbitService
	executor executor.Service
}

// NewServiceImpl : build a ServiceImpl
func NewServiceImpl(rabbit RabbitService, executor executor.Service) (*ServiceImpl, error) {
	return &ServiceImpl{rabbit: rabbit, executor:executor}, nil
}

// ListenForTasks : listen for and execute tasks
func (s *ServiceImpl) ListenForTasks() error {
	fmt.Println("Listening for tasks")
	go func() {
		for {
			for msg := range s.rabbit.GetWorkQueueChan() {
				taskSpec := new(task.Spec)
				err := taskSpec.UnmarshalBinary(msg.Body())
				if err == nil {
					go s.handleMsg(msg)
				} else {
					log.Printf("error unmarshalling msg %s : %s", string(msg.Body()[:]), err.Error())
				}
			}
			fmt.Println("Stopped listening for messages, waiting 2 seconds for reconnect...")
			time.Sleep(time.Second * 2)
		}
	}()
	fmt.Println("Awaiting")

	return nil
}

func (s *ServiceImpl) handleMsg(msg wabbit.Delivery) {
	fmt.Println("received msg")
	taskSpec := new(task.Spec)
	err := taskSpec.UnmarshalBinary(msg.Body())
	if err == nil {
		executable := executor.NewExecutableImpl(taskSpec)
		s.executor.ExecuteTask(executable)
		msg.Ack(false)
	} else {
		log.Printf("error unmarshalling msg %s : %s", string(msg.Body()[:]), err.Error())
	}
}