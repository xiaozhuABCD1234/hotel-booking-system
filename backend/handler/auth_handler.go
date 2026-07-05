// Package handler 提供 HTTP 请求处理，包括认证相关端点。
package handler

import (
	"errors"
	"net/http"

	"backend/auth"
	"backend/model"
	schema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler 认证 HTTP 处理器，处理注册、登录、令牌刷新和登出。
type AuthHandler struct {
	users repo.UserRepository
}

// NewAuthHandler 创建 AuthHandler 实例。
func NewAuthHandler(users repo.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

// ─── 请求/响应结构 ────────────────────────────────────────────

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	RealName string `json:"realName,omitempty"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"` // 秒
	TokenType   string `json:"tokenType"`
}

// passwordMinLen 注册密码最小长度。
const passwordMinLen = 6

// ─── Register ──────────────────────────────────────────────────

// Register 用户注册，创建账户并返回 JWT 访问令牌。
//
//	POST /api/v1/auth/register
func (h *AuthHandler) Register(c fiber.Ctx) error {
	ctx := c.Context()

	var req registerRequest
	if err := c.Bind().Body(&req); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if req.Username == "" || req.Password == "" {
		return model.SendError(c, http.StatusBadRequest, "Username and password are required")
	}
	if len(req.Password) < passwordMinLen {
		return model.SendError(c, http.StatusBadRequest, "Password must be at least 6 characters")
	}

	// 检查用户名是否已存在
	existing, err := h.users.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return model.SendError(c, http.StatusConflict, "Username already exists")
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := schema.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     schema.RoleCustomer,
		Points:   0,
		Status:   1,
	}

	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.RealName != "" {
		user.RealName = &req.RealName
	}

	if err := h.users.Create(ctx, &user); err != nil {
		return err
	}

	return h.issueAccessToken(c, user.ID.String(), string(user.Role))
}

// ─── Login ─────────────────────────────────────────────────────

// Login 用户登录，验证凭据并返回 JWT 访问令牌。
//
//	POST /api/v1/auth/login
func (h *AuthHandler) Login(c fiber.Ctx) error {
	ctx := c.Context()

	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	if req.Username == "" || req.Password == "" {
		return model.SendError(c, http.StatusBadRequest, "Username and password are required")
	}

	// 查找用户
	user, err := h.users.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.SendError(c, http.StatusUnauthorized, "Invalid username or password")
		}
		return err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return model.SendError(c, http.StatusUnauthorized, "Invalid username or password")
	}

	return h.issueAccessToken(c, user.ID.String(), string(user.Role))
}

// ─── Refresh ───────────────────────────────────────────────────

// Refresh 刷新令牌（暂未实现，保留接口待后续开发）。
//
//	POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	return model.SendError(c, http.StatusNotImplemented, "Refresh token not yet implemented")
}

// ─── Logout ────────────────────────────────────────────────────

// Logout 登出（无状态 JWT，服务端不存储令牌，客户端自行丢弃令牌即可）。
//
//	POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	return model.SendSuccess(c, model.WithMessage("Logged out — discard your access token on the client"))
}

// ─── 内部辅助 ─────────────────────────────────────────────────

// issueAccessToken 签发 JWT 访问令牌（15 分钟有效），返回 200。
// Register 与 Login 都返回令牌而非用户资源，故统一用 200 而非 201。
func (h *AuthHandler) issueAccessToken(c fiber.Ctx, userID, role string) error {
	accessToken, err := auth.GenerateAccessToken(userID, role)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(tokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(auth.AccessTokenExpiry.Seconds()),
		TokenType:   "Bearer",
	}))
}
