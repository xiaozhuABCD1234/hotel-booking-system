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

func (h *PersonHandler) GetByIDCard(c fiber.Ctx) error {
	ctx := c.Context()
	idCard := c.Params("idCard")

	person, err := h.persons.FindByIDCard(ctx, idCard)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(person))
}

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

func (h *RegionHandler) List(c fiber.Ctx) error {
	ctx := c.Context()

	regions, err := h.regions.FindAll(ctx)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(regions))
}

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

func (h *RegionHandler) ListProvinces(c fiber.Ctx) error {
	ctx := c.Context()

	provinces, err := h.regions.FindAllProvinces(ctx)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(provinces))
}

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
