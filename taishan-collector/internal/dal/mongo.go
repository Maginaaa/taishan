package dal

import (
	"collector/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	m *mongo.Client
)

func MustInitMongo() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Conf.Mongo.DSN).SetMaxPoolSize(config.Conf.Mongo.PoolSize))
	if err != nil {
		panic(fmt.Errorf("mongo err:%w", err))
	}

	err = m.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongodb initialized")
}

func GetMongo() *mongo.Client {
	return m
}

func MongoDB() string {
	return config.Conf.Mongo.DataBase
}

func GetMongoCollection(table string) *mongo.Collection {
	return m.Database(config.Conf.Mongo.DataBase).Collection(table)
}
