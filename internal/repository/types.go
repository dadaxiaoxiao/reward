package repository

import (
	"context"
	"github.com/dadaxiaoxiao/reward/internal/domain"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward domain.Reward) (int64, error)
	GetReward(ctx context.Context, rid int64) (domain.Reward, error)
	UpdateStatus(ctx context.Context, rid int64, status domain.RewardStatus) error
	// CachedCodeURL 缓存 支付链接
	CachedCodeURL(ctx context.Context, cu domain.CodeURL, reward domain.Reward) error
	GetCachedCodeURL(ctx context.Context, reward domain.Reward) (domain.CodeURL, error)
	DeleteCachedCodeURL(ctx context.Context, reward domain.Reward) error
}
