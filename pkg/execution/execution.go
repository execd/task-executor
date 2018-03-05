package execution

import (
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"github.com/wayofthepie/task-executor/pkg/executor"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"log"
)

const queueName = "work_queue"

type Service interface {
	ListenForTasks() error
}

type ServiceImpl struct {
	channel  wabbit.Channel
	queue    wabbit.Queue
	executor executor.Service
}

func NewServiceImpl(channel wabbit.Channel, ex executor.Service) (*ServiceImpl, error) {
	q, err := channel.QueueDeclare(
		queueName,
		wabbit.Option{
			"durable":    true,
			"autoDelete": false,
			"exclusive":  false,
			"noWait":     false,
		},
	)
	channel.Qos(
		5,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	return &ServiceImpl{channel: channel, queue: q, executor: ex}, nil
}

func (s *ServiceImpl) ListenForTasks() error {
	fmt.Println("Listening for tasks")
	msgs, err := s.channel.Consume(
		queueName,
		"",
		wabbit.Option{
			"auto-ack":  false,
			"exclusive": false,
			"no-local":  false,
			"no-wait":   false,
		},
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			fmt.Println("received msg")
			s.handleMsg(msg)
		}
	}()
	fmt.Println("Awaiting")

	return nil
}

func (s *ServiceImpl) handleMsg(msg wabbit.Delivery) {
	fmt.Println("Handling message")
	go func() {
		taskSpec := new(task.TaskSpec)
		err := taskSpec.UnmarshalBinary(msg.Body())
		if err == nil {
			s.executor.ExecuteTask(taskSpec)
			msg.Ack(false)
		} else {
			log.Printf("error unmarshalling msg %s : %s", string(msg.Body()[:]), err.Error())
		}
	}()
}
