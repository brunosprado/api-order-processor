package mongodb

import (
	"context"
	"time"

	"github.com/brunosprado/api-order-processor/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database = ""
)

func NewConnection(mongoURL string, dbName string, log log.Logger) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	//
	database = dbName

	if err != nil {
		log.Error().Sendf("error on MongoDB connection: %q", err)
		return nil, err
	}

	return client, nil
}
