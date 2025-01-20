package handler

import (
	"log"
	"net/http"

	"github.com/Soham-Maha/r2d-be/controllers"
	"github.com/Soham-Maha/r2d-be/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/joho/godotenv"
)

// Handler is the entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.String()
	app := setupApp()
	adaptor.FiberApp(app).ServeHTTP(w, r)
}

// setupApp initializes the Fiber application with routes
func setupApp() *fiber.App {
	// Initialize database
	db.InitDB()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Create a Fiber app
	app := fiber.New()

	// Define routes
	app.Post("/items", controllers.CreateItem)
	app.Get("/items", controllers.GetItems)
	app.Get("/items/:id", controllers.GetItem)
	app.Put("/items/:id", controllers.UpdateItem)
	app.Delete("/items/:id", controllers.DeleteItem)

	// Version and health check endpoints
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Welcome to Fiber on Vercel!",
		})
	})
	app.Get("/v1", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"version": "v1",
		})
	})
	app.Get("/v2", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"version": "v2",
		})
	})

	return app
}
