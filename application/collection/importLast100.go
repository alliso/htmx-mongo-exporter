package collection

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"mongoExporter/infrastructure/mongo"
)

func ImportLast100(c *fiber.Ctx) error {
	collectionName := c.Params("name")
	dbName := c.Params("dbname")
	if collectionName == "" || dbName == "" {
		log.Fatal("Collection or db not found")
		return nil
	}

	data, err := mongo.DbManager.FindLast100(dbName, collectionName, mongo.MongoRemote)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = mongo.DbManager.DeleteOldAndSaveAll(dbName, collectionName, data, mongo.MongoLocal)

	if err != nil {
		log.Fatal(err)
	}

	message, _ := json.Marshal(Toast{ShowMessage: "Imported last 100"})

	c.Set("HX-TRIGGER", string(message))
	return c.SendString("100 last")
}
