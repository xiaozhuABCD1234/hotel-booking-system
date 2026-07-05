// Package middleware 提供 Fiber v3 中间件，包括全局错误处理、日志、恢复等。
package middleware

import (
	"errors"
	"net/http"
	"os"

	"backend/model"
	"backend/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// ErrorHandler Fiber v3 全局错误处理器。
// 将各类错误统一转换为 model.Response 格式（SendError），
// 确保所有 API 响应都遵循统一的 {success, error, timestamp} 结构。
//
// 错误映射规则：
//   - gorm.ErrRecordNotFound → HTTP 404
//   - *fiber.Error → 使用其内置 StatusCode
//   - 其他 error → HTTP 500 Internal Server Error
func ErrorHandler(c fiber.Ctx, err error) error {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	// 1. GORM 记录未找到 → 404
	if errors.Is(err, gorm.ErrRecordNotFound) {
		code = http.StatusNotFound
		message = "Resource not found"
	}

	// 2. 非法状态流转 → 409 Conflict
	if errors.Is(err, service.ErrInvalidTransition) {
		code = http.StatusConflict
		message = err.Error()
	}

	// 3. Fiber 内置错误（fiber.NewError / fiber.ErrXxx）→ 使用其状态码
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	// 4. 生产环境隐藏内部错误详情
	if os.Getenv("APP_ENV") == "production" {
		message = "Internal Server Error"
	}

	return model.SendError(c, code, message)
}
