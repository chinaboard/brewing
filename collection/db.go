package collection

import (
	"context"
	"github.com/chinaboard/brewing/pkg/cfg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newDbClient() (*mongo.Client, error) {
	dsn := cfg.MongoDBConn
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(dsn),
		options.Client().SetMinPoolSize(10),
	)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, err
}
