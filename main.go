package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Soham-Maha/r2d-be/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var Collection *mongo.Collection

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectionString := "mongodb+srv://sohammaha15:OCu9cwjwaI5yEHDF@cluster0.0upma.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	Collection = client.Database("Road-to-devsoc-25").Collection("items")

	router := mux.NewRouter()

	router.HandleFunc("/items", createItem).Methods("POST")
	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/items/{id}", getItem).Methods("GET")
	router.HandleFunc("/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item model.Book
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.CreatedAt = time.Now()
	result, err := Collection.InsertOne(context.Background(), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	var items []model.Book
	cursor, err := Collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var item model.Book
		cursor.Decode(&item)
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item model.Book
	err = Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item model.Book
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":        item.Name,
			"description": item.Description,
			"price":       item.Price,
		},
	}

	result, err := Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		update,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
