package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brunosprado/api-order-processor/internal/infraestructure/server/http/mongodb"
	"github.com/streadway/amqp"
)

const (
	rabbitURL   = "amqp://user:password@rabbitmq:5672/"
	rabbitQueue = "order_events"
	mongoURL    = "mongodb://root:example@mongodb:27017"
	dbName      = "orders_db"
)

type orderEvent struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func main() {
	// MongoDB
	db, err := mongodb.NewConnection(mongoURL, dbName, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	clientStorage := mongodb.NewClientStorage(db, nil)

	// RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		rabbitQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		rabbitQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Worker started, waiting for messages...")

	go func() {
		for d := range msgs {
			var evt orderEvent
			if err := json.Unmarshal(d.Body, &evt); err != nil {
				log.Printf("Invalid message: %v", err)
				continue
			}
			log.Printf("Received order_id=%s, status=%s", evt.OrderID, evt.Status)
			time.Sleep(2 * time.Second)
			if err := clientStorage.UpdateOrderStatus(evt.OrderID, "PROCESSADO"); err != nil {
				log.Printf("Failed to update order: %v", err)
			} else {
				log.Printf("Order %s processed!", evt.OrderID)
				order, err := clientStorage.GetOrderById(evt.OrderID)
				if err != nil {
					log.Printf("Failed to fetch order from MongoDB: %v", err)
				} else {
					orderJson, _ := json.MarshalIndent(order, "", "  ")
					log.Printf("Order from MongoDB: %s", orderJson)
				}
			}
		}
	}()

	<-quit
	log.Println("Worker shutting down...")
}
