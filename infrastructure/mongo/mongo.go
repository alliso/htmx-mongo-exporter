package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mongoExporter/infrastructure/config"
)

var MongoRemote mongo.Client
var MongoLocal mongo.Client

var DbManager Db

type Db struct {
}

type IDbManager interface {
	GetMongoClient(mongoUri string) *mongo.Client
	GetDatabases(client mongo.Client) ([]string, error)
	GetCollections(dbName string, client mongo.Client) ([]string, error)
	FindAll(dbName string, collectionName string, client mongo.Client) ([]any, error)
	FindLast100(dbName string, collectionName string, client mongo.Client) ([]any, error)
	DeleteOldAndSaveAll(dbName string, collectionName string, data []any, client mongo.Client) error
}

func InitMongoLocal() {
	MongoLocal = *DbManager.GetMongoClient(config.MainConf.Local)
}
func (m Db) GetMongoClient(mongoUri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	return client
}

func (m Db) GetDatabases(client mongo.Client) ([]string, error) {
	return client.ListDatabaseNames(context.TODO(), bson.D{{}})
}

func (m Db) GetCollections(dbName string, client mongo.Client) ([]string, error) {
	return client.Database(dbName).ListCollectionNames(context.TODO(), bson.D{{}})
}

func (m Db) FindAll(dbName string, collectionName string, client mongo.Client) ([]any, error) {
	cursor, err := client.Database(dbName).Collection(collectionName).Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var data []any
	err = cursor.All(context.TODO(), &data)

	return data, err
}

func (m Db) FindLast100(dbName string, collectionName string, client mongo.Client) ([]any, error) {
	optsFind := options.Find().SetLimit(100).SetSort(bson.D{{"_id", -1}})
	cursor, err := client.Database(dbName).Collection(collectionName).Find(context.TODO(), bson.D{{}}, optsFind)

	if err != nil {
		return nil, err
	}

	var data []any
	err = cursor.All(context.TODO(), &data)

	return data, err
}

func (m Db) DeleteOldAndSaveAll(dbName string, collectionName string, data []any, client mongo.Client) error {
	err := client.Database(dbName).Collection(collectionName).Drop(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Database(dbName).Collection(collectionName).InsertMany(context.TODO(), data)

	return err
}
