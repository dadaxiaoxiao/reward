package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type RewardGORMDAO struct {
	db *gorm.DB
}

func NewRewardGORMDAO(db *gorm.DB) RewardDao {
	return &RewardGORMDAO{
		db: db,
	}
}

func (dao *RewardGORMDAO) Insert(ctx context.Context, r Reward) (int64, error) {
	now := time.Now().UnixMilli()
	r.Utime = now
	r.Ctime = now
	err := dao.db.WithContext(ctx).Create(&r).Error
	return r.Id, err
}

func (dao *RewardGORMDAO) GetReward(ctx context.Context, rid int64) (Reward, error) {
	var r Reward
	err := dao.db.WithContext(ctx).Where("id =?", rid).First(&r).Error
	return r, err
}

func (dao *RewardGORMDAO) UpdateStatus(ctx context.Context, rid int64, status uint8) error {
	return dao.db.WithContext(ctx).Model(&Reward{}).Updates(
		map[string]any{
			"status": status,
			"utime":  time.Now().UnixMilli(),
		}).Error
}
