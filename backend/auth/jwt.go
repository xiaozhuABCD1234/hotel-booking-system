// Package auth 提供 JWT 令牌生成与校验。
package auth

import (
	"errors"
	"os"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 是 JWT 令牌的自定义声明。
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AccessTokenExpiry 访问令牌过期时间。
const AccessTokenExpiry = 15 * time.Minute

// minSecretLen 是 JWT_SECRET 的最小字节数，防止过弱的签名密钥。
const minSecretLen = 16

var (
	// ErrSecretNotSet JWT 密钥未配置。
	ErrSecretNotSet = errors.New("JWT_SECRET environment variable is not set")

	// ErrSecretTooShort JWT 密钥过短，至少需要 16 字节。
	ErrSecretTooShort = errors.New("JWT_SECRET must be at least 16 bytes")

	// ErrInvalidClaims 令牌声明格式无效。
	ErrInvalidClaims = errors.New("invalid JWT claims")
)

// secretKey 以原子指针存储签名密钥，保证并发读写的线程安全。
// 测试中可多次调用 LoadSecret 切换密钥，运行期亦安全。
var secretKey atomic.Pointer[[]byte]

// LoadSecret 从环境变量 JWT_SECRET 加载签名密钥，未设置或过短时返回错误。
// 线程安全，可在运行期多次调用以热更新密钥。
func LoadSecret() error {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return ErrSecretNotSet
	}
	if len(key) < minSecretLen {
		return ErrSecretTooShort
	}
	b := []byte(key)
	secretKey.Store(&b)
	return nil
}

// GenerateAccessToken 生成 15 分钟有效的访问令牌。
func GenerateAccessToken(userID, role string) (string, error) {
	sk := secretKey.Load()
	if sk == nil {
		return "", ErrSecretNotSet
	}
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hotel-booking-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(*sk)
}

// ParseAccessToken 解析并校验访问令牌，返回声明信息。
func ParseAccessToken(tokenString string) (*Claims, error) {
	sk := secretKey.Load()
	if sk == nil {
		return nil, ErrSecretNotSet
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return *sk, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidClaims
}
