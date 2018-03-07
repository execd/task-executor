package event

import (
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
	"log"
)

const queueName = "work_queue"

type RabbitService interface {
	GetChannel() wabbit.Channel
	GetWorkQueueChan() <-chan wabbit.Delivery
}

type RabbitServiceImpl struct {
	connection    wabbit.Conn
	channel       wabbit.Channel
	workQueueChan <-chan wabbit.Delivery
	workQueueName string
}

// NewRabbitConnection : build a new connection to rabbitmq
func NewRabbitServiceImpl(address string) (*RabbitServiceImpl, error) {

	r := &RabbitServiceImpl{}
	r.initialize(address)

	return r, nil
}

func (r *RabbitServiceImpl) GetChannel() wabbit.Channel {
	return r.channel
}

func (r *RabbitServiceImpl) GetWorkQueueChan() <-chan wabbit.Delivery {
	return r.workQueueChan
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

	q, err := ch.QueueDeclare(
		queueName,
		wabbit.Option{
			"durable":    true,
			"autoDelete": false,
			"exclusive":  false,
			"noWait":     false,
		},
	)
	ch.Qos(
		5,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	msgs, err := ch.Consume(
		q.Name(),
		"",
		wabbit.Option{
			"auto-ack":  false,
			"exclusive": false,
			"no-local":  false,
			"no-wait":   false,
		},
	)
	if err != nil {
		panic("Could not consumer work_queue")
	}

	r.connection = conn
	r.channel = ch
	r.workQueueChan = msgs
	r.workQueueName = queueName
}
