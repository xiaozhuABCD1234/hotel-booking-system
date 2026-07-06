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
//
//	@Summary		查询全部评价
//	@Description	分页查询全部评价列表
//	@Tags			reviews
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]modelSchema.Review}
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reviews [get]
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
//
//	@Summary		查询评价详情
//	@Description	根据 UUID 查询单个评价信息（含关联用户、酒店、订单）
//	@Tags			reviews
//	@Produce		json
//	@Param			id		path		string	true	"评价 ID (UUID)"
//	@Success		200		{object}	model.Response{data=modelSchema.Review}
//	@Failure		400		{object}	model.Response	"无效的评价 ID"
//	@Failure		404		{object}	model.Response	"评价不存在"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reviews/{id} [get]
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
//
//	@Summary		创建评价
//	@Description	创建新评价记录
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			body	body		modelSchema.Review	true	"评价信息"
//	@Success		201		{object}	model.Response{data=modelSchema.Review}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reviews [post]
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

// updateReviewInput 评价更新请求体（仅允许修改评分和内容）。
type updateReviewInput struct {
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

// Update 更新评价（仅 Rating 和 Content）。
//
//	@Summary		更新评价
//	@Description	根据 UUID 更新评价的评分和内容
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"评价 ID (UUID)"
//	@Param			body	body		updateReviewInput		true	"评分和内容"
//	@Success		200		{object}	model.Response{data=modelSchema.Review}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reviews/{id} [put]
func (h *ReviewHandler) Update(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid review ID")
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
//
//	@Summary		删除评价
//	@Description	根据 UUID 硬删除评价
//	@Tags			reviews
//	@Produce		json
//	@Param			id		path		string	true	"评价 ID (UUID)"
//	@Success		200		{object}	model.Response
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reviews/{id} [delete]
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
//
//	@Summary		按酒店查询评价
//	@Description	根据酒店 ID 分页查询该酒店的全部评价
//	@Tags			reviews
//	@Produce		json
//	@Param			hotelID		query		string	true	"酒店 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]modelSchema.Review}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reviews/by-hotel [get]
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
//
//	@Summary		按用户查询评价
//	@Description	根据用户 ID 分页查询该用户的全部评价
//	@Tags			reviews
//	@Produce		json
//	@Param			userID		query		string	true	"用户 ID (UUID)"
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Success		200			{object}	model.Response{data=[]modelSchema.Review}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/reviews/by-user [get]
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
