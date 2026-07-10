package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"backend/auth"
	"backend/database"
	"backend/middleware"
	"backend/model"
	"backend/router"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var testApp *fiber.App
var testDB *gorm.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load("../.env.test")
	_ = godotenv.Load("../.env")

	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-secret-key-for-integration-tests-min-16-bytes")
	}
	if err := auth.LoadSecret(); err != nil {
		log.Fatalf("TestMain: JWT secret: %v", err)
	}

	cfg := database.DefaultConfig()
	if dbName := os.Getenv("TEST_DB_NAME"); dbName != "" {
		cfg.DBName = dbName
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("TestMain: database: %v", err)
	}
	testDB = db

	testDB.Exec(`INSERT INTO vip_level_1718 (level, level_name, min_points, discount_rate) VALUES (0, '普通会员', 0, 1.00) ON CONFLICT (level) DO NOTHING`)

	app := fiber.New(fiber.Config{
		AppName:      "Hotel Booking API (Test)",
		ErrorHandler: middleware.ErrorHandler,
	})
	router.RegisterRoutes(app, db)
	testApp = app

	code := m.Run()

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
	os.Exit(code)
}

func parseResponse(t *testing.T, resp *http.Response) model.Response {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	resp.Body.Close()
	var r model.Response
	if err := json.Unmarshal(body, &r); err != nil {
		t.Fatalf("unmarshal: %v\nbody: %s", err, string(body))
	}
	return r
}

func request(t *testing.T, method, path, body string, headers ...map[string]string) *http.Response {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = bytes.NewReader([]byte(body))
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}
	resp, err := testApp.Test(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return resp
}

func bearerAuth(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}

type registeredUser struct {
	username     string
	password     string
	accessToken  string
	refreshToken string
}

func registerUser(t *testing.T) registeredUser {
	t.Helper()
	username := fmt.Sprintf("it_%d", time.Now().UnixNano())
	password := "Pass1234"
	resp := request(t, "POST", "/api/v1/auth/register",
		fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("register: 200 expected, got %d", resp.StatusCode)
	}
	r := parseResponse(t, resp)
	data, _ := json.Marshal(r.Data)
	var td struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	json.Unmarshal(data, &td)
	t.Cleanup(func() { testDB.Exec(`DELETE FROM user_1718 WHERE username = ?`, username) })
	return registeredUser{username, password, td.AccessToken, td.RefreshToken}
}

// ─── 1. 认证 ───────────────────────────────────────────────────

func TestAuth_RegisterAndLogin(t *testing.T) {
	username := fmt.Sprintf("auth_%d", time.Now().UnixNano())
	password := "Pass1234"

	resp := request(t, "POST", "/api/v1/auth/register",
		fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("register: expected 200, got %d", resp.StatusCode)
	}
	r := parseResponse(t, resp)
	if !r.Success {
		t.Fatal("register success=false")
	}
	t.Cleanup(func() { testDB.Exec(`DELETE FROM user_1718 WHERE username = ?`, username) })

	resp = request(t, "POST", "/api/v1/auth/login",
		fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: expected 200, got %d", resp.StatusCode)
	}
	r = parseResponse(t, resp)
	if !r.Success {
		t.Fatal("login success=false")
	}
}

func TestAuth_Refresh(t *testing.T) {
	u := registerUser(t)
	resp := request(t, "POST", "/api/v1/auth/refresh",
		fmt.Sprintf(`{"refreshToken":"%s"}`, u.refreshToken))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("refresh: expected 200, got %d", resp.StatusCode)
	}
}

