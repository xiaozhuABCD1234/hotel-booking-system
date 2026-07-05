// Package middleware 提供 Fiber v3 中间件，包括全局错误处理、JWT 认证等。
package middleware

import (
	"net/http"
	"strings"

	"backend/auth"

	"github.com/gofiber/fiber/v3"
)

// contextKey 是 Fiber context 中存储认证信息的键。
const contextKey = "auth_claims"

// JWTAuth 返回 JWT 认证中间件。
// 从 Authorization: Bearer <token> 头提取并校验访问令牌，
// 校验通过后将 *auth.Claims 存入 fiber.Ctx 的上下文中。
func JWTAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 提取 Authorization 头
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(http.StatusUnauthorized, "missing authorization header")
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return fiber.NewError(http.StatusUnauthorized, "invalid authorization format, expected: Bearer <token>")
		}

		tokenString := parts[1]

		// 校验 token
		claims, err := auth.ParseAccessToken(tokenString)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, "invalid or expired token")
		}

		// 将 claims 存入 context，供后续 handler 使用
		c.Locals(contextKey, claims)

		return c.Next()
	}
}

// GetClaims 从 Fiber context 中提取认证声明。
// 仅在 JWTAuth 中间件之后调用有效，否则返回 nil。
func GetClaims(c fiber.Ctx) *auth.Claims {
	if claims, ok := c.Locals(contextKey).(*auth.Claims); ok {
		return claims
	}
	return nil
}
