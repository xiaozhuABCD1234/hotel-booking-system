// Package main 启动 Fiber v3 HTTP 服务，连接 PostgreSQL，注册路由与全局错误处理。
//
// @title                      Hotel Booking API
// @version                    1.0
// @description                酒店预订管理系统 API — 支持用户注册、酒店搜索、客房预订、订单管理、评价与数据统计
// @contact.name               上海电力大学
// @contact.url                https://www.shiep.edu.cn
// @license.name               MIT
// @host                       localhost:3000
// @BasePath                   /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                JWT Bearer Token，格式 "Bearer <token>"
package main

import (
	"log"
 	"time"

	"backend/auth"
	"backend/database"
	"backend/middleware"
	"backend/router"

	_ "backend/docs"
	_ "github.com/joho/godotenv/autoload"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"gorm.io/gorm"
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

	go cleanupJWTBlacklist(db)

	app := fiber.New(fiber.Config{
		AppName:      "Hotel Booking API",
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(recoverer.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
	}))

	router.RegisterRoutes(app, db)

	log.Fatal(app.Listen(":3000"))
}

// cleanupJWTBlacklist 每天凌晨 4:00 清理 jwt_blacklist_1718 中已过期的记录
func cleanupJWTBlacklist(db *gorm.DB) {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}

	d := time.NewTimer(next.Sub(now))
	<-d.C

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		result := db.Exec("DELETE FROM jwt_blacklist_1718 WHERE expires_at < NOW()")
		if result.Error != nil {
			log.Printf("cleanup JWT blacklist failed: %v", result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("cleaned %d expired JWT blacklist entries", result.RowsAffected)
		}
	}
}
