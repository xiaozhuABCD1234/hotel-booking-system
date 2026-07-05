package handler

import (
	"net/http"
	"strconv"

	"backend/model"
	"backend/model/view"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// ReportHandler 视图报表 HTTP 处理器，聚合所有只读视图的查询接口。
type ReportHandler struct {
	hotelSummaries    repo.HotelSummaryRepository
	roomDetails       repo.RoomDetailsRepository
	orderFull         repo.OrderFullRepository
	reviewFull        repo.ReviewFullRepository
	userVip           repo.UserVipRepository
	personInfo        repo.PersonInfoRepository
	guestBookingStats repo.GuestBookingStatsRepository
	myOrders          repo.MyOrdersRepository
}

// NewReportHandler 创建 ReportHandler 实例。
func NewReportHandler(
	hotelSummaries repo.HotelSummaryRepository,
	roomDetails repo.RoomDetailsRepository,
	orderFull repo.OrderFullRepository,
	reviewFull repo.ReviewFullRepository,
	userVip repo.UserVipRepository,
	personInfo repo.PersonInfoRepository,
	guestBookingStats repo.GuestBookingStatsRepository,
	myOrders repo.MyOrdersRepository,
) *ReportHandler {
	return &ReportHandler{
		hotelSummaries:    hotelSummaries,
		roomDetails:       roomDetails,
		orderFull:         orderFull,
		reviewFull:        reviewFull,
		userVip:           userVip,
		personInfo:        personInfo,
		guestBookingStats: guestBookingStats,
		myOrders:          myOrders,
	}
}

// parsePageQuery 解析并校验分页参数，返回 page、pageSize、offset。
func parsePageQuery(c fiber.Ctx) (page, pageSize, offset int, err error) {
	var pq PageQuery
	if err = c.Bind().Query(&pq); err != nil {
		return 0, 0, 0, err
	}
	if pq.Page <= 0 {
		pq.Page = 1
	}
	if pq.PageSize <= 0 {
		pq.PageSize = 10
	}
	offset = (pq.Page - 1) * pq.PageSize
	return pq.Page, pq.PageSize, offset, nil
}

// parseOptionalInt16 从查询参数解析可选的 int16 值。
func parseOptionalInt16(c fiber.Ctx, key string) *int16 {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 16)
	if err != nil || v <= 0 {
		return nil
	}
	sl := int16(v)
	return &sl
}

// parseOptionalFloat64 从查询参数解析可选的 float64 值。
func parseOptionalFloat64(c fiber.Ctx, key string) *float64 {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v <= 0 {
		return nil
	}
	return &v
}

// parseOptionalInt 从查询参数解析可选的 int 值。
func parseOptionalInt(c fiber.Ctx, key string) *int {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return nil
	}
	return &v
}

// ===================== HotelSummaries =====================

