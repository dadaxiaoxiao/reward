package repository

import (
	"context"
	"github.com/dadaxiaoxiao/reward/internal/domain"
	"github.com/dadaxiaoxiao/reward/internal/repository/cache"
	"github.com/dadaxiaoxiao/reward/internal/repository/dao"
)

type rewardRepository struct {
	dao   dao.RewardDao
	cache cache.RewardCache
}

func NewRewardRepository(dao dao.RewardDao, cache cache.RewardCache) RewardRepository {
	return &rewardRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *rewardRepository) CreateReward(ctx context.Context, reward domain.Reward) (int64, error) {
	return r.dao.Insert(ctx, r.toEntity(reward))

}

func (r *rewardRepository) GetReward(ctx context.Context, rid int64) (domain.Reward, error) {
	data, err := r.dao.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, err
	}
	return r.toDomain(data), nil
}

func (r *rewardRepository) UpdateStatus(ctx context.Context, rid int64, status domain.RewardStatus) error {
	return r.dao.UpdateStatus(ctx, rid, status.AsUint8())
}

func (r *rewardRepository) toDomain(reward dao.Reward) domain.Reward {
	return domain.Reward{
		Id:  reward.Id,
		Uid: reward.Uid,
		Traget: domain.Target{
			Biz:     reward.Biz,
			BizId:   reward.BizId,
			BizName: reward.BizName,
			Uid:     reward.TargetUid,
		},
		Amt:    reward.Amount,
		Status: domain.RewardStatus(reward.Status),
	}
}

func (r *rewardRepository) toEntity(reward domain.Reward) dao.Reward {
	return dao.Reward{
		Biz:       reward.Traget.Biz,
		BizId:     reward.Traget.BizId,
		BizName:   reward.Traget.BizName,
		TargetUid: reward.Traget.Uid,
		Status:    reward.Status.AsUint8(),
		Uid:       reward.Uid,
		Amount:    reward.Amt,
	}
}

func (r *rewardRepository) CachedCodeURL(ctx context.Context, cu domain.CodeURL, reward domain.Reward) error {
	return r.cache.CachedCodeURL(ctx, cu, reward)
}

func (r *rewardRepository) GetCachedCodeURL(ctx context.Context, reward domain.Reward) (domain.CodeURL, error) {
	return r.cache.GetCachedCodeURL(ctx, reward)
}

func (r *rewardRepository) DeleteCachedCodeURL(ctx context.Context, reward domain.Reward) error {
	return r.cache.DeleteCachedCodeURL(ctx, reward)
}
