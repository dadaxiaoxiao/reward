package service

import (
	"context"
	"github.com/dadaxiaoxiao/reward/internal/domain"
)

type WechatH5RewardService struct {
}

func (w WechatH5RewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	//TODO implement me
	panic("implement me")
}

func (w WechatH5RewardService) GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error) {
	//TODO implement me
	panic("implement me")
}

func (w WechatH5RewardService) UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {
	//TODO implement me
	panic("implement me")
}