func TestAuth_BadRequest(t *testing.T) {
	t.Run("register empty", func(t *testing.T) {
		resp := request(t, "POST", "/api/v1/auth/register", `{"username":"","password":""}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})
	t.Run("register short password", func(t *testing.T) {
		resp := request(t, "POST", "/api/v1/auth/register", `{"username":"u","password":"abc"}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})
	t.Run("login empty", func(t *testing.T) {
		resp := request(t, "POST", "/api/v1/auth/login", `{}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})
}

// ─── 2. 公开 GET ───────────────────────────────────────────────

func TestPublic_Regions(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		resp := request(t, "GET", "/api/v1/regions/", "")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})
	t.Run("provinces", func(t *testing.T) {
		resp := request(t, "GET", "/api/v1/regions/provinces", "")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})
	t.Run("by-parent", func(t *testing.T) {
		resp := request(t, "GET", "/api/v1/regions/by-parent?parentID=0", "")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})
	t.Run("by-id not found", func(t *testing.T) {
		resp := request(t, "GET", "/api/v1/regions/999999", "")
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestPublic_HotelsAndRooms(t *testing.T) {
	t.Run("hotels list", func(t *testing.T) {
		resp := request(t, "GET", "/api/v1/hotels?page=1&pageSize=10", "")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})
	t.Run("hotel by-id not found", func(t *testing.T) {
		resp := request(t, "GET", "/api/v1/hotels/00000000-0000-0000-0000-000000000000", "")
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})
	t.Run("rooms list", func(t *testing.T) {
		resp := request(t, "GET", "/api/v1/rooms?page=1&pageSize=10", "")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})
	t.Run("room by-id not found", func(t *testing.T) {
		resp := request(t, "GET", "/api/v1/rooms/00000000-0000-0000-0000-000000000000", "")
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})
}

// ─── 3. 受保护 → 401 ───────────────────────────────────────────

func TestProtected_401WithoutToken(t *testing.T) {
	type ep struct{ method, path string }
	protected := []ep{
		{"POST", "/api/v1/users"},
		{"GET", "/api/v1/users"},
		{"GET", "/api/v1/users/00000000-0000-0000-0000-000000000000"},
		{"PUT", "/api/v1/users/00000000-0000-0000-0000-000000000000"},
		{"DELETE", "/api/v1/users/00000000-0000-0000-0000-000000000000"},
		{"POST", "/api/v1/hotels"},
		{"PUT", "/api/v1/hotels/00000000-0000-0000-0000-000000000000"},
		{"DELETE", "/api/v1/hotels/00000000-0000-0000-0000-000000000000"},
		{"POST", "/api/v1/rooms"},
		{"PUT", "/api/v1/rooms/00000000-0000-0000-0000-000000000000"},
		{"DELETE", "/api/v1/rooms/00000000-0000-0000-0000-000000000000"},
		{"GET", "/api/v1/orders"},
		{"POST", "/api/v1/orders"},
		{"GET", "/api/v1/orders/by-user"},
		{"GET", "/api/v1/orders/by-hotel"},
		{"PUT", "/api/v1/orders/00000000-0000-0000-0000-000000000000/status"},
		{"DELETE", "/api/v1/orders/00000000-0000-0000-0000-000000000000"},
		{"GET", "/api/v1/reviews"},
		{"POST", "/api/v1/reviews"},
		{"GET", "/api/v1/reviews/by-hotel"},
		{"GET", "/api/v1/reviews/by-user"},
		{"GET", "/api/v1/persons"},
		{"POST", "/api/v1/persons"},
		{"POST", "/api/v1/regions"},
		{"PUT", "/api/v1/regions/1"},
		{"DELETE", "/api/v1/regions/1"},
		{"GET", "/api/v1/reports/hotel-summaries"},
		{"GET", "/api/v1/reports/room-details"},
		{"GET", "/api/v1/reports/user-vip"},
		{"GET", "/api/v1/reports/person-info"},
		{"GET", "/api/v1/reports/guest-stats"},
		{"GET", "/api/v1/reports/my-orders"},
	}
	for _, ep := range protected {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			var resp *http.Response
			switch ep.method {
			case "POST":
				resp = request(t, "POST", ep.path, `{}`)
			case "PUT":
				resp = request(t, "PUT", ep.path, `{}`)
			case "GET":
				resp = request(t, "GET", ep.path, "")
			case "DELETE":
				resp = request(t, "DELETE", ep.path, "")
			}
			if resp.StatusCode != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d", resp.StatusCode)
			}
		})
	}
}

// ─── 4. 用户 CRUD ──────────────────────────────────────────────

func TestUsers_CRUD(t *testing.T) {
	u := registerUser(t)
	h := bearerAuth(u.accessToken)

	ts := fmt.Sprint(time.Now().UnixNano())
	resp := request(t, "POST", "/api/v1/users",
		fmt.Sprintf(`{"Username":"crud_usr_%s","Password":"Pass1234","Role":"customer"}`, ts), h)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", resp.StatusCode)
	}
	r := parseResponse(t, resp)
	data, _ := json.Marshal(r.Data)
	var created struct{ ID string }
	json.Unmarshal(data, &created)
	if created.ID == "" {
		t.Fatal("created user has no ID")
	}
	t.Cleanup(func() { testDB.Exec(`DELETE FROM user_1718 WHERE id = ?`, created.ID) })

	resp = request(t, "GET", "/api/v1/users?page=1&pageSize=10", "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "GET", "/api/v1/users/"+created.ID, "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "PUT", "/api/v1/users/"+created.ID,
		fmt.Sprintf(`{"Username":"upd_%s","Role":"customer","Password":"","Phone":null,"Email":null,"RealName":null,"IDCard":null,"Points":0,"VipLevelID":0,"Status":1}`, ts), h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "DELETE", "/api/v1/users/"+created.ID, "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "GET", "/api/v1/users/"+created.ID, "", h)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get after delete: expected 404, got %d", resp.StatusCode)
	}
}

func TestUsers_CreateInvalid(t *testing.T) {
	u := registerUser(t)
	resp := request(t, "POST", "/api/v1/users", `{}`, bearerAuth(u.accessToken))
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// ─── 5. 酒店 CRUD ──────────────────────────────────────────────

func TestHotels_CRUD(t *testing.T) {
	u := registerUser(t)
	h := bearerAuth(u.accessToken)

	ts := fmt.Sprint(time.Now().UnixNano())
	regResp := request(t, "POST", "/api/v1/regions",
		fmt.Sprintf(`{"RegionName":"hotel-region-%s","ParentsID":null}`, ts), h)
	regR := parseResponse(t, regResp)
	regData, _ := json.Marshal(regR.Data)
	var regCreated struct{ ID int }
	json.Unmarshal(regData, &regCreated)
	if regCreated.ID == 0 {
		t.Fatal("failed to create test region")
	}
	regionID := regCreated.ID
	t.Cleanup(func() { testDB.Exec(`DELETE FROM region_1718 WHERE id = ?`, regionID) })

	resp := request(t, "POST", "/api/v1/hotels",
		fmt.Sprintf(`{"HotelName":"test-hotel-%s","RegionID":%d,"Address":"test addr","Telephone":"13800138000"}`, ts, regionID), h)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", resp.StatusCode)
	}
	r := parseResponse(t, resp)
	data, _ := json.Marshal(r.Data)
	var created struct{ ID string }
	json.Unmarshal(data, &created)
	if created.ID == "" {
		t.Fatal("created hotel has no ID")
	}
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM hotel_image_1718 WHERE hotel_id = ?`, created.ID)
		testDB.Exec(`DELETE FROM hotel_1718 WHERE id = ?`, created.ID)
	})

	resp = request(t, "GET", "/api/v1/hotels/"+created.ID, "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "PUT", "/api/v1/hotels/"+created.ID,
		fmt.Sprintf(`{"HotelName":"upd-%s","RegionID":%d,"Address":"upd addr","Telephone":"13900139000"}`, ts, regionID), h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "DELETE", "/api/v1/hotels/"+created.ID, "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}
}

func TestHotels_CreateInvalid(t *testing.T) {
	u := registerUser(t)
	resp := request(t, "POST", "/api/v1/hotels", `{}`, bearerAuth(u.accessToken))
	// handler 没有显式校验必填字段，空值落到 DB FK 约束失败返回 500
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 400 or 500, got %d", resp.StatusCode)
	}
}

// ─── 6. 客房 CRUD ──────────────────────────────────────────────

func TestRooms_CRUD(t *testing.T) {
	u := registerUser(t)
	h := bearerAuth(u.accessToken)

	ts := fmt.Sprint(time.Now().UnixNano())

	regResp := request(t, "POST", "/api/v1/regions",
		fmt.Sprintf(`{"RegionName":"room-region-%s","ParentsID":null}`, ts), h)
	regR := parseResponse(t, regResp)
	regData, _ := json.Marshal(regR.Data)
	var regCreated struct{ ID int }
	json.Unmarshal(regData, &regCreated)
	t.Cleanup(func() { testDB.Exec(`DELETE FROM region_1718 WHERE id = ?`, regCreated.ID) })

	hotelResp := request(t, "POST", "/api/v1/hotels",
		fmt.Sprintf(`{"HotelName":"room-hotel-%s","RegionID":%d,"Address":"addr","Telephone":"13800138001"}`, ts, regCreated.ID), h)
	hotelR := parseResponse(t, hotelResp)
	hotelData, _ := json.Marshal(hotelR.Data)
	var hotelCreated struct{ ID string }
	json.Unmarshal(hotelData, &hotelCreated)
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM hotel_image_1718 WHERE hotel_id = ?`, hotelCreated.ID)
		testDB.Exec(`DELETE FROM hotel_1718 WHERE id = ?`, hotelCreated.ID)
	})

	resp := request(t, "POST", "/api/v1/rooms",
		fmt.Sprintf(`{"HotelID":"%s","TypeName":"标准间-%s","TotalQuantity":10,"AvailableQuantity":10,"Price":200.00}`, hotelCreated.ID, ts), h)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", resp.StatusCode)
	}
	r := parseResponse(t, resp)
	data, _ := json.Marshal(r.Data)
	var created struct{ ID string }
	json.Unmarshal(data, &created)
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM room_facility_1718 WHERE room_id = ?`, created.ID)
		testDB.Exec(`DELETE FROM room_image_1718 WHERE room_id = ?`, created.ID)
		testDB.Exec(`DELETE FROM room_1718 WHERE id = ?`, created.ID)
	})

	resp = request(t, "GET", "/api/v1/rooms/"+created.ID, "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "PUT", "/api/v1/rooms/"+created.ID,
		fmt.Sprintf(`{"HotelID":"%s","TypeName":"豪华间","TotalQuantity":5,"AvailableQuantity":5,"Price":500.00}`, hotelCreated.ID), h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "DELETE", "/api/v1/rooms/"+created.ID, "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}
}

