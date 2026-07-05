package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
)

type CustomCtx struct {
	fiber.DefaultCtx
}

func (c *CustomCtx) CustomMethod() string {
	return "custom value"
}

func main() {
	app := fiber.NewWithCustomCtx(func(app *fiber.App) fiber.CustomCtx {
		return &CustomCtx{
			DefaultCtx: *fiber.NewDefaultCtx(app),
		}
	})

	app.Get("/", func(c fiber.Ctx) error {
		customCtx := c.(*CustomCtx)
		return c.SendString(customCtx.CustomMethod())
	})

	log.Fatal(app.Listen(":3000"))
}
