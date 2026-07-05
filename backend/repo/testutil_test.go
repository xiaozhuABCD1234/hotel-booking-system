package repo

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// testDB is the shared database connection for all repository integration tests.
// Initialized once in TestMain, reused across all tests.
var testDB *gorm.DB

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		envOrDefault("DB_HOST", "106.14.58.94"),
		envOrDefault("DB_PORT", "5432"),
		envOrDefault("DB_USER", "root"),
		envOrDefault("DB_PASSWORD", "ABCD1234"),
		envOrDefault("DB_NAME", "dev"),
		envOrDefault("DB_SSLMODE", "disable"),
		envOrDefault("DB_TIMEZONE", "Asia/Shanghai"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)

	testDB = db
	code := m.Run()

	sqlDB.Close()
	os.Exit(code)
}

// txRepo creates a transaction-scoped DB connection and returns it.
// The transaction is rolled back in t.Cleanup, so test data never persists.
func txRepo(t *testing.T) *gorm.DB {
	t.Helper()
	tx := testDB.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		tx.Rollback()
	})
	return tx
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
