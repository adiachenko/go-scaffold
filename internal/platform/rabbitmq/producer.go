package rabbitmq

import (
	"adiachenko/go-scaffold/internal/platform/jaeger"

	"github.com/Vinelab/tracing-go/formats"
	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type Message struct {
	Body         []byte
	RoutingKey   string
	Exchange     string
	ExchangeType string
}

func TopicMessage(body []byte, exchange, routingKey string) Message {
	return Message{
		Body:         body,
		Exchange:     exchange,
		RoutingKey:   routingKey,
		ExchangeType: "topic",
	}
}

func FanoutMessage(body []byte, exchange string) Message {
	return Message{
		Body:         body,
		Exchange:     exchange,
		RoutingKey:   "",
		ExchangeType: "fanout",
	}
}

func DirectMessage(body []byte, exchange, routingKey string) Message {
	return Message{
		Body:         body,
		Exchange:     exchange,
		RoutingKey:   routingKey,
		ExchangeType: "direct",
	}
}

func Produce(produced Message) error {
	conn, err := amqp.Dial(CompileRabbitMqServerURL())
	if err != nil {
		return err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		produced.Exchange,     // name
		produced.ExchangeType, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)

	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         produced.Body,
		DeliveryMode: amqp.Persistent,
	}

	if err := jaeger.Trace.Inject(&msg, formats.AMQP); err != nil {
		return err
	}

	err = ch.Publish(
		produced.Exchange,   // exchange
		produced.RoutingKey, // routing key
		false,               // mandatory
		false,               // immediate
		msg,
	)

	if err != nil {
		return err
	}

	logrus.Infof(
		" [x] Sent an AMQP message on %s exchange with a routing key of %s and trace uuid of %s",
		produced.Exchange,
		produced.RoutingKey,
		jaeger.Trace.UUID(),
	)

	return nil
}
