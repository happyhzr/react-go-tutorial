package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	todos := []*Todo{}

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": todos,
		})
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
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
		todo.ID = len(todos) + 1
		todos = append(todos, todo)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    todo,
		})
	})

	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}
		for _, todo := range todos {
			if todo.ID == id {
				todo.Completed = !todo.Completed
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"success": true,
					"data":    todo,
				})
			}
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Todo not found",
		})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}
		for i, todo := range todos {
			if todo.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"success": true,
				})
			}
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Todo not found",
		})
	})

	err = app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalln(err)
	}
}
