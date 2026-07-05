package handler

import (
	"net/http"
	"strconv"
	"time"

	"backend/model"
	schema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// PageQuery 分页查询参数
type PageQuery struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

// HotelHandler 酒店与客房 HTTP 处理器
type HotelHandler struct {
	hotels repo.HotelRepository
	rooms  repo.RoomRepository
}

// NewHotelHandler 创建 HotelHandler 实例
func NewHotelHandler(hotelRepo repo.HotelRepository, roomRepo repo.RoomRepository) *HotelHandler {
	return &HotelHandler{
		hotels: hotelRepo,
		rooms:  roomRepo,
	}
}

// List 查询酒店列表
func (h *HotelHandler) List(c fiber.Ctx) error {
	ctx := c.Context()

	var pq PageQuery
	if err := c.Bind().Query(&pq); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid query parameters")
	}
	if pq.Page <= 0 {
		pq.Page = 1
	}
	if pq.PageSize <= 0 {
		pq.PageSize = 10
	}
	offset := (pq.Page - 1) * pq.PageSize

	regionIDStr := c.Query("regionID")
	starLevelStr := c.Query("starLevel")
	keyword := c.Query("keyword")

	var regionID *int
	if regionIDStr != "" {
		if v, err := strconv.Atoi(regionIDStr); err == nil && v > 0 {
			regionID = &v
		}
	}

	var starLevel *int16
	if starLevelStr != "" {
		if v, err := strconv.ParseInt(starLevelStr, 10, 16); err == nil && v > 0 {
			sl := int16(v)
			starLevel = &sl
		}
	}

	hotels, total, err := h.hotels.FindAll(ctx, offset, pq.PageSize, regionID, starLevel, keyword)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(hotels), model.WithPagination(total, pq.Page, pq.PageSize))
}

// GetByID 根据 ID 查询酒店
func (h *HotelHandler) GetByID(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	hotel, err := h.hotels.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(hotel))
}

// Create 创建酒店
func (h *HotelHandler) Create(c fiber.Ctx) error {
	ctx := c.Context()

	var hotel schema.Hotel
	if err := c.Bind().Body(&hotel); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.hotels.Create(ctx, &hotel); err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(model.Response{
		Success:   true,
		Data:      hotel,
		Timestamp: time.Now(),
	})
}

// Update 更新酒店
func (h *HotelHandler) Update(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	var hotel schema.Hotel
	if err := c.Bind().Body(&hotel); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	hotel.ID = id

	if err := h.hotels.Update(ctx, &hotel); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(hotel))
}

// Delete 删除酒店（软删除）
func (h *HotelHandler) Delete(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	if err := h.hotels.Delete(ctx, id); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithMessage("Hotel deleted"))
}

// ListRooms 查询客房列表
func (h *HotelHandler) ListRooms(c fiber.Ctx) error {
	ctx := c.Context()

	var pq PageQuery
	if err := c.Bind().Query(&pq); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid query parameters")
	}
	if pq.Page <= 0 {
		pq.Page = 1
	}
	if pq.PageSize <= 0 {
		pq.PageSize = 10
	}
	offset := (pq.Page - 1) * pq.PageSize

	hotelIDStr := c.Query("hotelID")
	if hotelIDStr != "" {
		hotelID, err := uuid.Parse(hotelIDStr)
		if err != nil {
			return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
		}

		rooms, total, err := h.rooms.FindByHotelID(ctx, hotelID, offset, pq.PageSize)
		if err != nil {
			return err
		}

		return model.SendSuccess(c, model.WithData(rooms), model.WithPagination(total, pq.Page, pq.PageSize))
	}

	rooms, total, err := h.rooms.FindAll(ctx, offset, pq.PageSize)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(rooms), model.WithPagination(total, pq.Page, pq.PageSize))
}

// GetRoomByID 根据 ID 查询客房
func (h *HotelHandler) GetRoomByID(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid room ID")
	}

	room, err := h.rooms.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(room))
}

// CreateRoom 创建客房
func (h *HotelHandler) CreateRoom(c fiber.Ctx) error {
	ctx := c.Context()

	var room schema.Room
	if err := c.Bind().Body(&room); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.rooms.Create(ctx, &room); err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(model.Response{
		Success:   true,
		Data:      room,
		Timestamp: time.Now(),
	})
}

// UpdateRoom 更新客房
func (h *HotelHandler) UpdateRoom(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid room ID")
	}

	var room schema.Room
	if err := c.Bind().Body(&room); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	room.ID = id

	if err := h.rooms.Update(ctx, &room); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(room))
}

// DeleteRoom 删除客房（软删除）
func (h *HotelHandler) DeleteRoom(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid room ID")
	}

	if err := h.rooms.Delete(ctx, id); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithMessage("Room deleted"))
}