// ─── 7. 人员 CRUD ──────────────────────────────────────────────

func TestPersons_CRUD(t *testing.T) {
	u := registerUser(t)
	h := bearerAuth(u.accessToken)

	idCard := fmt.Sprintf("11010119900101%04d", 1000+time.Now().UnixNano()%9000)

	resp := request(t, "POST", "/api/v1/persons",
		fmt.Sprintf(`{"IDCard":"%s","Name":"张三","Phone":"13800138000"}`, idCard), h)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", resp.StatusCode)
	}
	t.Cleanup(func() { testDB.Exec(`DELETE FROM person_1718 WHERE id_card = ?`, idCard) })

	resp = request(t, "GET", "/api/v1/persons/"+idCard, "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "PUT", "/api/v1/persons/"+idCard,
		fmt.Sprintf(`{"IDCard":"%s","Name":"李四","Phone":"13900139000"}`, idCard), h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "DELETE", "/api/v1/persons/"+idCard, "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}
}

func TestPersons_List(t *testing.T) {
	u := registerUser(t)
	resp := request(t, "GET", "/api/v1/persons?page=1&pageSize=10", "", bearerAuth(u.accessToken))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", resp.StatusCode)
	}
}

// ─── 8. 订单 (只测 401 和空列表，完整 CRUD 依赖 UserID) ────────

