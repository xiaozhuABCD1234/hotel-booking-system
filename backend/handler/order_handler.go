package handler

import (
	"errors"
	"net/http"
	"time"

	appmodel "backend/model"
	model "backend/model/schema"
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

// List 查询全部订单（分页）
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

// GetByID 根据 ID 查询订单详情
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

// Create 创建订单（事务，含入住人）
func (h *OrderHandler) Create(c fiber.Ctx) error {
	var order model.Order
	if err := c.Bind().Body(&order); err != nil {
		return appmodel.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	now := time.Now()
	if order.CreateAt.IsZero() {
		order.CreateAt = now
	}
	if order.UpdateAt.IsZero() {
		order.UpdateAt = now
	}

	if err := h.orders.Create(c.Context(), &order); err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(appmodel.Response{
		Success:   true,
		Data:      order,
		Timestamp: time.Now(),
	})
}

// UpdateStatus 更新订单状态
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

// Delete 硬删除订单
func (h *OrderHandler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order id")
	}

	if err := h.orders.Delete(c.Context(), id); err != nil {
		return err
	}

	return appmodel.SendSuccess(c, appmodel.WithMessage("Order deleted"))
}

// userIDQuery 按用户 ID 查询参数
type userIDQuery struct {
	UserID uuid.UUID `query:"userID"`
	paginationQuery
}

// ListByUserID 根据用户 ID 查询订单列表（分页）
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

// ListByHotelID 根据酒店 ID 查询订单列表（分页）
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
