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
	CurrentDb string
}

type IDbManager interface {
	GetMongoClient(mongoUri string) *mongo.Client
	SetCurrentDb(dbName string)
	GetDatabases(client mongo.Client) ([]string, error)
	GetCollections(client mongo.Client) ([]string, error)
	FindAll(collectionName string, client mongo.Client) ([]any, error)
	FindLast100(collectionName string, client mongo.Client) ([]any, error)
	DeleteOldAndSaveAll(collectionName string, data []any, client mongo.Client) error
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

func (m *Db) SetCurrentDb(dbName string) {
	m.CurrentDb = dbName
}

func (m Db) GetDatabases(client mongo.Client) ([]string, error) {
	return client.ListDatabaseNames(context.TODO(), bson.D{{}})
}

func (m Db) GetCollections(client mongo.Client) ([]string, error) {
	return client.Database(m.CurrentDb).ListCollectionNames(context.TODO(), bson.D{{}})
}

func (m Db) FindAll(collectionName string, client mongo.Client) ([]any, error) {
	cursor, err := client.Database(m.CurrentDb).Collection(collectionName).Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var data []any
	err = cursor.All(context.TODO(), &data)

	return data, err
}

func (m Db) FindLast100(collectionName string, client mongo.Client) ([]any, error) {
	optsFind := options.Find().SetLimit(100).SetSort(bson.D{{"_id", -1}})
	cursor, err := client.Database(m.CurrentDb).Collection(collectionName).Find(context.TODO(), bson.D{{}}, optsFind)

	if err != nil {
		return nil, err
	}

	var data []any
	err = cursor.All(context.TODO(), &data)

	return data, err
}

func (m Db) DeleteOldAndSaveAll(collectionName string, data []any, client mongo.Client) error {
	err := client.Database(m.CurrentDb).Collection(collectionName).Drop(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Database(m.CurrentDb).Collection(collectionName).InsertMany(context.TODO(), data)

	return err
}
