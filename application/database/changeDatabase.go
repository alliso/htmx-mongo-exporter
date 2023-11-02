package database

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"mongoExporter/infrastructure/config"
	"mongoExporter/infrastructure/mongo"
)

func ChangeDatabase(c *fiber.Ctx) error {
	dbName := c.Params("name")

	if dbName == "" {
		log.Fatal("DB not found in remotes")
		return nil
	}

	for _, remote := range config.MainConf.Remotes {
		if remote.Name == dbName {
			mongo.MongoRemote = *mongo.DbManager.GetMongoClient(remote.Uri)
		}
	}

	databases, _ := mongo.DbManager.GetDatabases(mongo.MongoRemote)

	return c.Render("db/index", fiber.Map{
		"SiblingDatabases": databases,
	})
}
