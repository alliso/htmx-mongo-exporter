package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	mongoClient *mongo.Client
	currentDb   string
}

type DbManager interface {
	SetMongoClient(mongoUri string)
	SetCurrentDb(dbName string)
	GetDatabases() ([]string, error)
	GetCollections(dbName string) ([]string, error)
	FindAll(collectionName string) ([]any, error)
	FindLast100(collectionName string) ([]any, error)
	DeleteOldAndSaveAll(collectionName string, data []any) error
}

func (m Db) SetMongoClient(mongoUri string) {
	m.mongoClient, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
}

func (m Db) SetCurrentDb(dbName string) {
	m.currentDb = dbName
}

func (m Db) GetDatabases() ([]string, error) {
	return m.mongoClient.ListDatabaseNames(context.TODO(), bson.D{{}})
}

func (m Db) GetCollections(dbName string) ([]string, error) {
	return m.mongoClient.Database(dbName).ListCollectionNames(context.TODO(), bson.D{{}})
}

func (m Db) FindAll(collectionName string) ([]any, error) {
	cursor, err := m.mongoClient.Database(m.currentDb).Collection(collectionName).Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var data []any
	err = cursor.All(context.TODO(), &data)

	return data, err
}

func (m Db) FindLast100(collectionName string) ([]any, error) {
	optsFind := options.Find().SetLimit(100).SetSort(bson.D{{"_id", -1}})
	cursor, err := m.mongoClient.Database(m.currentDb).Collection(collectionName).Find(context.TODO(), bson.D{{}}, optsFind)

	if err != nil {
		return nil, err
	}

	var data []any
	err = cursor.All(context.TODO(), &data)

	return data, err
}

func (m Db) DeleteOldAndSaveAll(collectionName string, data []any) error {
	err := m.mongoClient.Database(m.currentDb).Collection(collectionName).Drop(context.TODO())

	if err != nil {
		return err
	}

	_, err = m.mongoClient.Database(m.currentDb).Collection(collectionName).InsertMany(context.TODO(), data)

	return err
}
