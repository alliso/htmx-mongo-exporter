package collection

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"mongoExporter/infrastructure/mongo"
)

func ImportFullCollection(c *fiber.Ctx) error {
	collectionName := c.Params("name")
	dbName := c.Params("dbname")
	if collectionName == "" || dbName == "" {
		log.Fatal("Collection not found")
		return nil
	}

	data, err := mongo.DbManager.FindAll(dbName, collectionName, mongo.MongoRemote)

	if err != nil {
		log.Fatal(err)
		return nil
	}
	dbs, err := mongo.DbManager.GetDatabases(mongo.MongoLocal)
	log.Println(dbs)
	err = mongo.DbManager.DeleteOldAndSaveAll(dbName, collectionName, data, mongo.MongoLocal)
	if err != nil {
		log.Fatal(err)
	}

	message, _ := json.Marshal(Toast{ShowMessage: "Imported full Collection"})

	c.Set("HX-TRIGGER", string(message))
	return c.SendStatus(200)
}