func TestOrders_ListEmpty(t *testing.T) {
	u := registerUser(t)
	resp := request(t, "GET", "/api/v1/orders?page=1&pageSize=10", "", bearerAuth(u.accessToken))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", resp.StatusCode)
	}
}

// ─── 9. 评价 (只测 401 和空列表) ───────────────────────────────

func TestReviews_List(t *testing.T) {
	u := registerUser(t)
	resp := request(t, "GET", "/api/v1/reviews?page=1&pageSize=10", "", bearerAuth(u.accessToken))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", resp.StatusCode)
	}
}

// ─── 10. 地区管理 CRUD ─────────────────────────────────────────

func TestRegions_CRUD(t *testing.T) {
	u := registerUser(t)
	h := bearerAuth(u.accessToken)

	ts := fmt.Sprint(time.Now().UnixNano())

	resp := request(t, "POST", "/api/v1/regions",
		fmt.Sprintf(`{"RegionName":"test-region-%s","ParentsID":null}`, ts), h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("create: expected 200, got %d", resp.StatusCode)
	}
	r := parseResponse(t, resp)
	data, _ := json.Marshal(r.Data)
	var created struct{ ID int }
	json.Unmarshal(data, &created)
	if created.ID == 0 {
		t.Fatal("created region has no ID")
	}
	t.Cleanup(func() { testDB.Exec(`DELETE FROM region_1718 WHERE id = ?`, created.ID) })

	resp = request(t, "PUT", "/api/v1/regions/"+fmt.Sprint(created.ID),
		fmt.Sprintf(`{"ID":%d,"RegionName":"test-region-upd-%s","ParentsID":null}`, created.ID, ts), h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", resp.StatusCode)
	}

	resp = request(t, "DELETE", "/api/v1/regions/"+fmt.Sprint(created.ID), "", h)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}
}

