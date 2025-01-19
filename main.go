package main

import (
	"log"
	"os"

	"github.com/Soham-Maha/r2d-be/controllers"
	"github.com/Soham-Maha/r2d-be/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	db.InitDB()

	app := fiber.New()

	//router := mux.NewRouter()

	app.Post("/items", controllers.CreateItem)
	app.Get("/items", controllers.GetItems)
	app.Get("/items/{id}", controllers.GetItem)
	app.Put("/items/{id}", controllers.UpdateItem)
	app.Delete("/items/{id}", controllers.DeleteItem)

	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(app.Listen(":"+port))
}
