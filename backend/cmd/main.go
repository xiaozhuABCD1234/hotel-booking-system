package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载 .env 文件: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("连接数据库失败(DSN解析/初始化): %v", err)
	}

	// 获取底层 *sql.DB 并 Ping 以验证实际网络连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层 DB 失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping 数据库失败(网络/权限/库不存在): %v", err)
	}

	fmt.Println("数据库连接成功！")
}
