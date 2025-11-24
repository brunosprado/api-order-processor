package rabbitmq

import (
	"encoding/json"

	"github.com/brunosprado/api-order-processor/domain"
	"github.com/streadway/amqp"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewPublisher(amqpURL, queue string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	return &Publisher{conn: conn, channel: ch, queue: queue}, nil
}

func (p *Publisher) PublishOrderEvent(order domain.Order, status string) error {
	msg := map[string]interface{}{
		"order_id": order.OrderID,
		"status":   status,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.channel.Publish(
		"",
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) Close() {
	p.channel.Close()
	p.conn.Close()
}
