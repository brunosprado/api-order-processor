package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/brunosprado/api-order-processor/domain/client"
	"github.com/brunosprado/api-order-processor/internal/infraestructure/server/http"
	"github.com/brunosprado/api-order-processor/internal/infraestructure/server/http/mongodb"
	"github.com/brunosprado/api-order-processor/internal/infraestructure/server/rabbitmq"
	"github.com/brunosprado/api-order-processor/pkg/log"
	"github.com/joho/godotenv"
)

const (
	envApplicationPort = "API_PORT"
	defaultPort        = "8081"
)

var (
	version, build, date string
)

func main() {
	_ = godotenv.Load()
	logger := log.NewZeroLog("api-order-processor", version, log.Level("info"))

	logger.Info().Sendf("API Order Processor starting - version:%s, build:%s, date:%s", version, build, date)

	db, err := mongodb.NewConnection("mongodb://root:example@mongodb:27017", "orders_db", logger)
	if err != nil {
		logger.Fatal().Err(err).Sendf("error connecting mongoDB: %q", err)
		return
	}

	clientStorage := mongodb.NewClientStorage(db, logger)

	rabbitPublisher, err := rabbitmq.NewPublisher("amqp://user:password@rabbitmq:5672/", "order_events")
	if err != nil {
		logger.Fatal().Err(err).Send("Failed to create RabbitMQ publisher")
		return
	}

	// Services
	clientService := client.NewService(clientStorage, rabbitPublisher)

	// Handler
	handler := http.NewHandler(clientService, logger)

	// Server
	server := http.New(defaultPort, handler, logger)
	server.ListenAndServe()

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	server.Shutdown()

}
