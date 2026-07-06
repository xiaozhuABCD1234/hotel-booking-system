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
//
//	@Summary		酒店摘要报表
//	@Description	分页查询酒店摘要视图，支持按省份、城市、区县、星级、价格范围筛选
//	@Tags			reports
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			province	query		string	false	"省份"
//	@Param			city		query		string	false	"城市"
//	@Param			district	query		string	false	"区县"
//	@Param			starLevel	query		int		false	"星级"
//	@Param			minPrice	query		number	false	"最低价格"
//	@Param			maxPrice	query		number	false	"最高价格"
//	@Success		200			{object}	model.Response{data=[]view.HotelSummary}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/hotel-summaries [get]
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
//
//	@Summary		客房详情报表
//	@Description	分页查询客房详情视图，支持按省份、城市、区县、星级、价格范围筛选
//	@Tags			reports
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			province	query		string	false	"省份"
//	@Param			city		query		string	false	"城市"
//	@Param			district	query		string	false	"区县"
//	@Param			starLevel	query		int		false	"星级"
//	@Param			minPrice	query		number	false	"最低价格"
//	@Param			maxPrice	query		number	false	"最高价格"
//	@Success		200			{object}	model.Response{data=[]view.RoomDetails}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/room-details [get]
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
//
//	@Summary		按酒店查询客房详情
//	@Description	根据酒店 ID 分页查询该酒店的客房详情
//	@Tags			reports
//	@Produce		json
//	@Param			hotelID		query		string	true	"酒店 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]view.RoomDetails}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/room-details/by-hotel [get]
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
//
//	@Summary		按用户查询完整订单
//	@Description	根据用户 ID 分页查询该用户的订单完整视图（含酒店、客房、入住人等信息）
//	@Tags			reports
//	@Produce		json
//	@Param			userID		query		string	true	"用户 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]view.OrderFull}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/order-full/by-user [get]
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
//
//	@Summary		按酒店查询完整订单
//	@Description	根据酒店 ID 分页查询该酒店的订单完整视图
//	@Tags			reports
//	@Produce		json
//	@Param			hotelID		query		string	true	"酒店 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]view.OrderFull}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/order-full/by-hotel [get]
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
//
//	@Summary		按酒店查询完整评价
//	@Description	根据酒店 ID 分页查询该酒店的评价完整视图
//	@Tags			reports
//	@Produce		json
//	@Param			hotelID		query		string	true	"酒店 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]view.ReviewFull}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/review-full/by-hotel [get]
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
//
//	@Summary		按用户查询完整评价
//	@Description	根据用户 ID 分页查询该用户的评价完整视图
//	@Tags			reports
//	@Produce		json
//	@Param			userID		query		string	true	"用户 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]view.ReviewFull}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/review-full/by-user [get]
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
//
//	@Summary		用户 VIP 报表
//	@Description	分页查询用户 VIP 视图，支持按角色筛选
//	@Tags			reports
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			role		query		string	false	"角色筛选"
//	@Success		200			{object}	model.Response{data=[]view.UserVip}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/user-vip [get]
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
//
//	@Summary		人员信息报表
//	@Description	分页查询人员信息视图，支持按性别、年龄范围筛选
//	@Tags			reports
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			gender		query		string	false	"性别筛选"
//	@Param			minAge		query		int		false	"最小年龄"
//	@Param			maxAge		query		int		false	"最大年龄"
//	@Success		200			{object}	model.Response{data=[]view.PersonInfo}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/person-info [get]
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
//
//	@Summary		客人预订统计
//	@Description	分页查询客人预订统计视图，支持按年龄组、性别、偏好城市筛选
//	@Tags			reports
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			ageGroup	query		string	false	"年龄组"
//	@Param			gender		query		string	false	"性别"
//	@Param			favCity		query		string	false	"偏好城市"
//	@Success		200			{object}	model.Response{data=[]view.GuestBookingStats}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/guest-stats [get]
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
//
//	@Summary		消费排行
//	@Description	查询消费金额最高的客人列表
//	@Tags			reports
//	@Produce		json
//	@Param			limit	query		int	false	"返回数量"		default(10)
//	@Success		200		{object}	model.Response{data=[]view.GuestBookingStats}
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/guest-stats/top [get]
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
//
//	@Summary		我的订单
//	@Description	根据用户 ID 分页查询订单视图，支持按状态筛选
//	@Tags			reports
//	@Produce		json
//	@Param			userID		query		string	true	"用户 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			status		query		string	false	"订单状态筛选"
//	@Success		200			{object}	model.Response{data=[]view.MyOrders}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reports/my-orders [get]
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
