package client

import (
	domainPkg "github.com/brunosprado/api-order-processor/domain"
	"github.com/brunosprado/api-order-processor/internal/infraestructure/server/rabbitmq"
)

type service struct {
	clientStorage domainPkg.ClientStorage
	publisher     *rabbitmq.Publisher
}

func NewService(clientStorage domainPkg.ClientStorage, publisher *rabbitmq.Publisher) *service {
	return &service{
		clientStorage: clientStorage,
		publisher:     publisher,
	}
}

func (s *service) PostOrder(order domainPkg.Order) error {
	err := s.clientStorage.PersistOrder(order)
	if err != nil {
		return err
	}
	// Publica evento na fila RabbitMQ
	if s.publisher != nil {
		if err := s.publisher.PublishOrderEvent(order, "PROCESSANDO"); err != nil {
			return err
		}
	}
	return nil
}
