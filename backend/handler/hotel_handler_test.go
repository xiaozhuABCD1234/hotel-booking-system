package handler_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"backend/handler"
	schema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
)

type mockHotelRepo struct {
	repo.HotelRepository

	calledMinPrice     *float64
	calledMaxPrice     *float64
	calledCheckInDate  *time.Time
	calledCheckOutDate *time.Time
}

func (m *mockHotelRepo) FindAll(_ context.Context, _, _ int, _ *int, _ *int16, _ string, minPrice, maxPrice *float64, checkInDate, checkOutDate *time.Time) ([]schema.Hotel, int64, error) {
	m.calledMinPrice = minPrice
	m.calledMaxPrice = maxPrice
	m.calledCheckInDate = checkInDate
	m.calledCheckOutDate = checkOutDate
	return nil, 0, nil
}

func setupHotelListApp(t *testing.T, hotelRepo *mockHotelRepo) *fiber.App {
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
	h := handler.NewHotelHandler(hotelRepo, nil, nil, nil)
	app.Get("/hotels", h.List)
	return app
}

func getHotels(t *testing.T, app *fiber.App, query string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest("GET", "/hotels"+query, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func TestHotelHandler_List_ValidParams(t *testing.T) {
	mock := &mockHotelRepo{}
	app := setupHotelListApp(t, mock)

	resp := getHotels(t, app, "?minPrice=200&maxPrice=500&checkInDate=2026-07-15&checkOutDate=2026-07-17")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	if mock.calledMinPrice == nil || *mock.calledMinPrice != 200.0 {
		t.Errorf("expected minPrice=200, got %v", mock.calledMinPrice)
	}
	if mock.calledMaxPrice == nil || *mock.calledMaxPrice != 500.0 {
		t.Errorf("expected maxPrice=500, got %v", mock.calledMaxPrice)
	}
	if mock.calledCheckInDate == nil {
		t.Fatal("expected checkInDate to be passed to repo, got nil")
	}
	expected := mock.calledCheckInDate.Format("2006-01-02")
	if expected != "2026-07-15" {
		t.Errorf("expected checkInDate=2026-07-15, got %s", expected)
	}
	if mock.calledCheckOutDate == nil {
		t.Fatal("expected checkOutDate to be passed to repo, got nil")
	}
	expected = mock.calledCheckOutDate.Format("2006-01-02")
	if expected != "2026-07-17" {
		t.Errorf("expected checkOutDate=2026-07-17, got %s", expected)
	}
}

func TestHotelHandler_List_NegativePrice(t *testing.T) {
	mock := &mockHotelRepo{}
	app := setupHotelListApp(t, mock)

	resp := getHotels(t, app, "?minPrice=-100")

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHotelHandler_List_InvalidDateFormat(t *testing.T) {
	mock := &mockHotelRepo{}
	app := setupHotelListApp(t, mock)

	resp := getHotels(t, app, "?checkInDate=2026/07/15")

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHotelHandler_List_InvertedDates(t *testing.T) {
	mock := &mockHotelRepo{}
	app := setupHotelListApp(t, mock)

	resp := getHotels(t, app, "?checkInDate=2026-07-17&checkOutDate=2026-07-15")

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHotelHandler_List_PastDate(t *testing.T) {
	mock := &mockHotelRepo{}
	app := setupHotelListApp(t, mock)

	resp := getHotels(t, app, "?checkInDate=2020-01-01")

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
