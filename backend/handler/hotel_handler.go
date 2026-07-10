package handler

import (
	"net/http"
	"strconv"
	"time"

	"backend/model"
	schema "backend/model/schema"
	"backend/repo"
	"backend/service"

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
	hotels      repo.HotelRepository
	rooms       repo.RoomRepository
	hotelImages *repo.HotelImageRepo
	cos         *service.COSService
}

// NewHotelHandler 创建 HotelHandler 实例
func NewHotelHandler(hotelRepo repo.HotelRepository, roomRepo repo.RoomRepository, hotelImageRepo *repo.HotelImageRepo, cosSvc *service.COSService) *HotelHandler {
	return &HotelHandler{
		hotels:      hotelRepo,
		rooms:       roomRepo,
		hotelImages: hotelImageRepo,
		cos:         cosSvc,
	}
}

// List 查询酒店列表（支持按区域、星级、关键词、价格范围、入住日期筛选）。
//
//	@Summary		查询酒店列表
//	@Description	分页查询酒店列表，支持按区域、星级、关键词、价格范围、入住日期筛选
//	@Tags			hotels
//	@Produce		json
//	@Param			page			query		int		false	"页码"			default(1)
//	@Param			pageSize		query		int		false	"每页数量"		default(10)
//	@Param			regionID		query		int		false	"区域 ID"
//	@Param			starLevel		query		int		false	"星级"
//	@Param			keyword			query		string	false	"搜索关键词"
//	@Param			minPrice		query		number	false	"最低价格"
//	@Param			maxPrice		query		number	false	"最高价格"
//	@Param			checkInDate		query		string	false	"入住日期 (YYYY-MM-DD)"
//	@Param			checkOutDate	query		string	false	"离店日期 (YYYY-MM-DD)"
//	@Success		200				{object}	model.Response{data=[]schema.Hotel}
//	@Failure		400				{object}	model.Response
//	@Failure		500				{object}	model.Response
//	@Router			/hotels [get]
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
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")
	checkInDateStr := c.Query("checkInDate")
	checkOutDateStr := c.Query("checkOutDate")

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

	var minPrice, maxPrice *float64
	if minPriceStr != "" {
		v, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil || v < 0 {
			return model.SendError(c, http.StatusBadRequest, "minPrice must be a non-negative number")
		}
		minPrice = &v
	}
	if maxPriceStr != "" {
		v, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil || v < 0 {
			return model.SendError(c, http.StatusBadRequest, "maxPrice must be a non-negative number")
		}
		maxPrice = &v
	}
	if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
		return model.SendError(c, http.StatusBadRequest, "minPrice cannot exceed maxPrice")
	}

	var checkInDate, checkOutDate *time.Time
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if checkInDateStr != "" {
		v, err := time.ParseInLocation("2006-01-02", checkInDateStr, time.Local)
		if err != nil {
			return model.SendError(c, http.StatusBadRequest, "checkInDate must be in YYYY-MM-DD format")
		}
		if v.Before(today) {
			return model.SendError(c, http.StatusBadRequest, "checkInDate cannot be in the past")
		}
		checkInDate = &v
	}
	if checkOutDateStr != "" {
		v, err := time.ParseInLocation("2006-01-02", checkOutDateStr, time.Local)
		if err != nil {
			return model.SendError(c, http.StatusBadRequest, "checkOutDate must be in YYYY-MM-DD format")
		}
		if v.Before(today) {
			return model.SendError(c, http.StatusBadRequest, "checkOutDate cannot be in the past")
		}
		checkOutDate = &v
	}
	// Require both dates together — single date is silently useless
	if (checkInDate != nil) != (checkOutDate != nil) {
		return model.SendError(c, http.StatusBadRequest, "Both checkInDate and checkOutDate are required for date filtering")
	}
	if checkInDate != nil && checkOutDate != nil && !checkInDate.Before(*checkOutDate) {
		return model.SendError(c, http.StatusBadRequest, "checkInDate must be before checkOutDate")
	}

	hotels, total, err := h.hotels.FindAll(ctx, offset, pq.PageSize, regionID, starLevel, keyword, minPrice, maxPrice, checkInDate, checkOutDate)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(hotels), model.WithPagination(total, pq.Page, pq.PageSize))
}

