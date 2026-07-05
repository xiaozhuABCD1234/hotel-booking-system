// Package router 提供 RESTful API 路由注册，将所有 handler 挂载到 Fiber App。
package router

import (
	"backend/handler"
	"backend/repo"
	"backend/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// RegisterRoutes 创建所有 repo 和 handler 实例，注册路由到 /api/v1。
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

	// ─── 创建 handlers ─────────────────────────────────────────
	userH := handler.NewUserHandler(userRepo, vipLevelRepo)
	hotelH := handler.NewHotelHandler(hotelRepo, roomRepo)
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

	// ─── /api/v1 路由组 ────────────────────────────────────────
	v1 := app.Group("/api/v1")

	// 用户管理
	users := v1.Group("/users")
	users.Get("/", userH.List)
	users.Get("/:id", userH.GetByID)
	users.Post("/", userH.Create)
	users.Put("/:id", userH.Update)
	users.Delete("/:id", userH.Delete)

	// 酒店管理
	hotels := v1.Group("/hotels")
	hotels.Get("/", hotelH.List)
	hotels.Get("/:id", hotelH.GetByID)
	hotels.Post("/", hotelH.Create)
	hotels.Put("/:id", hotelH.Update)
	hotels.Delete("/:id", hotelH.Delete)

	// 客房管理
	rooms := v1.Group("/rooms")
	rooms.Get("/", hotelH.ListRooms)
	rooms.Get("/:id", hotelH.GetRoomByID)
	rooms.Post("/", hotelH.CreateRoom)
	rooms.Put("/:id", hotelH.UpdateRoom)
	rooms.Delete("/:id", hotelH.DeleteRoom)

	// 订单管理
	orders := v1.Group("/orders")
	orders.Get("/", orderH.List)
	orders.Get("/:id", orderH.GetByID)
	orders.Post("/", orderH.Create)
	orders.Put("/:id/status", orderH.UpdateStatus)
	orders.Delete("/:id", orderH.Delete)
	orders.Get("/by-user", orderH.ListByUserID)
	orders.Get("/by-hotel", orderH.ListByHotelID)

	// 评价管理
	reviews := v1.Group("/reviews")
	reviews.Get("/", reviewH.List)
	reviews.Get("/:id", reviewH.GetByID)
	reviews.Post("/", reviewH.Create)
	reviews.Put("/:id", reviewH.Update)
	reviews.Delete("/:id", reviewH.Delete)
	reviews.Get("/by-hotel", reviewH.ListByHotelID)
	reviews.Get("/by-user", reviewH.ListByUserID)

	// 人员管理
	persons := v1.Group("/persons")
	persons.Get("/", personH.List)
	persons.Get("/:idCard", personH.GetByIDCard)
	persons.Post("/", personH.Create)
	persons.Put("/:idCard", personH.Update)
	persons.Delete("/:idCard", personH.Delete)

	// 地区管理
	regions := v1.Group("/regions")
	regions.Get("/", regionH.List)
	// 固定路径必须在 /:id 之前注册，否则会被 :id 参数捕获
	regions.Get("/provinces", regionH.ListProvinces)
	regions.Get("/by-parent", regionH.ListByParent)
	regions.Get("/:id", regionH.GetByID)
	regions.Post("/", regionH.Create)
	regions.Put("/:id", regionH.Update)
	regions.Delete("/:id", regionH.Delete)

	// 汇总报表
	reports := v1.Group("/reports")
	reports.Get("/hotel-summaries", reportH.HotelSummaries)
	reports.Get("/room-details", reportH.RoomDetails)
	reports.Get("/room-details/by-hotel", reportH.RoomDetailsByHotelID)
	reports.Get("/order-full/by-user", reportH.OrderFullByUserID)
	reports.Get("/order-full/by-hotel", reportH.OrderFullByHotelID)
	reports.Get("/review-full/by-hotel", reportH.ReviewFullByHotelID)
	reports.Get("/review-full/by-user", reportH.ReviewFullByUserID)
	reports.Get("/user-vip", reportH.UserVipList)
	reports.Get("/person-info", reportH.PersonInfoList)
	reports.Get("/guest-stats", reportH.GuestStats)
	reports.Get("/guest-stats/top", reportH.TopGuests)
	reports.Get("/my-orders", reportH.MyOrders)
}
