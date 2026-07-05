package model

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Response 统一响应格式
type Response struct {
	Success    bool           `json:"success"`
	Data       any            `json:"data,omitempty"`
	Message    string         `json:"message,omitempty"`
	Error      *ErrorObj      `json:"error,omitempty"`
	Pagination *PaginationObj `json:"pagination,omitempty"`
	Timestamp  time.Time      `json:"timestamp"`
}

type ErrorObj struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type PaginationObj struct {
	CurrentPage  int   `json:"currentPage"`
	TotalPages   int   `json:"totalPages"`
	TotalItems   int64 `json:"totalItems"`
	ItemsPerPage int   `json:"itemsPerPage"`
	HasNext      bool  `json:"hasNext"`
	HasPrev      bool  `json:"hasPrev"`
}

// Option 响应构建选项
type Option func(*Response)

// WithData 设置响应数据
func WithData(data any) Option {
	return func(r *Response) { r.Data = data }
}

// WithMessage 设置响应消息
func WithMessage(msg string) Option {
	return func(r *Response) { r.Message = msg }
}

// WithPagination 设置分页信息
func WithPagination(total int64, page, perPage int) Option {
	return func(r *Response) {
		r.Pagination = NewPagination(total, page, perPage)
	}
}

// WithErrorCode 设置自定义错误码（覆盖默认的 HTTP StatusText）
func WithErrorCode(code string) Option {
	return func(r *Response) {
		if r.Error != nil {
			r.Error.Code = code
		}
	}
}

// WithErrorDetails 设置错误详情
func WithErrorDetails(details any) Option {
	return func(r *Response) {
		if r.Error != nil {
			r.Error.Details = details
		}
	}
}

// SendSuccess 成功响应
func SendSuccess(c fiber.Ctx, opts ...Option) error {
	r := Response{
		Success:   true,
		Timestamp: time.Now(),
	}
	for _, opt := range opts {
		opt(&r)
	}
	return c.Status(http.StatusOK).JSON(r)
}

// SendError 错误响应
func SendError(c fiber.Ctx, code int, message string, opts ...Option) error {
	r := Response{
		Success:   false,
		Message:   message,
		Error:     &ErrorObj{Code: http.StatusText(code), Message: message},
		Timestamp: time.Now(),
	}
	for _, opt := range opts {
		opt(&r)
	}
	return c.Status(code).JSON(r)
}

// NewPagination 计算分页信息
func NewPagination(total int64, page, perPage int) *PaginationObj {
	if perPage <= 0 {
		perPage = 10
	}
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	return &PaginationObj{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   total,
		ItemsPerPage: perPage,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}
