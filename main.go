package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"log"
	"mongoExporter/infrastructure/config"
	"mongoExporter/infrastructure/mongo"
)

type Toast struct {
	ShowMessage string `json:"showMessage"`
}

func main() {
	config.LoadConfig()
	mongo.InitMongoLocal()
	initServer()
}

func initServer() {
	app := fiber.New(fiber.Config{
		Views: html.New("./public", ".html"),
	})

	app.Use(logger.New())

	app.Get("/", renderIndex)
	app.Get("/database/:name", changeDatabase)
	app.Get("/collection/:name", getCollections)
	app.Post("/collection/:name", importLast100)
	app.Post("/collection/:name/full", importFullCollection)

	log.Fatal(app.Listen(":4000"))
}

func renderIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Remotes": config.MainConf.Remotes,
	})
}

func changeDatabase(c *fiber.Ctx) error {
	dbName := c.Params("name")

	if dbName == "" {
		log.Fatal("DB not found in remotes")
		return nil
	}

	for _, remote := range config.MainConf.Remotes {
		if remote.Name == dbName {
			mongo.DbManager.SetCurrentDb(remote.Name)
			mongo.MongoRemote = *mongo.DbManager.GetMongoClient(remote.Uri)
		}
	}

	databases, _ := mongo.DbManager.GetDatabases(mongo.MongoRemote)

	return c.Render("db/index", fiber.Map{
		"DbName":           mongo.DbManager.CurrentDb,
		"SiblingDatabases": databases,
	})
}

func getCollections(c *fiber.Ctx) error {
	dbName := c.Params("name")

	if dbName == "" {
		log.Fatal("DB not found in remotes")
		return nil
	}

	mongo.DbManager.SetCurrentDb(dbName)

	collections, _ := mongo.DbManager.GetCollections(mongo.MongoRemote)

	return c.Render("collection/index", fiber.Map{
		"Collections": collections,
	})
}

func getFullCollection(collectionName string) ([]any, error) {
	return mongo.DbManager.FindAll(collectionName, mongo.MongoRemote)
}

func getLast100(collectionName string) ([]any, error) {
	return mongo.DbManager.FindLast100(collectionName, mongo.MongoRemote)
}

func importFullCollection(c *fiber.Ctx) error {
	collectionName := c.Params("name")
	if collectionName == "" {
		log.Fatal("Collection not found")
		return nil
	}

	data, err := getFullCollection(collectionName)

	if err != nil {
		log.Fatal(err)
		return nil
	}
	dbs, err := mongo.DbManager.GetDatabases(mongo.MongoLocal)
	log.Println(dbs)
	err = mongo.DbManager.DeleteOldAndSaveAll(collectionName, data, mongo.MongoLocal)
	if err != nil {
		log.Fatal(err)
	}

	message, _ := json.Marshal(Toast{ShowMessage: "Imported full Collection"})

	c.Set("HX-TRIGGER", string(message))
	return c.SendStatus(200)
}

func importLast100(c *fiber.Ctx) error {
	collectionName := c.Params("name")
	if collectionName == "" {
		log.Fatal("Collection not found")
		return nil
	}

	data, err := getLast100(collectionName)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = mongo.DbManager.DeleteOldAndSaveAll(collectionName, data, mongo.MongoLocal)

	if err != nil {
		log.Fatal(err)
	}

	message, _ := json.Marshal(Toast{ShowMessage: "Imported last 100"})

	c.Set("HX-TRIGGER", string(message))
	return c.SendString("100 last")
}