// HotelSummaries 查询酒店摘要视图列表，支持按地区、星级、价格筛选。
func (h *ReportHandler) HotelSummaries(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	province := c.Query("province")
	city := c.Query("city")
	district := c.Query("district")
	starLevel := parseOptionalInt16(c, "starLevel")
	minPrice := parseOptionalFloat64(c, "minPrice")
	maxPrice := parseOptionalFloat64(c, "maxPrice")

	results, total, err := h.hotelSummaries.FindAll(ctx, offset, pageSize, province, city, district, starLevel, minPrice, maxPrice)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// ===================== RoomDetails =====================

// RoomDetails 查询客房详情视图列表，支持按地区、星级、价格筛选。
func (h *ReportHandler) RoomDetails(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	province := c.Query("province")
	city := c.Query("city")
	district := c.Query("district")
	starLevel := parseOptionalInt16(c, "starLevel")
	minPrice := parseOptionalFloat64(c, "minPrice")
	maxPrice := parseOptionalFloat64(c, "maxPrice")

	results, total, err := h.roomDetails.FindAll(ctx, offset, pageSize, province, city, district, starLevel, minPrice, maxPrice)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// RoomDetailsByHotelID 根据酒店 ID 查询客房详情视图列表。
func (h *ReportHandler) RoomDetailsByHotelID(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	hotelIDStr := c.Query("hotelID")
	if hotelIDStr == "" {
		return model.SendError(c, http.StatusBadRequest, "hotelID is required")
	}
	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	results, total, err := h.roomDetails.FindByHotelID(ctx, hotelID, offset, pageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// ===================== OrderFull =====================

// OrderFullByUserID 根据用户 ID 查询订单完整视图列表。
func (h *ReportHandler) OrderFullByUserID(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	userIDStr := c.Query("userID")
	if userIDStr == "" {
		return model.SendError(c, http.StatusBadRequest, "userID is required")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid user ID")
	}

	results, total, err := h.orderFull.FindByUserID(ctx, userID, offset, pageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// OrderFullByHotelID 根据酒店 ID 查询订单完整视图列表。
func (h *ReportHandler) OrderFullByHotelID(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	hotelIDStr := c.Query("hotelID")
	if hotelIDStr == "" {
		return model.SendError(c, http.StatusBadRequest, "hotelID is required")
	}
	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	results, total, err := h.orderFull.FindByHotelID(ctx, hotelID, offset, pageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// ===================== ReviewFull =====================

// ReviewFullByHotelID 根据酒店 ID 查询评价完整视图列表。
func (h *ReportHandler) ReviewFullByHotelID(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	hotelIDStr := c.Query("hotelID")
	if hotelIDStr == "" {
		return model.SendError(c, http.StatusBadRequest, "hotelID is required")
	}
	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	results, total, err := h.reviewFull.FindByHotelID(ctx, hotelID, offset, pageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// ReviewFullByUserID 根据用户 ID 查询评价完整视图列表。
func (h *ReportHandler) ReviewFullByUserID(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	userIDStr := c.Query("userID")
	if userIDStr == "" {
		return model.SendError(c, http.StatusBadRequest, "userID is required")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid user ID")
	}

	results, total, err := h.reviewFull.FindByUserID(ctx, userID, offset, pageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// ===================== UserVip =====================

// UserVipList 查询用户 VIP 视图列表，支持按角色筛选。
func (h *ReportHandler) UserVipList(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	role := c.Query("role")

	results, total, err := h.userVip.FindAll(ctx, offset, pageSize, role)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// ===================== PersonInfo =====================

// PersonInfoList 查询人员信息视图列表，支持按性别、年龄范围筛选。
func (h *ReportHandler) PersonInfoList(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	gender := c.Query("gender")
	minAge := parseOptionalInt(c, "minAge")
	maxAge := parseOptionalInt(c, "maxAge")

	results, total, err := h.personInfo.FindAll(ctx, offset, pageSize, gender, minAge, maxAge)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// ===================== GuestBookingStats =====================

// GuestStats 查询客人预订统计视图列表，支持按年龄组、性别、偏好城市筛选。
func (h *ReportHandler) GuestStats(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	ageGroup := c.Query("ageGroup")
	gender := c.Query("gender")
	favCity := c.Query("favCity")

	results, total, err := h.guestBookingStats.FindAll(ctx, offset, pageSize, ageGroup, gender, favCity)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}

// TopGuests 查询消费最高的客人列表。
func (h *ReportHandler) TopGuests(c fiber.Ctx) error {
	ctx := c.Context()

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	results, err := h.guestBookingStats.FindTopGuests(ctx, limit)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results))
}

// ===================== MyOrders =====================

// MyOrders 根据用户 ID 查询我的订单视图列表，支持按状态筛选。
func (h *ReportHandler) MyOrders(c fiber.Ctx) error {
	ctx := c.Context()
	page, pageSize, offset, err := parsePageQuery(c)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}

	userIDStr := c.Query("userID")
	if userIDStr == "" {
		return model.SendError(c, http.StatusBadRequest, "userID is required")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid user ID")
	}

	status := c.Query("status")

	var results []view.MyOrders
	var total int64
	if status != "" {
		results, total, err = h.myOrders.FindByUserIDAndStatus(ctx, userID, status, offset, pageSize)
	} else {
		results, total, err = h.myOrders.FindByUserID(ctx, userID, offset, pageSize)
	}
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(results), model.WithPagination(total, page, pageSize))
}
