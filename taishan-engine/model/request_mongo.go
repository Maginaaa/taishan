package model

import (
	"context"
	"engine/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	m *mongo.Client
)

func MustInitMongo() {
	var err error

	m, err = mongo.Connect(context.Background(), options.Client().ApplyURI(config.Conf.Mongo.DSN).SetMaxPoolSize(config.Conf.Mongo.PoolSize))
	if err != nil {
		panic(fmt.Errorf("mongo err:%w", err))
	}

	err = m.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongodb initialized")
}

func GetMongo() *mongo.Client {
	return m
}

func GetMongoCollection(table string) *mongo.Collection {
	return m.Database(config.Conf.Mongo.DataBase).Collection(table)
}
