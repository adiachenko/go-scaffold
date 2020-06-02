package rabbitmq

import (
	"net"
	"time"

	"adiachenko/go-scaffold/internal/platform/jaeger"

	"github.com/Vinelab/tracing-go/formats"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	ExchangeTopic  = "topic"
	ExchangeFanout = "fanout"
	ExchangeDirect = "direct"
)

type Subscriber struct {
	ExchangeName string
	ExchangeType string
	QueueName    string
	BindingKeys  []string
	Handler      func(delivery *amqp.Delivery)
}

func Subscribe(subscriber Subscriber) error {
	config := amqp.Config{
		Heartbeat: 120 * time.Second,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 120*time.Second)
		},
	}

	conn, err := amqp.DialConfig(CompileRabbitMqServerURL(), config)
	if err != nil {
		return err
	}
	notify := conn.NotifyClose(make(chan *amqp.Error, 1))
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		subscriber.ExchangeName,
		subscriber.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(
		subscriber.QueueName, // talent.scoring.range
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for _, bindingKey := range subscriber.BindingKeys {
		logrus.Infof("Binding queue \"%s\" to exchange \"%s\" with routing key \"%s\"", queue.Name, subscriber.ExchangeName, bindingKey)

		err = ch.QueueBind(queue.Name, bindingKey, subscriber.ExchangeName, false, nil)
		if err != nil {
			return err
		}
	}

	// Define the max number of unacknowledged deliveries that are permitted on a channel
	if err := ch.Qos(1, 0, false); err != nil {
		logrus.WithError(err).Fatal(err)
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	logrus.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for {
		select {
		case err = <-notify:
			logrus.WithError(err).Panic(err.Error())
		case d := <-msgs:
			spanCtx, err := jaeger.Trace.Extract(&d, formats.AMQP)
			if err != nil {
				logrus.Fatal(err)
			}

			span := jaeger.Trace.StartSpan(queue.Name, spanCtx)

			span.Tag("type", "amqp")
			span.Tag("msg_body", string(d.Body))

			subscriber.Handler(&d)

			jaeger.Trace.RootSpan().Finish()
			jaeger.Trace.Flush()
		default:
			break
		}
	}

	return nil
}

func ProcessSuccess(delivery *amqp.Delivery) {
	if err := delivery.Ack(false); err != nil {
		acknowledgementError(err)
	}
}

func ProcessError(delivery *amqp.Delivery, err error) {
	jaeger.Trace.RootSpan().Tag("error", "true")
	logrus.WithError(err).Error(err.Error())

	// Ensure we don't loop over the same erroneous message indefinitely without some rest period
	time.Sleep(5 * time.Second)

	if err := delivery.Reject(true); err != nil {
		acknowledgementError(err)
	}
}

func acknowledgementError(err error) {
	jaeger.Trace.RootSpan().Tag("error", "true")
	jaeger.Trace.RootSpan().Finish()
	jaeger.Trace.Flush()

	logrus.WithError(err).Fatal(err.Error())
}
