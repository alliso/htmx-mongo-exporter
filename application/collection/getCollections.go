package collection

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"mongoExporter/infrastructure/mongo"
)

func GetCollections(c *fiber.Ctx) error {
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