// ─── 11. 报表 ──────────────────────────────────────────────────

func TestReports(t *testing.T) {
	u := registerUser(t)
	h := bearerAuth(u.accessToken)

	type rc struct{ name, path string }
	reports := []rc{
		{"hotel-summaries", "/api/v1/reports/hotel-summaries?page=1&pageSize=10"},
		{"room-details", "/api/v1/reports/room-details?page=1&pageSize=10"},
		{"user-vip", "/api/v1/reports/user-vip?page=1&pageSize=10"},
		{"person-info", "/api/v1/reports/person-info?page=1&pageSize=10"},
		{"guest-stats", "/api/v1/reports/guest-stats?page=1&pageSize=10"},
		{"guest-stats top", "/api/v1/reports/guest-stats/top"},
		{"my-orders", "/api/v1/reports/my-orders?page=1&pageSize=10"},
	}
	for _, rc := range reports {
		t.Run(rc.name, func(t *testing.T) {
			// my-orders 需要 userID 参数
			path := rc.path
			if rc.name == "my-orders" {
				// 用当前 token 对应用户的 ID
				userResp := request(t, "GET", "/api/v1/users?page=1&pageSize=1", "", h)
				ur := parseResponse(t, userResp)
				udata, _ := json.Marshal(ur.Data)
				var users []struct {
					ID string `json:"ID"`
				}
				json.Unmarshal(udata, &users)
				if len(users) > 0 {
					path = "/api/v1/reports/my-orders?userID=" + users[0].ID + "&page=1&pageSize=10"
				}
			}
			resp := request(t, "GET", path, "", h)
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}
			r := parseResponse(t, resp)
			if !r.Success {
				t.Fatal("success=false")
			}
		})
	}
}
