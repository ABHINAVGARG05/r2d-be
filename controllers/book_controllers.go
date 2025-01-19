package controllers

import (
	"context"
	"time"

	"github.com/Soham-Maha/r2d-be/db"
	"github.com/Soham-Maha/r2d-be/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateItem(c* fiber.Ctx) error {
	var book model.Book

	if err:= c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"invalid request Body",
		})
	}

	book.CreatedAt = time.Now()

	result, err := db.Collection.InsertOne(context.Background(), book)
	
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	book.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusOK).JSON(book)
}

func GetItems(c* fiber.Ctx) error {
	var books []model.Book
	cursor, err := db.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var book model.Book
		cursor.Decode(&book)
		books = append(books, book)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data":books,
	})
}

func GetItem(c* fiber.Ctx) error {

	bookIdParams := c.Params("bookId")

	id, err := primitive.ObjectIDFromHex(bookIdParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	var book model.Book

	err = db.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&book)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":err.Error(),
				"Message":"Book not found",
			})
		}
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"error":err.Error(),
			"Message":"Internal server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data":book,
	})
}

func UpdateItem(c* fiber.Ctx) error {

	idParams := c.Params("id")

	id, err := primitive.ObjectIDFromHex(idParams)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	var book model.Book

	if err:= c.BodyParser(&book); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	update := bson.M{
		"$set": bson.M{
			"name":        book.Name,
			"description": book.Description,
			"price":       book.Price,
		},
	}

	result, err := db.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		update,
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	if result.ModifiedCount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message":"Book not found",
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"Message":"Updated Item Details",
	})
}

func DeleteItem(c* fiber.Ctx) error {

	idParams:= c.Params("id")

	id, err := primitive.ObjectIDFromHex(idParams)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	result, err := db.Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message":"Book not found",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
