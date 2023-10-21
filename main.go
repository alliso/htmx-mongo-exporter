package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type remote struct {
	Name string
	Uri  string
}

type Data struct {
	Config struct {
		Local   string
		Remotes []remote
	}
}

type Database struct {
	Name string
}

var config = &Data{}

var currentDb remote = remote{}

var clientGlobal mongo.Client

var mongoLocal mongo.Client

var globalDb string

func main() {
	loadConfig()
	initMongoLocal()
	initServer()
}

func loadConfig() {
	buf, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = yaml.Unmarshal(buf, config)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(config.Config.Remotes)
}

func initMongoLocal() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Config.Local))

	if err != nil {
		panic(err)
	}

	mongoLocal = *client
}

func initServer() {
	app := fiber.New(fiber.Config{
		Views: html.New("./public", ".html"),
	})

	app.Use(logger.New())

	app.Get("/", renderIndex)
	app.Get("/database/:name", changeDatabase)
	app.Get("/collection/:name", getCollections)
	app.Post("/collection/:name", importCollection)
	app.Post("/collection/:name/full", importFullCollection)

	log.Fatal(app.Listen(":4000"))
}

func renderIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Remotes": config.Config.Remotes,
	})
}

func changeDatabase(c *fiber.Ctx) error {
	dbName := c.Params("name")

	if dbName == "" {
		log.Fatal("DB not found in remotes")
		return nil
	}

	for _, remote := range config.Config.Remotes {
		if remote.Name == dbName {
			currentDb = remote
		}
	}

	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(currentDb.Uri))
	clientGlobal = *client

	filter := bson.D{{}}

	databases, _ := clientGlobal.ListDatabases(context.TODO(), filter)

	return c.Render("db/index", fiber.Map{
		"DbName":           currentDb.Name,
		"SiblingDatabases": databases.Databases,
	})
}

func getCollections(c *fiber.Ctx) error {
	dbName := c.Params("name")

	if dbName == "" {
		log.Fatal("DB not found in remotes")
		return nil
	}

	globalDb = dbName

	collections, _ := clientGlobal.Database(dbName).ListCollectionNames(context.TODO(), bson.D{{}})

	return c.Render("collection/index", fiber.Map{
		"Collections": collections,
	})
}

func importFullCollection(c *fiber.Ctx) error {
	collectionName := c.Params("name")
	if collectionName == "" {
		log.Fatal("Collection not found")
		return nil
	}

	cursor, _ := clientGlobal.Database(globalDb).Collection(collectionName).Find(context.TODO(), bson.D{{}})

	var data []any
	err := cursor.All(context.TODO(), &data)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	createOrDropCollection(globalDb, collectionName)
	_, err = mongoLocal.Database(globalDb).Collection(collectionName).InsertMany(context.TODO(), data)

	if err != nil {
		log.Fatal(err)
	}

	return c.SendString("FULL")
}

func createOrDropCollection(dbName string, collectionName string) {
	err := mongoLocal.Database(dbName).Collection(collectionName).Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func importCollection(c *fiber.Ctx) error {
	fmt.Println("IMPORTED")
	return c.SendString("IMPORTED")
}
