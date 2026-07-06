package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// BlacklistRepository 黑名单存储接口，支持测试 mock。
type BlacklistRepository interface {
	Insert(ctx context.Context, jti string, expiresAt time.Time) error
	Exists(ctx context.Context, jti string) (bool, error)
}

type BlacklistRepo struct {
	db *gorm.DB
}

func NewBlacklistRepo(db *gorm.DB) *BlacklistRepo {
	return &BlacklistRepo{db: db}
}

func (r *BlacklistRepo) Insert(ctx context.Context, jti string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO jwt_blacklist_1718 (jti, expires_at) VALUES (?, ?) ON CONFLICT (jti) DO NOTHING",
		jti, expiresAt,
	).Error
}

func (r *BlacklistRepo) Exists(ctx context.Context, jti string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM jwt_blacklist_1718 WHERE jti = ?", jti,
	).Scan(&count).Error
	return count > 0, err
}