// GetByID 根据 ID 查询酒店详情。
//
//	@Summary		查询酒店详情
//	@Description	根据 UUID 查询单个酒店信息
//	@Tags			hotels
//	@Produce		json
//	@Param			id		path		string	true	"酒店 ID (UUID)"
//	@Success		200		{object}	model.Response{data=schema.Hotel}
//	@Failure		400		{object}	model.Response	"无效的酒店 ID"
//	@Failure		404		{object}	model.Response	"酒店不存在"
//	@Failure		500		{object}	model.Response
//	@Router			/hotels/{id} [get]
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

// Create 创建酒店，返回 201 Created。
//
//	@Summary		创建酒店
//	@Description	创建新酒店记录
//	@Tags			hotels
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.Hotel	true	"酒店信息"
//	@Success		201		{object}	model.Response{data=schema.Hotel}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/hotels [post]
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

// updateHotelInput 仅包含客户端可更新的酒店字段。
type updateHotelInput struct {
	HotelName   string  `json:"hotelName,omitempty"`
	RegionID    *int    `json:"regionID,omitempty"`
	Address     string  `json:"address,omitempty"`
	Telephone   string  `json:"telephone,omitempty"`
	StarLevel   *int16  `json:"starLevel,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Update 更新酒店信息。
//
//	@Summary		更新酒店信息
//	@Description	根据 UUID 更新酒店信息
//	@Tags			hotels
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"酒店 ID (UUID)"
//	@Param			body	body		schema.Hotel	true	"更新的酒店信息"
//	@Success		200		{object}	model.Response{data=schema.Hotel}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/hotels/{id} [put]
func (h *HotelHandler) Update(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	existing, err := h.hotels.FindByID(ctx, id)
	if err != nil {
		return err
	}

	var input updateHotelInput
	if err := c.Bind().Body(&input); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if input.HotelName != "" {
		existing.HotelName = input.HotelName
	}
	if input.RegionID != nil {
		existing.RegionID = *input.RegionID
	}
	if input.Address != "" {
		existing.Address = input.Address
	}
	if input.Telephone != "" {
		existing.Telephone = input.Telephone
	}
	if input.StarLevel != nil {
		existing.StarLevel = input.StarLevel
	}
	if input.Description != nil {
		existing.Description = input.Description
	}

	if err := h.hotels.Update(ctx, existing); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(existing))
}

// Delete 删除酒店（软删除）。
//
//	@Summary		删除酒店
//	@Description	根据 UUID 软删除酒店
//	@Tags			hotels
//	@Produce		json
//	@Param			id		path		string	true	"酒店 ID (UUID)"
//	@Success		200		{object}	model.Response
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/hotels/{id} [delete]
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

// ListRooms 查询客房列表（支持按酒店筛选）。
//
//	@Summary		查询客房列表
//	@Description	分页查询客房列表，支持按酒店 ID 筛选
//	@Tags			rooms
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			hotelID		query		string	false	"酒店 ID (UUID)"
//	@Success		200			{object}	model.Response{data=[]schema.Room}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Router			/rooms [get]
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

// GetRoomByID 根据 ID 查询客房详情。
//
//	@Summary		查询客房详情
//	@Description	根据 UUID 查询单个客房信息
//	@Tags			rooms
//	@Produce		json
//	@Param			id		path		string	true	"客房 ID (UUID)"
//	@Success		200		{object}	model.Response{data=schema.Room}
//	@Failure		400		{object}	model.Response	"无效的客房 ID"
//	@Failure		404		{object}	model.Response	"客房不存在"
//	@Failure		500		{object}	model.Response
//	@Router			/rooms/{id} [get]
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

// CreateRoom 创建客房，返回 201 Created。
//
//	@Summary		创建客房
//	@Description	为指定酒店创建新房型
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.Room	true	"客房信息"
//	@Success		201		{object}	model.Response{data=schema.Room}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/rooms [post]
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

// updateRoomInput 仅包含客户端可更新的客房字段。
type updateRoomInput struct {
	HotelID           *uuid.UUID `json:"hotelID,omitempty"`
	TypeName          string     `json:"typeName,omitempty"`
	TotalQuantity     *int32     `json:"totalQuantity,omitempty"`
	AvailableQuantity *int32     `json:"availableQuantity,omitempty"`
	Price             *float64   `json:"price,omitempty"`
	WeekendPrice      *float64   `json:"weekendPrice,omitempty"`
	Description       *string    `json:"description,omitempty"`
}

// UpdateRoom 更新客房信息。
//
//	@Summary		更新客房信息
//	@Description	根据 UUID 更新客房信息
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"客房 ID (UUID)"
//	@Param			body	body		schema.Room	true	"更新的客房信息"
//	@Success		200		{object}	model.Response{data=schema.Room}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/rooms/{id} [put]
func (h *HotelHandler) UpdateRoom(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid room ID")
	}

	existing, err := h.rooms.FindByID(ctx, id)
	if err != nil {
		return err
	}

	var input updateRoomInput
	if err := c.Bind().Body(&input); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if input.HotelID != nil {
		existing.HotelID = *input.HotelID
	}
	if input.TypeName != "" {
		existing.TypeName = input.TypeName
	}
	if input.TotalQuantity != nil {
		existing.TotalQuantity = *input.TotalQuantity
	}
	if input.AvailableQuantity != nil {
		existing.AvailableQuantity = *input.AvailableQuantity
	}
	if input.Price != nil {
		existing.Price = *input.Price
	}
	if input.WeekendPrice != nil {
		existing.WeekendPrice = input.WeekendPrice
	}
	if input.Description != nil {
		existing.Description = input.Description
	}

	if err := h.rooms.Update(ctx, existing); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(existing))
}

// DeleteRoom 删除客房（软删除）。
//
//	@Summary		删除客房
//	@Description	根据 UUID 软删除客房
//	@Tags			rooms
//	@Produce		json
//	@Param			id		path		string	true	"客房 ID (UUID)"
//	@Success		200		{object}	model.Response
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/rooms/{id} [delete]
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

// UploadImage 上传酒店图片到 COS 并保存记录。
//
//	@Summary		上传酒店图片
//	@Description	上传图片到腾讯云 COS 并关联到指定酒店。
//	@Tags			hotels
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string	true	"酒店 ID (UUID)"
//	@Param			file	formData	file	true	"图片文件"
//	@Success		200		{object}	model.Response{data=object{url=string,key=string}}
//	@Failure		400		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/hotels/{id}/images [post]
func (h *HotelHandler) UploadImage(c fiber.Ctx) error {
	ctx := c.Context()

	hotelID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "file is required")
	}

	result, err := h.cos.Upload(ctx, file)
	if err != nil {
		return model.SendError(c, http.StatusInternalServerError, "upload failed: "+err.Error())
	}

	if err := h.hotelImages.Create(ctx, &schema.HotelImage{
		HotelID:  hotelID,
		ImageURL: result.URL,
	}); err != nil {
		return model.SendError(c, http.StatusInternalServerError, "failed to save image record: "+err.Error())
	}

	return model.SendSuccess(c, model.WithData(fiber.Map{
		"url": result.URL,
		"key": result.Key,
	}))
}

// DeleteImage 删除酒店图片（从 COS 和数据库中同时删除）。
//
//	@Summary		删除酒店图片
//	@Description	根据酒店 ID 和图片 URL 删除酒店图片（同时删除 COS 对象和数据库记录）。
//	@Tags			hotels
//	@Produce		json
//	@Param			id			path		string	true	"酒店 ID (UUID)"
//	@Param			imageUrl	query		string	true	"图片 URL"
//	@Success		200			{object}	model.Response
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/hotels/{id}/images [delete]
func (h *HotelHandler) DeleteImage(c fiber.Ctx) error {
	ctx := c.Context()

	hotelID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid hotel ID")
	}

	imageURL := c.Query("imageUrl")
	if imageURL == "" {
		return model.SendError(c, http.StatusBadRequest, "imageUrl is required")
	}

	// 先删 COS 对象
	key := service.KeyFromURL(imageURL)
	if key != "" {
		if err := h.cos.Delete(ctx, key); err != nil {
			return model.SendError(c, http.StatusInternalServerError, "failed to delete COS object: "+err.Error())
		}
	}

	// 再删数据库记录
	if err := h.hotelImages.Delete(ctx, hotelID, imageURL); err != nil {
		return model.SendError(c, http.StatusInternalServerError, "failed to delete image record: "+err.Error())
	}

	return model.SendSuccess(c, model.WithMessage("Image deleted"))
}
