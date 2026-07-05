// Package main 启动 Fiber v3 HTTP 服务，连接 PostgreSQL，注册路由与全局错误处理。
package main

import (
	"log"

	"backend/auth"
	"backend/database"
	"backend/middleware"
	"backend/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	if err := auth.LoadSecret(); err != nil {
		log.Fatalf("JWT configuration failed: %v", err)
	}
	db, err := database.Connect(database.DefaultConfig())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Database connected")

	app := fiber.New(fiber.Config{
		AppName:      "Hotel Booking API",
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(recoverer.New())
	app.Use(logger.New())

	router.RegisterRoutes(app, db)

	log.Fatal(app.Listen(":3000"))
}
