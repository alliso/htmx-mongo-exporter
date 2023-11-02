package collection

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"mongoExporter/infrastructure/mongo"
)

func GetCollections(c *fiber.Ctx) error {
	dbName := c.Params("dbname")

	if dbName == "" {
		log.Fatal("DB not found in remotes")
		return nil
	}

	collectionNames, _ := mongo.DbManager.GetCollections(dbName, mongo.MongoRemote)

	collections := make(map[string]string)
	for _, names := range collectionNames {
		collections[names] = dbName
	}

	return c.Render("collection/index", fiber.Map{
		"DbName":      dbName,
		"Collections": collections,
	})
}
