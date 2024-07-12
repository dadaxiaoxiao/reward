package domain

// 目标
type Target struct {
	// 打赏业务
	Biz   string
	BizId int64
	// 可选 ，打赏内容
	BizName string

	// 打赏的目标用户
	Uid int64
}

// Reward 打赏
type Reward struct {
	Id  int64
	Uid int64
	// 打赏目标
	Traget Target
	Amt    int64
	Status RewardStatus
}

// Completed 支付完成
//
// 支付成功，支付失败，都是处理了支付回调
func (r Reward) Completed() bool {
	return r.Status == RewardStatusFailed || r.Status == RewardStatusPayed
}

type RewardStatus uint8

func (r RewardStatus) AsUint8() uint8 {
	return uint8(r)
}

const (
	RewardStatusUnknown = iota
	RewardStatusInit
	RewardStatusPayed
	RewardStatusFailed
)

// CodeURL 二维码url
type CodeURL struct {
	// 打赏id
	Rid int64
	// 二维码url
	URL string
}
