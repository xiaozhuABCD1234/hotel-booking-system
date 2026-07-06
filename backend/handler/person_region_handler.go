package handler

import (
	"net/http"
	"strconv"
	"time"

	"backend/model"
	schema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
)

// ===================== PersonHandler =====================

type PersonHandler struct {
	persons repo.PersonRepository
}

func NewPersonHandler(persons repo.PersonRepository) *PersonHandler {
	return &PersonHandler{persons: persons}
}

// List 查询人员列表（分页，支持关键词搜索）。
//
//	@Summary		查询人员列表
//	@Description	分页查询人员列表，支持按姓名、身份证号等关键词搜索
//	@Tags			persons
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			keyword		query		string	false	"搜索关键词"
//	@Success		200			{object}	model.Response{data=[]schema.Person}
//	@Failure		400			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/persons [get]
func (h *PersonHandler) List(c fiber.Ctx) error {
	ctx := c.Context()

	var pq PageQuery
	if err := c.Bind().Query(&pq); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid pagination parameters")
	}
	if pq.Page <= 0 {
		pq.Page = 1
	}
	if pq.PageSize <= 0 {
		pq.PageSize = 10
	}
	offset := (pq.Page - 1) * pq.PageSize
	keyword := c.Query("keyword", "")

	persons, total, err := h.persons.FindAll(ctx, offset, pq.PageSize, keyword)
	if err != nil {
		return err
	}

	return model.SendSuccess(c,
		model.WithData(persons),
		model.WithPagination(total, pq.Page, pq.PageSize),
	)
}

// GetByIDCard 根据身份证号查询人员详情。
//
//	@Summary		查询人员详情
//	@Description	根据身份证号查询单个人员信息
//	@Tags			persons
//	@Produce		json
//	@Param			idCard	path		string	true	"身份证号"
//	@Success		200		{object}	model.Response{data=schema.Person}
//	@Failure		404		{object}	model.Response	"人员不存在"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/persons/{idCard} [get]
func (h *PersonHandler) GetByIDCard(c fiber.Ctx) error {
	ctx := c.Context()
	idCard := c.Params("idCard")

	person, err := h.persons.FindByIDCard(ctx, idCard)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(person))
}

// Create 创建人员记录，返回 201 Created。
//
//	@Summary		创建人员
//	@Description	创建新人员记录
//	@Tags			persons
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.Person	true	"人员信息"
//	@Success		201		{object}	model.Response{data=schema.Person}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/persons [post]
func (h *PersonHandler) Create(c fiber.Ctx) error {
	ctx := c.Context()

	var person schema.Person
	if err := c.Bind().Body(&person); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.persons.Create(ctx, &person); err != nil {
		return err
	}

	// SendSuccess 固定返回 200，Create 需要 201，因此手动构造同格式响应
	r := model.Response{
		Success:   true,
		Data:      person,
		Message:   "Person created successfully",
		Timestamp: time.Now(),
	}
	return c.Status(http.StatusCreated).JSON(r)
}

// Update 根据身份证号更新人员信息。
//
//	@Summary		更新人员信息
//	@Description	根据身份证号更新人员信息
//	@Tags			persons
//	@Accept			json
//	@Produce		json
//	@Param			idCard	path		string			true	"身份证号"
//	@Param			body	body		schema.Person	true	"更新的人员信息"
//	@Success		200		{object}	model.Response{data=schema.Person}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/persons/{idCard} [put]
func (h *PersonHandler) Update(c fiber.Ctx) error {
	ctx := c.Context()
	idCard := c.Params("idCard")

	var person schema.Person
	if err := c.Bind().Body(&person); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}
	person.IDCard = idCard

	if err := h.persons.Update(ctx, &person); err != nil {
		return err
	}

	return model.SendSuccess(c,
		model.WithData(person),
		model.WithMessage("Person updated successfully"),
	)
}

// Delete 根据身份证号删除人员记录。
//
//	@Summary		删除人员
//	@Description	根据身份证号删除人员记录
//	@Tags			persons
//	@Produce		json
//	@Param			idCard	path		string	true	"身份证号"
//	@Success		200		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/persons/{idCard} [delete]
func (h *PersonHandler) Delete(c fiber.Ctx) error {
	ctx := c.Context()
	idCard := c.Params("idCard")

	if err := h.persons.Delete(ctx, idCard); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithMessage("Person deleted successfully"))
}

// ===================== RegionHandler =====================

type RegionHandler struct {
	regions repo.RegionRepository
}

