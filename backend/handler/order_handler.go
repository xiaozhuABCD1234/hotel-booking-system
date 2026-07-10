package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"backend/middleware"
	appmodel "backend/model"
	model "backend/model/schema"
	_ "backend/model/view"
	"backend/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// OrderHandler 订单 HTTP 处理器
type OrderHandler struct {
	orders *service.OrderService
}

// NewOrderHandler 构造 OrderHandler
func NewOrderHandler(orderSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{orders: orderSvc}
}

// paginationQuery 分页查询参数
type paginationQuery struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

// statusBody 状态更新请求体
type statusBody struct {
	Status model.OrderStatus `json:"status"`
}

// isFKeyViolation 判断 GORM 错误是否为外键约束违反。
func isFKeyViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "violates foreign key constraint")
}

// List 查询全部订单（分页）。
//
// @Summary		查询全部订单
// @Description	分页查询全部订单列表
// @Tags			orders
// @Produce		json
// @Param			page		query		int		false	"页码"			default(1)
// @Param			pageSize	query		int		false	"每页数量"		default(10)
// @Success		200			{object}	model.Response{data=[]view.OrderSummary}
// @Failure		500			{object}	model.Response
// @Security		BearerAuth
// @Router			/orders [get]
func (h *OrderHandler) List(c fiber.Ctx) error {
	var q paginationQuery
	if err := c.Bind().Query(&q); err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid request query")
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	offset := (q.Page - 1) * q.PageSize

	orders, total, err := h.orders.ListAll(c.Context(), offset, q.PageSize)
	if err != nil {
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithData(orders), appmodel.WithPagination(total, q.Page, q.PageSize))
}

// GetByID 根据 ID 查询订单详情。
//
//	@Summary		查询订单详情
//	@Description	根据 UUID 查询单个订单信息
//	@Tags			orders
//	@Produce		json
//	@Param			id		path		string	true	"订单 ID (UUID)"
//	@Success		200		{object}	model.Response{data=model.Order}
//	@Failure		400		{object}	model.Response	"无效的订单 ID"
//	@Failure		404		{object}	model.Response	"订单不存在"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/orders/{id} [get]
func (h *OrderHandler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order id")
	}

	order, err := h.orders.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithData(order))
}

// createOrderRequest 下单请求 DTO（匹配前端 BookingView 提交的格式）。
// 后端负责从 JWT 提取 userId、映射字段名、构造 guests 数组。
type createOrderRequest struct {
	RoomID       uuid.UUID `json:"roomId"`
	CheckInDate  string    `json:"checkInDate"`
	CheckOutDate string    `json:"checkOutDate"`
	GuestName    string    `json:"guestName"`
	GuestPhone   string    `json:"guestPhone"`
	GuestIDCard  string    `json:"guestIdCard"`
	RoomCount    int32     `json:"roomCount"`
	TotalPrice   float64   `json:"totalPrice"`
	ActualPrice  float64   `json:"actualPrice"`
}

// Create 创建订单（事务，含入住人），返回 201 Created。
//
//	@Summary		创建订单
//	@Description	创建新订单（含入住人信息），返回 201 Created
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createOrderRequest	true	"订单信息"
//	@Success		201		{object}	model.Response{data=model.Order}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		401		{object}	model.Response	"未认证"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/orders [post]
func (h *OrderHandler) Create(c fiber.Ctx) error {
	// 从 JWT 提取 userId
	claims := middleware.GetClaims(c)
	if claims == nil {
		return appmodel.SendError(c, http.StatusUnauthorized, "authentication required")
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "invalid user ID in token")
	}

	// 绑定前端格式的请求体
	var req createOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	// 解析日期
	checkInDate, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid checkInDate format (expected YYYY-MM-DD)")
	}
	checkOutDate, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid checkOutDate format (expected YYYY-MM-DD)")
	}

	now := time.Now()

	// 构造 Order 模型（VIP 折扣由 trg_order_vip_discount_1718 触发器自动计算）
	order := model.Order{
		UserID:       userID,
		RoomID:       req.RoomID,
		Quantity:     req.RoomCount,
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
		TotalPrice:   req.TotalPrice,
		ActualPrice:  req.ActualPrice,
		Status:       model.OrderPending,
		CreateAt:     now,
		UpdateAt:     now,
	}

	// 构造入住人（前端一个表单对应一个住客）
	if req.GuestName != "" || req.GuestIDCard != "" {
		if req.GuestName == "" || req.GuestIDCard == "" {
			return appmodel.SendError(c, http.StatusBadRequest, "guestName and guestIdCard are both required")
		}
		order.Guests = []model.OrderGuest{
			{
				IDCard: req.GuestIDCard,
				Person: model.Person{
					IDCard: req.GuestIDCard,
					Name:   req.GuestName,
					Phone:  &req.GuestPhone,
				},
			},
		}
	}

	if err := h.orders.Create(c.Context(), &order); err != nil {
		if isFKeyViolation(err) {
			return appmodel.SendError(c, http.StatusBadRequest, "Referenced user or room not found")
		}
		// Service 层校验错误（如无效状态、日期范围等）返回 400
		return appmodel.SendError(c, http.StatusBadRequest, err.Error())
	}

	return c.Status(http.StatusCreated).JSON(appmodel.Response{
		Success:   true,
		Data:      order,
		Timestamp: time.Now(),
	})
}

