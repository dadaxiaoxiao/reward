package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dadaxiaoxiao/reward/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type rewardCache struct {
	client redis.Cmdable
}

func NewRewardCache(client redis.Cmdable) RewardCache {
	return &rewardCache{
		client: client,
	}
}

func (c *rewardCache) GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	key := c.codeURLKey(r)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.CodeURL{}, err
	}
	var res domain.CodeURL
	err = json.Unmarshal(data, &res)
	return res, err
}

func (c *rewardCache) CachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error {
	key := c.codeURLKey(r)
	data, err := json.Marshal(cu)
	if err != nil {
		return err
	}
	// 如果担心 30 分钟刚好是微信订单过期的问题，那么可以设置成 29 分钟
	return c.client.Set(ctx, key, data, time.Minute*29).Err()
}

func (c *rewardCache) DeleteCachedCodeURL(ctx context.Context, r domain.Reward) error {
	key := c.codeURLKey(r)
	return c.client.Del(ctx, key).Err()
}

func (c *rewardCache) codeURLKey(r domain.Reward) string {
	return fmt.Sprintf("reward:code_url:%s:%d:%d", r.Traget.Biz, r.Traget.BizId, r.Uid)
}
