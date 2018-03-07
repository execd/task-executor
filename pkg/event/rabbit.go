package event

import (
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"log"
)

const queueName = "work_queue"

// RabbitService : interface for building rabbit services
type RabbitService interface {
	GetWorkQueueChan() <-chan wabbit.Delivery
	PublishTaskStatus(info *task.Info) error
}

// RabbitServiceImpl : service to interact with rabbit
type RabbitServiceImpl struct {
	connection          wabbit.Conn
	channel             wabbit.Channel
	taskStatusQueue     wabbit.Queue
	taskStatusQueueName string
	workQueueChan       <-chan wabbit.Delivery
	workQueueName       string
}

// NewRabbitServiceImpl : build a new connection to rabbitmq
func NewRabbitServiceImpl(address string) (*RabbitServiceImpl, error) {

	r := &RabbitServiceImpl{}
	r.initialize(address)

	return r, nil
}

// GetWorkQueueChan : get the work queue channel
func (r *RabbitServiceImpl) GetWorkQueueChan() <-chan wabbit.Delivery {
	return r.workQueueChan
}

// PublishTaskStatus : publish task status on the task status queue
func (r *RabbitServiceImpl) PublishTaskStatus(info *task.Info) error {
	data, err := info.MarshalBinary()
	if err != nil {
		return err
	}
	opts := wabbit.Option{
		"contentType": "application/json",
	}
	return r.channel.Publish("", r.taskStatusQueueName, data, opts)
}

func (r *RabbitServiceImpl) initialize(address string) {
	c := make(chan wabbit.Error)
	fmt.Println("Initializing")
	go func() {
		err := <-c
		log.Println("reconnect: ", err.Error())
		r.initialize(address)
	}()

	conn, err := amqp.Dial(address)
	if err != nil {
		panic("cannot connect")
	}
	conn.NotifyClose(c)

	ch, err := conn.Channel()
	if err != nil {
		panic("cannot create channel")
	}

	r.connection = conn
	r.channel = ch
	r.initializeWorkQueueConsumer()
	r.declareTaskQueue()
}

func (r *RabbitServiceImpl) initializeWorkQueueConsumer() {
	workQueue, err := r.channel.QueueDeclare(
		queueName,
		wabbit.Option{
			"durable":    true,
			"autoDelete": false,
			"exclusive":  false,
			"noWait":     false,
		},
	)
	r.channel.Qos(
		5,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic("Could not setup work_queue")
	}

	workQueueChan, err := r.channel.Consume(
		workQueue.Name(),
		"",
		wabbit.Option{
			"auto-ack":  false,
			"exclusive": false,
			"no-local":  false,
			"no-wait":   false,
		},
	)
	if err != nil {
		panic("Could not setup work_queue consumer")
	}
	r.workQueueChan = workQueueChan
	r.workQueueName = workQueue.Name()
}

func (r *RabbitServiceImpl) declareTaskQueue() {
	name := "task_status_queue"
	taskStatusQueue, err := r.channel.QueueDeclare(
		name,
		wabbit.Option{
			"durable":    true,
			"autoDelete": false,
			"exclusive":  false,
			"noWait":     false,
		},
	)
	if err != nil {
		panic("Could not setup work_queue")
	}
	r.taskStatusQueue = taskStatusQueue
	r.taskStatusQueueName = name
}
