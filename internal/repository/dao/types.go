package dao

import "context"

type RewardDao interface {
	Insert(ctx context.Context, r Reward) (int64, error)
	GetReward(ctx context.Context, rid int64) (Reward, error)
	UpdateStatus(ctx context.Context, rid int64, status uint8) error
}

type Reward struct {
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Biz     string `gorm:"column:biz;index:biz_biz_id"`
	BizId   int64  `gorm:"column:biz_id;index:biz_biz_id"`
	BizName string `gorm:"column:biz_name;type:varchar(128)"`
	// 被打赏的人
	TargetUid int64 `gorm:"column:target_uid;index"`
	// 使用 RewardStatus
	Status uint8 `gorm:"column:status"`
	// 打赏的人
	Uid    int64 `gorm:"column:uid"`
	Amount int64 `gorm:"column:amount"`
	Utime  int64
	Ctime  int64
}
