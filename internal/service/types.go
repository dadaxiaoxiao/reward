package service

import (
	"context"
	"github.com/dadaxiaoxiao/reward/internal/domain"
)

// RewardService 打赏服务
//
// 这里统一了一个抽象打赏的接口
// 不同渠道的打赏，各自实现这个接口
// ps: 因为支付渠道的参数通常不一致，所以支付微服务那里，没法统一一个渠道支付的接口
type RewardService interface {
	// PreReward 准备打赏
	// 对标到创建一个打赏的订单
	PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	// GetReward 获取打赏信息
	GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error)
	// UpdateReward 更新打赏状态
	UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error
}