// UpdateStatus 更新订单状态。
//
//	@Summary		更新订单状态
//	@Description	根据 UUID 更新订单状态（如已预订→已入住）
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"订单 ID (UUID)"
//	@Param			body	body		statusBody	true	"状态更新请求"
//	@Success		200		{object}	model.Response{data=model.Order}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		409		{object}	model.Response	"非法状态转换"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/orders/{id}/status [put]
func (h *OrderHandler) UpdateStatus(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order id")
	}

	var body statusBody
	if err := c.Bind().Body(&body); err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.orders.UpdateStatus(c.Context(), id, body.Status); err != nil {
		if errors.Is(err, service.ErrInvalidTransition) {
			return fiber.NewError(fiber.StatusConflict, "invalid status transition")
		}
		return err
	}

	order, err := h.orders.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithData(order))
}

// Delete 硬删除订单。
//
//	@Summary		删除订单
//	@Description	根据 UUID 硬删除订单
//	@Tags			orders
//	@Produce		json
//	@Param			id		path		string	true	"订单 ID (UUID)"
//	@Success		200		{object}	model.Response
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/orders/{id} [delete]
func (h *OrderHandler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order id")
	}

	if err := h.orders.Delete(c.Context(), id); err != nil {
		if isFKeyViolation(err) {
			return appmodel.SendError(c, http.StatusConflict, "Cannot delete order with existing reviews, please delete reviews first")
		}
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithMessage("Order deleted"))
}

// userIDQuery 按用户 ID 查询参数
type userIDQuery struct {
	UserID uuid.UUID `query:"userID"`
	paginationQuery
}

// ListByUserID 根据用户 ID 查询订单列表（分页）。
//
//	@Summary		按用户查询订单
//	@Description	根据用户 ID 分页查询该用户的全部订单
//	@Tags			orders
//	@Produce		json
//	@Param			userID		query		string	true	"用户 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]model.Order}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/orders/by-user [get]
func (h *OrderHandler) ListByUserID(c fiber.Ctx) error {
	var q userIDQuery
	if err := c.Bind().Query(&q); err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid request query")
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	offset := (q.Page - 1) * q.PageSize

	orders, total, err := h.orders.ListByUser(c.Context(), q.UserID, offset, q.PageSize)
	if err != nil {
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithData(orders), appmodel.WithPagination(total, q.Page, q.PageSize))
}

// hotelIDQuery 按酒店 ID 查询参数
type hotelIDQuery struct {
	HotelID uuid.UUID `query:"hotelID"`
	paginationQuery
}

// ListByHotelID 根据酒店 ID 查询订单列表（分页）。
//
//	@Summary		按酒店查询订单
//	@Description	根据酒店 ID 分页查询该酒店的全部订单
//	@Tags			orders
//	@Produce		json
//	@Param			hotelID		query		string	true	"酒店 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]model.Order}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/orders/by-hotel [get]
func (h *OrderHandler) ListByHotelID(c fiber.Ctx) error {
	var q hotelIDQuery
	if err := c.Bind().Query(&q); err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid request query")
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	offset := (q.Page - 1) * q.PageSize

	orders, total, err := h.orders.ListByHotel(c.Context(), q.HotelID, offset, q.PageSize)
	if err != nil {
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithData(orders), appmodel.WithPagination(total, q.Page, q.PageSize))
}

// Detail 查询订单完整详情（下单人与入住人明确区分，走 fn_order_detail_1718）。
//
//	@Summary		查询订单完整详情
//	@Description	根据 UUID 查询订单完整信息，下单人与入住人字段明确区分，入住人已聚合
//	@Tags			orders
//	@Produce		json
//	@Param			id		path		string	true	"订单 ID (UUID)"
//	@Success		200		{object}	model.Response{data=view.OrderDetail}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/orders/{id}/detail [get]
func (h *OrderHandler) Detail(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order id")
	}

	detail, err := h.orders.GetDetail(c.Context(), id)
	if err != nil {
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithData(detail))
}

// summaryQuery 订单概览查询参数
type summaryQuery struct {
	Status string `query:"status"`
	paginationQuery
}

// ListSummaries 查询订单概览列表（支持按状态筛选，分页）。
//
//	@Summary		查询订单概览列表
//	@Description	分页查询订单概览列表，字段精简适合管理端，支持按状态筛选
//	@Tags			orders
//	@Produce		json
//	@Param			status		query		string	false	"订单状态"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]view.OrderSummary}
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/orders/summaries [get]
func (h *OrderHandler) ListSummaries(c fiber.Ctx) error {
	var q summaryQuery
	if err := c.Bind().Query(&q); err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid request query")
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	offset := (q.Page - 1) * q.PageSize

	results, total, err := h.orders.ListSummaries(c.Context(), q.Status, offset, q.PageSize)
	if err != nil {
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithData(results), appmodel.WithPagination(total, q.Page, q.PageSize))
}
