package main

import (
	"log"
	"net/http"

	"github.com/Soham-Maha/r2d-be/controllers"
	"github.com/Soham-Maha/r2d-be/db"
	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()

	router.HandleFunc("/items", controllers.CreateItem).Methods("POST")
	router.HandleFunc("/items", controllers.GetItems).Methods("GET")
	router.HandleFunc("/items/{id}", controllers.GetItem).Methods("GET")
	router.HandleFunc("/items/{id}", controllers.UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", controllers.DeleteItem).Methods("DELETE")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
