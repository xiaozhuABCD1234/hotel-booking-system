// Package service 提供业务逻辑层，位于 handler 和 repo 之间。
package service

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/utils"

	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// defaultPresignedExpiry 预签名 URL 默认有效期（永久密钥下设为 10 年）。
const defaultPresignedExpiry = 10 * 365 * 24 * time.Hour

// COSService 腾讯云对象存储（COS）业务封装。
type COSService struct {
	client    *cos.Client
	secretID  string
	secretKey string
}

// NewCOSService 根据环境变量创建 COSService。
// 未配置 COS_BUCKET_URL 时返回 nil（不阻塞启动）。
func NewCOSService() *COSService {
	bucketURL := utils.GetEnv("COS_BUCKET_URL", "")
	if bucketURL == "" {
		return nil
	}

	u, err := url.Parse(bucketURL)
	if err != nil {
		return nil
	}

	secretID := os.Getenv("SECRETID")
	secretKey := os.Getenv("SECRETKEY")

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	return &COSService{
		client:    client,
		secretID:  secretID,
		secretKey: secretKey,
	}
}

// UploadResult 上传返回值。
type UploadResult struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	FileName string `json:"filename"`
}

// Upload 上传单个文件到 COS，返回长期有效的预签名 URL。
func (s *COSService) Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error) {
	ext := filepath.Ext(file.Filename)
	key := "uploads/" + uuid.New().String() + ext

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	if _, err = s.client.Object.Put(ctx, key, src, nil); err != nil {
		return nil, err
	}

	presignedURL, err := s.client.Object.GetPresignedURL(
		ctx, http.MethodGet, key,
		s.secretID, s.secretKey,
		defaultPresignedExpiry, nil,
	)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:      presignedURL.String(),
		Key:      key,
		FileName: file.Filename,
	}, nil
}

// PresignedURL 生成对象的预签名下载 URL。
// expire <= 0 时使用默认有效期（10 年）。
func (s *COSService) PresignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if expire <= 0 {
		expire = defaultPresignedExpiry
	}

	presignedURL, err := s.client.Object.GetPresignedURL(
		ctx, http.MethodGet, key,
		s.secretID, s.secretKey,
		expire, nil,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// Delete 从 COS 删除指定 key 的对象。
func (s *COSService) Delete(ctx context.Context, key string) error {
	_, err := s.client.Object.Delete(ctx, key)
	return err
}

// KeyFromURL 从 COS 对象完整 URL（含预签名参数）中提取 object key。
// 例如 "https://bucket.cos.region.com/uploads/a.png?sign=..." → "uploads/a.png"
func KeyFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(u.Path, "/")
}
