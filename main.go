package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"log"
	"mongoExporter/application/collection"
	"mongoExporter/application/database"
	"mongoExporter/infrastructure/config"
	"mongoExporter/infrastructure/mongo"
)

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
	app.Get("/database/:name", database.ChangeDatabase)
	app.Get("/collection/:name", collection.GetCollections)
	app.Post("/collection/:name", collection.ImportLast100)
	app.Post("/collection/:name/full", collection.ImportFullCollection)

	log.Fatal(app.Listen(":4000"))
}

func renderIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Remotes": config.MainConf.Remotes,
	})
}
