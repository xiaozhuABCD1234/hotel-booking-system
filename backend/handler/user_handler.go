package handler

import (
	"net/http"
	"time"

	"backend/model"
	schema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler 用户 CRUD HTTP 处理器。
type UserHandler struct {
	users     repo.UserRepository
	vipLevels repo.VipLevelRepository
}

// NewUserHandler 创建 UserHandler 实例。
func NewUserHandler(userRepo repo.UserRepository, vipLevelRepo repo.VipLevelRepository) *UserHandler {
	return &UserHandler{
		users:     userRepo,
		vipLevels: vipLevelRepo,
	}
}

// List 查询用户列表（支持分页与按角色筛选）。
//
//	@Summary		查询用户列表
//	@Description	分页查询用户列表，支持按角色筛选
//	@Tags			users
//	@Produce		json
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			pageSize	query		int		false	"每页数量"		default(10)
//	@Param			role		query		string	false	"角色筛选"
//	@Success		200			{object}	model.Response{data=[]schema.User}
//	@Failure		500			{object}	model.Response
//	@Security		BearerAuth
//	@Router			/users [get]
func (h *UserHandler) List(c fiber.Ctx) error {
	ctx := c.Context()

	var query struct {
		Page     int    `query:"page"`
		PageSize int    `query:"pageSize"`
		Role     string `query:"role"`
	}
	if err := c.Bind().Query(&query); err != nil {
		return err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize

	var rolePtr *schema.UserRole
	if query.Role != "" {
		role := schema.UserRole(query.Role)
		rolePtr = &role
	}

	users, total, err := h.users.FindAll(ctx, offset, query.PageSize, rolePtr)
	if err != nil {
		return err
	}

	return model.SendSuccess(c,
		model.WithData(users),
		model.WithPagination(total, query.Page, query.PageSize),
	)
}

// GetByID 根据 ID 查询单个用户。
//
//	@Summary		查询用户详情
//	@Description	根据 UUID 查询单个用户信息
//	@Tags			users
//	@Produce		json
//	@Param			id		path		string	true	"用户 ID (UUID)"
//	@Success		200		{object}	model.Response{data=schema.User}
//	@Failure		400		{object}	model.Response	"无效的用户 ID"
//	@Failure		404		{object}	model.Response	"用户不存在"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/users/{id} [get]
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}

	user, err := h.users.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(user))
}

// Create 创建用户，返回 201 Created。
//
//	@Summary		创建用户
//	@Description	创建新用户（管理员操作），密码自动 bcrypt 哈希
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.User	true	"用户信息"
//	@Success		201		{object}	model.Response{data=schema.User}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/users [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	ctx := c.Context()

	var user schema.User
	if err := c.Bind().Body(&user); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if user.Username == "" || user.Password == "" {
		return model.SendError(c, http.StatusBadRequest, "Username and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := h.users.Create(ctx, &user); err != nil {
		return err
	}

	r := model.Response{
		Success:   true,
		Data:      user,
		Timestamp: time.Now(),
	}
	return c.Status(http.StatusCreated).JSON(r)
}

// Update 根据 ID 更新用户信息。
//
//	@Summary		更新用户信息
//	@Description	根据 UUID 更新用户信息
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"用户 ID (UUID)"
//	@Param			body	body		schema.User	true	"更新的用户信息"
//	@Success		200		{object}	model.Response{data=schema.User}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/users/{id} [put]
//
// updateUserInput 仅包含客户端可更新的用户字段。
type updateUserInput struct {
	Username    string           `json:"username,omitempty"`
	OldPassword string           `json:"oldPassword,omitempty"`
	Password    string           `json:"password,omitempty"`
	Phone       *string          `json:"phone,omitempty"`
	Email       *string          `json:"email,omitempty"`
	RealName    *string          `json:"realName,omitempty"`
	IDCard      *string          `json:"idCard,omitempty"`
	Role        *schema.UserRole `json:"role,omitempty"`
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid user ID")
	}

	// 查询已有记录
	existing, err := h.users.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 绑定请求体中允许更新的字段
	var input updateUserInput
	if err := c.Bind().Body(&input); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	// 仅合并请求中显式提供的字段
	if input.Username != "" {
		existing.Username = input.Username
	}
	if input.Password != "" {
		if input.OldPassword == "" {
			return model.SendError(c, http.StatusBadRequest, "Old password is required to change password")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(input.OldPassword)); err != nil {
			return model.SendError(c, http.StatusBadRequest, "Old password is incorrect")
		}
		if len(input.Password) < passwordMinLen {
			return model.SendError(c, http.StatusBadRequest, "Password must be at least 6 characters")
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		existing.Password = string(hashed)
	}
	if input.Phone != nil {
		existing.Phone = input.Phone
	}
	if input.Email != nil {
		existing.Email = input.Email
	}
	if input.RealName != nil {
		existing.RealName = input.RealName
	}
	if input.IDCard != nil {
		existing.IDCard = input.IDCard
	}
	if input.Role != nil {
		existing.Role = *input.Role
	}

	if err := h.users.Update(ctx, existing); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(existing))
}

// Delete 根据 ID 软删除用户。
//
//	@Summary		删除用户
//	@Description	根据 UUID 软删除用户
//	@Tags			users
//	@Produce		json
//	@Param			id		path		string	true	"用户 ID (UUID)"
//	@Success		200		{object}	model.Response
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Security		BearerAuth
//	@Router			/users/{id} [delete]
func (h *UserHandler) Delete(c fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}

	if err := h.users.Delete(ctx, id); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithMessage("User deleted"))
}
