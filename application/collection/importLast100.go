package collection

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"mongoExporter/infrastructure/mongo"
)

func ImportLast100(c *fiber.Ctx) error {
	collectionName := c.Params("name")
	if collectionName == "" {
		log.Fatal("Collection not found")
		return nil
	}

	data, err := mongo.DbManager.FindLast100(collectionName, mongo.MongoRemote)

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
