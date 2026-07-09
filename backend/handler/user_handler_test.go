package handler_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"backend/handler"
	schema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// ─── Mock Repository ───────────────────────────────────────────

type mockUserUpdateRepo struct {
	user        *schema.User
	updatedUser *schema.User
	repo.UserRepository
}

func (m *mockUserUpdateRepo) FindByID(_ context.Context, _ uuid.UUID) (*schema.User, error) {
	return m.user, nil
}

func (m *mockUserUpdateRepo) Update(_ context.Context, user *schema.User) error {
	m.updatedUser = user
	return nil
}

type mockVipLevelRepo struct {
	repo.VipLevelRepository
}

// ─── Test Helpers ──────────────────────────────────────────────

func setupUserUpdateApp(t *testing.T, userRepo *mockUserUpdateRepo) *fiber.App {
	t.Helper()
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
	h := handler.NewUserHandler(userRepo, &mockVipLevelRepo{})
	app.Put("/users/:id", h.Update)
	return app
}

func putJSON(t *testing.T, app *fiber.App, path, body string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest("PUT", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

// ─── Helper Functions ──────────────────────────────────────────

func strPtr(s string) *string { return &s }

func float64Ptr(f float64) *float64 { return &f }

// ─── Tests ─────────────────────────────────────────────────────

func TestUpdate_NewFields(t *testing.T) {
	userID := uuid.New()
	existingUser := &schema.User{
		ID:       userID,
		Username: "testuser",
		Phone:    strPtr("13800138000"),
	}
	mock := &mockUserUpdateRepo{user: existingUser}
	app := setupUserUpdateApp(t, mock)

	body := `{"occupation":"工程师","education":"本科","income":15000}`
	resp := putJSON(t, app, "/users/"+userID.String(), body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if mock.updatedUser == nil {
		t.Fatal("expected Update to be called, but updatedUser is nil")
	}
	if mock.updatedUser.Occupation == nil || *mock.updatedUser.Occupation != "工程师" {
		t.Errorf("expected occupation=工程师, got %v", mock.updatedUser.Occupation)
	}
	if mock.updatedUser.Education == nil || *mock.updatedUser.Education != "本科" {
		t.Errorf("expected education=本科, got %v", mock.updatedUser.Education)
	}
	if mock.updatedUser.Income == nil || *mock.updatedUser.Income != 15000 {
		t.Errorf("expected income=15000, got %v", mock.updatedUser.Income)
	}
}

func TestUpdate_PartialFields(t *testing.T) {
	userID := uuid.New()
	existingUser := &schema.User{
		ID:       userID,
		Username: "testuser",
	}
	mock := &mockUserUpdateRepo{user: existingUser}
	app := setupUserUpdateApp(t, mock)

	body := `{"occupation":"测试员"}`
	resp := putJSON(t, app, "/users/"+userID.String(), body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if mock.updatedUser.Occupation == nil || *mock.updatedUser.Occupation != "测试员" {
		t.Errorf("expected occupation=测试员, got %v", mock.updatedUser.Occupation)
	}
	if mock.updatedUser.Education != nil {
		t.Errorf("expected education=nil (not provided), got %v", mock.updatedUser.Education)
	}
	if mock.updatedUser.Income != nil {
		t.Errorf("expected income=nil (not provided), got %v", mock.updatedUser.Income)
	}
}

func TestUpdate_NoNewFields(t *testing.T) {
	userID := uuid.New()
	origOccupation := "医生"
	existingUser := &schema.User{
		ID:         userID,
		Username:   "testuser",
		Occupation: &origOccupation,
	}
	mock := &mockUserUpdateRepo{user: existingUser}
	app := setupUserUpdateApp(t, mock)

	body := `{"phone":"13900139000"}`
	resp := putJSON(t, app, "/users/"+userID.String(), body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if mock.updatedUser.Occupation == nil || *mock.updatedUser.Occupation != "医生" {
		t.Errorf("expected occupation=医生 (unchanged), got %v", mock.updatedUser.Occupation)
	}
}