func NewRegionHandler(regions repo.RegionRepository) *RegionHandler {
	return &RegionHandler{regions: regions}
}

// List 查询全部地区列表。
//
//	@Summary		查询全部地区
//	@Description	查询全部地区列表
//	@Tags			regions
//	@Produce		json
//	@Success		200	{object}	model.Response{data=[]schema.Region}
//	@Failure		500	{object}	model.Response
//	@Router			/regions [get]
func (h *RegionHandler) List(c fiber.Ctx) error {
	ctx := c.Context()

	regions, err := h.regions.FindAll(ctx)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(regions))
}

// GetByID 根据 ID 查询地区详情。
//
//	@Summary		查询地区详情
//	@Description	根据地区 ID 查询单个地区信息
//	@Tags			regions
//	@Produce		json
//	@Param			id		path		int		true	"地区 ID"
//	@Success		200		{object}	model.Response{data=schema.Region}
//	@Failure		400		{object}	model.Response	"无效的地区 ID"
//	@Failure		404		{object}	model.Response	"地区不存在"
//	@Failure		500		{object}	model.Response
//	@Router			/regions/{id} [get]
func (h *RegionHandler) GetByID(c fiber.Ctx) error {
	ctx := c.Context()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid region ID")
	}

	region, err := h.regions.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(region))
}

// ListProvinces 查询全部省级地区。
//
//	@Summary		查询省级地区
//	@Description	查询全部省级（顶层）地区列表
//	@Tags			regions
//	@Produce		json
//	@Success		200	{object}	model.Response{data=[]schema.Region}
//	@Failure		500	{object}	model.Response
//	@Router			/regions/provinces [get]
func (h *RegionHandler) ListProvinces(c fiber.Ctx) error {
	ctx := c.Context()

	provinces, err := h.regions.FindAllProvinces(ctx)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(provinces))
}

// ListByParent 根据上级地区 ID 查询子地区。
//
//	@Summary		按上级查询子地区
//	@Description	根据 parentID 查询下一级地区列表
//	@Tags			regions
//	@Produce		json
//	@Param			parentID	query		int		true	"上级地区 ID"
//	@Success		200			{object}	model.Response{data=[]schema.Region}
//	@Failure		400			{object}	model.Response	"无效的 parentID"
//	@Failure		500			{object}	model.Response
//	@Router			/regions/by-parent [get]
func (h *RegionHandler) ListByParent(c fiber.Ctx) error {
	ctx := c.Context()
	parentID, err := strconv.Atoi(c.Query("parentID", "0"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid parent ID")
	}

	regions, err := h.regions.FindByParentID(ctx, parentID)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(regions))
}

// Create 创建地区记录。
//
//	@Summary		创建地区
//	@Description	创建新地区记录
//	@Tags			regions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.Region	true	"地区信息"
//	@Success		200		{object}	model.Response{data=schema.Region}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/regions [post]
func (h *RegionHandler) Create(c fiber.Ctx) error {
	ctx := c.Context()

	var region schema.Region
	if err := c.Bind().Body(&region); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.regions.Create(ctx, &region); err != nil {
		return err
	}

	return model.SendSuccess(c,
		model.WithData(region),
		model.WithMessage("Region created successfully"),
	)
}

// Update 根据 ID 更新地区信息。
//
//	@Summary		更新地区信息
//	@Description	根据地区 ID 更新地区信息
//	@Tags			regions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"地区 ID"
//	@Param			body	body		schema.Region	true	"更新的地区信息"
//	@Success		200		{object}	model.Response{data=schema.Region}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/regions/{id} [put]
func (h *RegionHandler) Update(c fiber.Ctx) error {
	ctx := c.Context()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid region ID")
	}

	var region schema.Region
	if err := c.Bind().Body(&region); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}
	region.ID = id

	if err := h.regions.Update(ctx, &region); err != nil {
		return err
	}

	return model.SendSuccess(c,
		model.WithData(region),
		model.WithMessage("Region updated successfully"),
	)
}

// Delete 根据 ID 删除地区记录。
//
//	@Summary		删除地区
//	@Description	根据地区 ID 删除地区记录
//	@Tags			regions
//	@Produce		json
//	@Param			id		path		int		true	"地区 ID"
//	@Success		200		{object}	model.Response
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/regions/{id} [delete]
func (h *RegionHandler) Delete(c fiber.Ctx) error {
	ctx := c.Context()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid region ID")
	}

	if err := h.regions.Delete(ctx, id); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithMessage("Region deleted successfully"))
}
