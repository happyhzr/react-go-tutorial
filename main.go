package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Todo struct {
	ID        bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	Completed bool          `json:"completed" bson:"completed"`
	Body      string        `json:"body" bson:"body"`
}

var collection *mongo.Collection

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to MongoDB!")

	collection = client.Database("golang_db").Collection("todos")

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	err = app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalln(err)
	}
}

func getTodos(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close(ctx)
	todos := []*Todo{}
	for cur.Next(ctx) {
		todo := &Todo{}
		err := cur.Decode(todo)
		if err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    todos,
	})
}

func createTodo(c *fiber.Ctx) error {
	todo := &Todo{}
	err := c.BodyParser(todo)
	if err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Body is required",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, todo)
	if err != nil {
		return err
	}
	todo.ID = result.InsertedID.(bson.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    todo,
	})
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid ID",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.D{{"$set", bson.D{{"completed", true}}}}
	_, err = collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid ID",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"_id", objectID}}
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
