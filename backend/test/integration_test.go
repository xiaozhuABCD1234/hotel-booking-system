// Package test 提供 HTTP 端到端集成测试，使用 Fiber v3 测试助手和真实数据库。
package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"backend/database"
	"backend/middleware"
	"backend/model"
	"backend/router"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// testApp 是共享的 Fiber 应用实例，在 TestMain 中初始化，供所有测试用例复用。
var testApp *fiber.App

// testDB 是共享的数据库连接，在 TestMain 中初始化。
var testDB *gorm.DB

// TestMain 初始化测试环境：连接数据库、创建 Fiber 应用、注册路由，然后运行所有测试。
func TestMain(m *testing.M) {
	// 0. 加载 .env.test 环境变量（优先于 .env 确保测试隔离）
	_ = godotenv.Load("../.env.test")
	_ = godotenv.Load("../.env")

	// 1. 数据库配置：优先使用 TEST_DB_NAME，否则回退到默认配置中的 DB_NAME
	cfg := database.DefaultConfig()
	if testDBName := os.Getenv("TEST_DB_NAME"); testDBName != "" {
		cfg.DBName = testDBName
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("TestMain: failed to connect to database: %v", err)
	}
	testDB = db

	// 2. 创建 Fiber App，注册全局错误处理器
	app := fiber.New(fiber.Config{
		AppName:      "Hotel Booking API (Test)",
		ErrorHandler: middleware.ErrorHandler,
	})

	// 3. 注册所有路由
	router.RegisterRoutes(app, db)

	testApp = app

	// 4. 运行测试
	code := m.Run()

	// 5. 清理：关闭底层 sql.DB 连接
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	os.Exit(code)
}

// parseResponse 读取 HTTP 响应体并解析为统一的 model.Response 格式。
func parseResponse(t *testing.T, resp *http.Response) model.Response {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	resp.Body.Close()

	var result model.Response
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, string(body))
	}
	return result
}

// TestHealthCheck 测试地区省份列表接口，验证基本连通性和 200 响应。
func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/regions/provinces", nil)
	resp, err := testApp.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	result := parseResponse(t, resp)
	if !result.Success {
		t.Fatalf("expected success=true, got false, message: %s", result.Message)
	}
}

// TestGetUsers 测试用户列表分页接口，验证返回 200 且包含分页信息。
func TestGetUsers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&pageSize=10", nil)
	resp, err := testApp.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	result := parseResponse(t, resp)
	if !result.Success {
		t.Fatalf("expected success=true, got false, message: %s", result.Message)
	}

	if result.Pagination == nil {
		t.Fatalf("expected pagination info, got nil")
	}
	if result.Pagination.CurrentPage != 1 {
		t.Fatalf("expected currentPage=1, got %d", result.Pagination.CurrentPage)
	}
	if result.Pagination.ItemsPerPage != 10 {
		t.Fatalf("expected itemsPerPage=10, got %d", result.Pagination.ItemsPerPage)
	}
}

// TestGetHotels 测试酒店列表分页接口，验证返回 200。
func TestGetHotels(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/hotels?page=1&pageSize=10", nil)
	resp, err := testApp.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	result := parseResponse(t, resp)
	if !result.Success {
		t.Fatalf("expected success=true, got false, message: %s", result.Message)
	}
}

// TestGetUserByID_NotFound 测试查询不存在的用户，验证返回 404。
func TestGetUserByID_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/00000000-0000-0000-0000-000000000000", nil)
	resp, err := testApp.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}

	result := parseResponse(t, resp)
	if result.Success {
		t.Fatalf("expected success=false for 404, got true")
	}
}

// TestCreateUser_InvalidBody 测试创建用户时提交空请求体，验证返回 400。
func TestCreateUser_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := testApp.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}

	result := parseResponse(t, resp)
	if result.Success {
		t.Fatalf("expected success=false for 400, got true")
	}
}
