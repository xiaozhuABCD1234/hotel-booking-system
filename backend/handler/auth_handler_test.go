package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"backend/auth"
	"backend/handler"
	schema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ─── Mock Repository ───────────────────────────────────────────

type mockUserRepo struct {
	findByUsernameFunc func(ctx context.Context, username string) (*schema.User, error)
	createFunc         func(ctx context.Context, user *schema.User) error
}

var _ repo.UserRepository = (*mockUserRepo)(nil)

func (m *mockUserRepo) Create(ctx context.Context, user *schema.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}
func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*schema.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*schema.User, error) {
	if m.findByUsernameFunc != nil {
		return m.findByUsernameFunc(ctx, username)
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockUserRepo) FindByPhone(ctx context.Context, phone string) (*schema.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockUserRepo) FindAll(ctx context.Context, offset, limit int, role *schema.UserRole) ([]schema.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Update(ctx context.Context, user *schema.User) error { return nil }
func (m *mockUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	return nil
}
func (m *mockUserRepo) UpdatePoints(ctx context.Context, userID uuid.UUID, points int32) error {
	return nil
}
func (m *mockUserRepo) UpdateVipLevel(ctx context.Context, userID uuid.UUID, vipLevel int16) error {
	return nil
}
func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

// ─── Test Helpers ──────────────────────────────────────────────

func setupAuthApp(t *testing.T, repo repo.UserRepository) *fiber.App {
	t.Helper()
	if err := auth.LoadSecret(); err != nil {
		t.Fatalf("setupAuthApp: failed to load JWT secret: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := http.StatusInternalServerError
			if fe, ok := err.(*fiber.Error); ok {
				code = fe.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	api := app.Group("/api/v1")
	authGroup := api.Group("/auth")
	h := handler.NewAuthHandler(repo)
	authGroup.Post("/register", h.Register)
	authGroup.Post("/login", h.Login)
	authGroup.Post("/refresh", h.Refresh)
	authGroup.Post("/logout", h.Logout)

	return app
}

func postJSON(t *testing.T, app *fiber.App, path, body string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

// assertTokenBody 校验响应体为成功的令牌响应，且 accessToken 非空、tokenType 为 Bearer。
func assertTokenBody(t *testing.T, resp *http.Response) {
	t.Helper()
	var got struct {
		Success bool `json:"success"`
		Data    struct {
			AccessToken string `json:"accessToken"`
			TokenType   string `json:"tokenType"`
			ExpiresIn   int    `json:"expiresIn"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if !got.Success {
		t.Errorf("expected success=true, got false")
	}
	if got.Data.AccessToken == "" {
		t.Errorf("expected non-empty accessToken, got empty")
	}
	if got.Data.TokenType != "Bearer" {
		t.Errorf("expected tokenType=Bearer, got %q", got.Data.TokenType)
	}
	if got.Data.ExpiresIn <= 0 {
		t.Errorf("expected positive expiresIn, got %d", got.Data.ExpiresIn)
	}
}

var errDB = errors.New("database connection refused")

// ─── Register Tests ────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-register-success")
	mock := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*schema.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createFunc: func(ctx context.Context, user *schema.User) error {
			user.ID = uuid.New()
			return nil
		},
	}
	app := setupAuthApp(t, mock)

	resp := postJSON(t, app, "/api/v1/auth/register", `{"username":"testuser","password":"secret123"}`)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	assertTokenBody(t, resp)
}

func TestRegister_MissingFields(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-register-missing")
	app := setupAuthApp(t, &mockUserRepo{})

	resp := postJSON(t, app, "/api/v1/auth/register", `{"username":"","password":""}`)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-register-short")
	app := setupAuthApp(t, &mockUserRepo{})

	resp := postJSON(t, app, "/api/v1/auth/register", `{"username":"testuser","password":"abc"}`)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for short password, got %d", resp.StatusCode)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-register-dup")
	existingUser := &schema.User{ID: uuid.New(), Username: "existing"}
	mock := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*schema.User, error) {
			if username == "existing" {
				return existingUser, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
	}
	app := setupAuthApp(t, mock)

	resp := postJSON(t, app, "/api/v1/auth/register", `{"username":"existing","password":"secret123"}`)

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", resp.StatusCode)
	}
}

func TestRegister_DatabaseError(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-register-dberr")
	mock := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*schema.User, error) {
			return nil, errDB
		},
	}
	app := setupAuthApp(t, mock)

	resp := postJSON(t, app, "/api/v1/auth/register", `{"username":"testuser","password":"secret123"}`)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

// ─── Login Tests ───────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-login-success")
	password := "secret123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userID := uuid.New()

	mock := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*schema.User, error) {
			return &schema.User{
				ID:       userID,
				Username: "testuser",
				Password: string(hashed),
				Role:     schema.RoleCustomer,
				Status:   1,
			}, nil
		},
	}
	app := setupAuthApp(t, mock)

	resp := postJSON(t, app, "/api/v1/auth/login", `{"username":"testuser","password":"secret123"}`)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	assertTokenBody(t, resp)
}

func TestLogin_WrongPassword(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-login-wrong")
	password := "correct"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mock := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*schema.User, error) {
			return &schema.User{
				ID:       uuid.New(),
				Username: "testuser",
				Password: string(hashed),
				Role:     schema.RoleCustomer,
				Status:   1,
			}, nil
		},
	}
	app := setupAuthApp(t, mock)

	resp := postJSON(t, app, "/api/v1/auth/login", `{"username":"testuser","password":"wrongpassword"}`)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-login-notfound")
	mock := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*schema.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	app := setupAuthApp(t, mock)

	resp := postJSON(t, app, "/api/v1/auth/login", `{"username":"nonexistent","password":"secret123"}`)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLogin_DatabaseError(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-login-dberr")
	mock := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*schema.User, error) {
			return nil, errDB
		},
	}
	app := setupAuthApp(t, mock)

	resp := postJSON(t, app, "/api/v1/auth/login", `{"username":"testuser","password":"secret123"}`)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

// ─── Refresh Test ──────────────────────────────────────────────

func TestRefresh_NotImplemented(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-refresh")
	app := setupAuthApp(t, &mockUserRepo{})

	resp := postJSON(t, app, "/api/v1/auth/refresh", `{"refreshToken":"any"}`)

	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", resp.StatusCode)
	}
}

// ─── Logout Test ───────────────────────────────────────────────

func TestLogout_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-logout")
	app := setupAuthApp(t, &mockUserRepo{})

	resp := postJSON(t, app, "/api/v1/auth/logout", "")

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
