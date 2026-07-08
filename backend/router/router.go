// Package router 提供 RESTful API 路由注册，将所有 handler 挂载到 Fiber App。
package router

import (
	"backend/handler"
	"backend/middleware"
	"backend/repo"
	"backend/service"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// RegisterRoutes 创建所有 repo、service 和 handler 实例，注册路由到 /api/v1。
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// ─── 创建 repos ────────────────────────────────────────────
	userRepo := repo.NewUserRepo(db)
	vipLevelRepo := repo.NewVipLevelRepo(db)
	hotelRepo := repo.NewHotelRepo(db)
	roomRepo := repo.NewRoomRepo(db)
	orderRepo := repo.NewOrderRepo(db)
	reviewRepo := repo.NewReviewRepo(db)
	personRepo := repo.NewPersonRepo(db)
	regionRepo := repo.NewRegionRepo(db)
	blacklistRepo := repo.NewBlacklistRepo(db)
	hotelImageRepo := repo.NewHotelImageRepo(db)

	// 视图 repos
	hotelSummaryRepo := repo.NewHotelSummaryRepo(db)
	roomDetailsRepo := repo.NewRoomDetailsRepo(db)
	orderFullRepo := repo.NewOrderFullRepo(db)
	reviewFullRepo := repo.NewReviewFullRepo(db)
	userVipRepo := repo.NewUserVipRepo(db)
	personInfoRepo := repo.NewPersonInfoRepo(db)
	guestBookingStatsRepo := repo.NewGuestBookingStatsRepo(db)
	myOrdersRepo := repo.NewMyOrdersRepo(db)

	// ─── 创建 services ───────────────────────────────────────
	orderSvc := service.NewOrderService(orderRepo)
	cosSvc := service.NewCOSService()

	// ─── 创建 handlers ─────────────────────────────────────────
	authH := handler.NewAuthHandler(userRepo, blacklistRepo)
	userH := handler.NewUserHandler(userRepo, vipLevelRepo)
	hotelH := handler.NewHotelHandler(hotelRepo, roomRepo, hotelImageRepo, cosSvc)
	orderH := handler.NewOrderHandler(orderSvc)
	reviewH := handler.NewReviewHandler(reviewRepo)
	personH := handler.NewPersonHandler(personRepo)
	regionH := handler.NewRegionHandler(regionRepo)
	reportH := handler.NewReportHandler(
		hotelSummaryRepo, roomDetailsRepo,
		orderFullRepo, reviewFullRepo,
		userVipRepo, personInfoRepo,
		guestBookingStatsRepo, myOrdersRepo,
	)

	// ─── Swagger 文档（不依赖 /api/v1 前缀）─────────────────────
	app.Get("/swagger/*", swaggo.HandlerDefault)

	// ─── /api/v1 路由组 ────────────────────────────────────────
	v1 := app.Group("/api/v1")

	// ═══════════════════════════════════════════════════════════
	// 认证路由（公开，无需 JWT）
	// ═══════════════════════════════════════════════════════════
	auth := v1.Group("/auth")
	auth.Post("/register", authH.Register)
	auth.Post("/login", authH.Login)
	auth.Post("/refresh", authH.Refresh)
	auth.Post("/logout", authH.Logout)

	// ═══════════════════════════════════════════════════════════
	// 公开路由（无需认证）
	// ═══════════════════════════════════════════════════════════
	// 地区查询（下拉选择等场景无需登录）
	regions := v1.Group("/regions")
	regions.Get("/", regionH.List)
	regions.Get("/provinces", regionH.ListProvinces)
	regions.Get("/by-parent", regionH.ListByParent)
	regions.Get("/:id", regionH.GetByID)

	// 酒店/客房浏览（未登录用户也可搜索）
	hotels := v1.Group("/hotels")
	hotels.Get("/", hotelH.List)
	hotels.Get("/:id", hotelH.GetByID)

	rooms := v1.Group("/rooms")
	rooms.Get("/", hotelH.ListRooms)
	rooms.Get("/:id", hotelH.GetRoomByID)

	// ═══════════════════════════════════════════════════════════
	// 需要认证的路由（JWT 中间件保护）
	// ═══════════════════════════════════════════════════════════
	protected := v1.Group("", middleware.JWTAuth(blacklistRepo))

	// 用户管理（需认证）
	protectedUsers := protected.Group("/users")
	protectedUsers.Post("/", userH.Create)
	protectedUsers.Get("/", userH.List)
	protectedUsers.Get("/:id", userH.GetByID)
	protectedUsers.Put("/:id", userH.Update)
	protectedUsers.Delete("/:id", userH.Delete)

	// 酒店/客房管理（需认证）
	protectedHotels := protected.Group("/hotels")
	protectedHotels.Post("/", hotelH.Create)
	protectedHotels.Put("/:id", hotelH.Update)
	protectedHotels.Delete("/:id", hotelH.Delete)
	if cosSvc != nil {
		protectedHotels.Post("/:id/images", hotelH.UploadImage)
		protectedHotels.Delete("/:id/images", hotelH.DeleteImage)
	}

	protectedRooms := protected.Group("/rooms")
	protectedRooms.Post("/", hotelH.CreateRoom)
	protectedRooms.Put("/:id", hotelH.UpdateRoom)
	protectedRooms.Delete("/:id", hotelH.DeleteRoom)

	// 订单管理（需认证）
	protectedOrders := protected.Group("/orders")
	protectedOrders.Get("/", orderH.List)
	// 静态路由必须在 :id 参数路由之前注册
	protectedOrders.Get("/by-user", orderH.ListByUserID)
	protectedOrders.Get("/by-hotel", orderH.ListByHotelID)
	protectedOrders.Get("/:id", orderH.GetByID)
	protectedOrders.Post("/", orderH.Create)
	protectedOrders.Put("/:id/status", orderH.UpdateStatus)
	protectedOrders.Delete("/:id", orderH.Delete)

	// 评价管理（需认证）
	protectedReviews := protected.Group("/reviews")
	protectedReviews.Get("/", reviewH.List)
	// 静态路由必须在 :id 参数路由之前注册
	protectedReviews.Get("/by-hotel", reviewH.ListByHotelID)
	protectedReviews.Get("/by-user", reviewH.ListByUserID)
	protectedReviews.Get("/:id", reviewH.GetByID)
	protectedReviews.Post("/", reviewH.Create)
	protectedReviews.Put("/:id", reviewH.Update)
	protectedReviews.Delete("/:id", reviewH.Delete)

	// 人员管理（需认证）
	protectedPersons := protected.Group("/persons")
	protectedPersons.Get("/", personH.List)
	protectedPersons.Get("/:idCard", personH.GetByIDCard)
	protectedPersons.Post("/", personH.Create)
	protectedPersons.Put("/:idCard", personH.Update)
	protectedPersons.Delete("/:idCard", personH.Delete)

	// 地区管理（需认证：增删改）
	protectedRegions := protected.Group("/regions")
	protectedRegions.Post("/", regionH.Create)
	protectedRegions.Put("/:id", regionH.Update)
	protectedRegions.Delete("/:id", regionH.Delete)

	// 汇总报表（需认证）
	protectedReports := protected.Group("/reports")
	protectedReports.Get("/hotel-summaries", reportH.HotelSummaries)
	protectedReports.Get("/room-details", reportH.RoomDetails)
	protectedReports.Get("/room-details/by-hotel", reportH.RoomDetailsByHotelID)
	protectedReports.Get("/order-full/by-user", reportH.OrderFullByUserID)
	protectedReports.Get("/order-full/by-hotel", reportH.OrderFullByHotelID)
	protectedReports.Get("/review-full/by-hotel", reportH.ReviewFullByHotelID)
	protectedReports.Get("/review-full/by-user", reportH.ReviewFullByUserID)
	protectedReports.Get("/user-vip", reportH.UserVipList)
	protectedReports.Get("/person-info", reportH.PersonInfoList)
	protectedReports.Get("/guest-stats", reportH.GuestStats)
	protectedReports.Get("/guest-stats/top", reportH.TopGuests)
	protectedReports.Get("/my-orders", reportH.MyOrders)
}
