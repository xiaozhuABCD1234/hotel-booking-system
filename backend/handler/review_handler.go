package handler

import (
	"net/http"
	"time"

	"backend/model"
	modelSchema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// ReviewHandler 评价 HTTP 处理器，封装 ReviewRepo 的 CRUD 操作。
type ReviewHandler struct {
	reviews repo.ReviewRepository
}

// NewReviewHandler 创建 ReviewHandler 实例。
func NewReviewHandler(reviewRepo repo.ReviewRepository) *ReviewHandler {
	return &ReviewHandler{reviews: reviewRepo}
}

// List 查询全部评价列表（分页）。
// Query: page (default 1), pageSize (default 10)
func (h *ReviewHandler) List(c fiber.Ctx) error {
	ctx := c.Context()

	var q struct {
		Page     int `query:"page"`
		PageSize int `query:"pageSize"`
	}
	if err := c.Bind().Query(&q); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid query parameters")
	}

	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	offset := (q.Page - 1) * q.PageSize
	reviews, total, err := h.reviews.FindAll(ctx, offset, q.PageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c,
		model.WithData(reviews),
		model.WithPagination(total, q.Page, q.PageSize),
	)
}

// GetByID 根据评价 ID 查询详情，预加载 User、Hotel、Order。
func (h *ReviewHandler) GetByID(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid review ID")
	}

	review, err := h.reviews.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(review))
}

// Create 创建评价，返回 201 Created。
func (h *ReviewHandler) Create(c fiber.Ctx) error {
	ctx := c.Context()

	var review modelSchema.Review
	if err := c.Bind().Body(&review); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.reviews.Create(ctx, &review); err != nil {
		return err
	}

	// SendSuccess 固定返回 200，Create 需要 201，因此手动构造同格式响应
	r := model.Response{
		Success:   true,
		Data:      review,
		Timestamp: time.Now(),
	}
	return c.Status(http.StatusCreated).JSON(r)
}

// Update 更新评价（仅 Rating 和 Content）。
func (h *ReviewHandler) Update(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid review ID")
	}

	type updateReviewInput struct {
		Rating  int    `json:"rating"`
		Content string `json:"content"`
	}

	var input updateReviewInput
	if err := c.Bind().Body(&input); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	review, err := h.reviews.FindByID(ctx, id)
	if err != nil {
		return err
	}

	review.Rating = int16(input.Rating)
	review.Content = &input.Content

	if err := h.reviews.Update(ctx, review); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(review))
}

// Delete 根据 ID 硬删除评价。
func (h *ReviewHandler) Delete(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid review ID")
	}

	if err := h.reviews.Delete(ctx, id); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithMessage("Review deleted successfully"))
}

// ListByHotelID 根据酒店 ID 查询评价列表（分页）。
// Query: hotelID (required), page (default 1), pageSize (default 10)
func (h *ReviewHandler) ListByHotelID(c fiber.Ctx) error {
	ctx := c.Context()

	var q struct {
		HotelID  string `query:"hotelID"`
		Page     int    `query:"page"`
		PageSize int    `query:"pageSize"`
	}
	if err := c.Bind().Query(&q); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid query parameters")
	}

	hotelID, err := uuid.Parse(q.HotelID)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	offset := (q.Page - 1) * q.PageSize
	reviews, total, err := h.reviews.FindByHotelID(ctx, hotelID, offset, q.PageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c,
		model.WithData(reviews),
		model.WithPagination(total, q.Page, q.PageSize),
	)
}

// ListByUserID 根据用户 ID 查询评价列表（分页）。
// Query: userID (required), page (default 1), pageSize (default 10)
func (h *ReviewHandler) ListByUserID(c fiber.Ctx) error {
	ctx := c.Context()

	var q struct {
		UserID   string `query:"userID"`
		Page     int    `query:"page"`
		PageSize int    `query:"pageSize"`
	}
	if err := c.Bind().Query(&q); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid query parameters")
	}

	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid user ID")
	}

	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	offset := (q.Page - 1) * q.PageSize
	reviews, total, err := h.reviews.FindByUserID(ctx, userID, offset, q.PageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c,
		model.WithData(reviews),
		model.WithPagination(total, q.Page, q.PageSize),
	)
}
