package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Soham-Maha/r2d-be/controllers"
	"github.com/Soham-Maha/r2d-be/db"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()

	router.HandleFunc("/items", controllers.CreateItem).Methods("POST")
	router.HandleFunc("/items", controllers.GetItems).Methods("GET")
	router.HandleFunc("/items/{id}", controllers.GetItem).Methods("GET")
	router.HandleFunc("/items/{id}", controllers.UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", controllers.DeleteItem).Methods("DELETE")

	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
