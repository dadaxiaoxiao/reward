package cache

import (
	"context"
	"github.com/dadaxiaoxiao/reward/internal/domain"
)

type RewardCache interface {
	GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	CachedCodeURL(ctx context.Context, cu domain.CodeURL, reward domain.Reward) error
	DeleteCachedCodeURL(ctx context.Context, r domain.Reward) error
}
