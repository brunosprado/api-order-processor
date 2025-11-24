package mongodb

import (
	"context"

	"github.com/brunosprado/api-order-processor/domain"
	"github.com/brunosprado/api-order-processor/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type clientStorage struct {
	clientMongo *mongo.Client
	log         log.Logger
}

func NewClientStorage(clientMongo *mongo.Client, log log.Logger) *clientStorage {
	return &clientStorage{
		clientMongo: clientMongo,
		log:         log,
	}
}

func (c *clientStorage) PersistOrder(order domain.Order) error {
	collection := c.clientMongo.Database("orders_db").Collection("orders")
	_, err := collection.InsertOne(context.Background(), order)
	if err != nil {
		c.log.Error().Sendf("Erro ao persistir pedido: %v", err)
		return err
	}
	return nil
}

func (c *clientStorage) UpdateOrderStatus(orderID, status string) error {
	collection := c.clientMongo.Database("orders_db").Collection("orders")
	filter := bson.M{"_id": orderID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.log.Error().Sendf("Erro ao atualizar status do pedido: %v", err)
		return err
	}
	return nil
}

func (c *clientStorage) GetOrderById(orderID string) (*domain.Order, error) {
	collection := c.clientMongo.Database("orders_db").Collection("orders")
	filter := bson.M{"_id": orderID}
	var order domain.Order
	err := collection.FindOne(context.Background(), filter).Decode(&order)
	if err != nil {
		c.log.Error().Sendf("Erro ao buscar pedido por ID: %v", err)
		return nil, err
	}
	return &order, nil
}
