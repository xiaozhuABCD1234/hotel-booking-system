// Package handler 提供 HTTP 请求处理，包括认证相关端点。
package handler

import (
	"errors"
	"net/http"
	"time"

	"backend/auth"
	"backend/model"
	schema "backend/model/schema"
	"backend/repo"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler 认证 HTTP 处理器，处理注册、登录、令牌刷新、登出。
type AuthHandler struct {
	users     repo.UserRepository
	blacklist repo.BlacklistRepository
}

// NewAuthHandler 创建 AuthHandler 实例。
func NewAuthHandler(users repo.UserRepository, blacklist repo.BlacklistRepository) *AuthHandler {
	return &AuthHandler{users: users, blacklist: blacklist}
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
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type logoutRequest struct {
	AccessToken string `json:"accessToken"`
}

// passwordMinLen 注册密码最小长度。
const passwordMinLen = 6

// ─── Register ──────────────────────────────────────────────────

// Register 用户注册，创建账户并返回 JWT 访问令牌。
//
//	@Summary		用户注册
//	@Description	创建新账户并返回 JWT 访问令牌
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		registerRequest					true	"注册信息"
//	@Success		200		{object}	model.Response{data=tokenResponse}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		409		{object}	model.Response	"用户名已存在"
//	@Failure		500		{object}	model.Response	"服务器错误"
//	@Router			/auth/register [post]
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

	return h.issueTokens(c, user.ID.String(), string(user.Role))
}

// ─── Login ─────────────────────────────────────────────────────

// Login 用户登录，验证凭据并返回 JWT 访问令牌。
//
//	@Summary		用户登录
//	@Description	验证用户名密码并返回 JWT 访问令牌
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		loginRequest					true	"登录凭据"
//	@Success		200		{object}	model.Response{data=tokenResponse}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		401		{object}	model.Response	"用户名或密码错误"
//	@Failure		500		{object}	model.Response	"服务器错误"
//	@Router			/auth/login [post]
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

	return h.issueTokens(c, user.ID.String(), string(user.Role))
}

// ─── Refresh ───────────────────────────────────────────────────

// Refresh 刷新令牌，用 refreshToken 换取新的 accessToken + refreshToken。
//
//	@Summary		刷新令牌
//	@Description	用刷新令牌换取新的访问令牌和刷新令牌
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		refreshRequest				true	"刷新令牌"
//	@Success		200		{object}	model.Response{data=tokenResponse}
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		401		{object}	model.Response	"刷新令牌无效或已过期"
//	@Router			/auth/refresh [post]
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req refreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}
	if req.RefreshToken == "" {
		return model.SendError(c, http.StatusBadRequest, "refreshToken is required")
	}

	claims, err := auth.ParseAccessToken(req.RefreshToken)
	if err != nil {
		return model.SendError(c, http.StatusUnauthorized, "Invalid or expired refresh token")
	}

	return h.issueTokens(c, claims.UserID, claims.Role)
}

// ─── Logout ────────────────────────────────────────────────────

// Logout 登出，将当前访问令牌加入黑名单，使其立即失效。
//
//	@Summary		登出
//	@Description	将当前访问令牌加入黑名单，使其立即失效
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		logoutRequest				true	"访问令牌"
//	@Success		200		{object}	model.Response
//	@Failure		400		{object}	model.Response	"请求参数无效"
//	@Failure		401		{object}	model.Response	"令牌无效或已过期"
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req logoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return model.SendError(c, http.StatusBadRequest, "Invalid request body")
	}
	if req.AccessToken == "" {
		return model.SendError(c, http.StatusBadRequest, "accessToken is required")
	}

	claims, err := auth.ParseAccessToken(req.AccessToken)
	if err != nil {
		return model.SendError(c, http.StatusUnauthorized, "Invalid or expired token")
	}

	// 校验 token 类型：accessToken 有效期 15 分钟，refreshToken 7 天
	if claims.ExpiresAt.Time.After(time.Now().Add(1 * time.Hour)) {
		return model.SendError(c, http.StatusBadRequest, "expected an access token, not a refresh token")
	}

	if err := h.blacklist.Insert(c.Context(), claims.ID, claims.ExpiresAt.Time); err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithMessage("Logged out successfully"))
}

// ─── 内部辅助 ─────────────────────────────────────────────────

// issueTokens 签发 accessToken（15 分钟）和 refreshToken（7 天），返回 200。
func (h *AuthHandler) issueTokens(c fiber.Ctx, userID, role string) error {
	accessToken, err := auth.GenerateAccessToken(userID, role)
	if err != nil {
		return err
	}
	refreshToken, err := auth.GenerateRefreshToken(userID, role)
	if err != nil {
		return err
	}

	return model.SendSuccess(c, model.WithData(tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}))
}
