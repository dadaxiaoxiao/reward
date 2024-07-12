package service

import (
	"context"
	"errors"
	"fmt"
	accountv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/account/v1"
	pmtv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/payment/v1"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/reward/internal/domain"
	"github.com/dadaxiaoxiao/reward/internal/repository"
	"go.opentelemetry.io/otel"
	"strconv"
	"strings"
)

// WechatNativeRewardService 微信 Native打赏
type WechatNativeRewardService struct {
	paymentCli            pmtv1.WechatPaymentServiceClient
	accountCli            accountv1.AccountServiceClient
	repo                  repository.RewardRepository
	l                     accesslog.Logger
	paymentStatusToStatus map[pmtv1.PaymentStatus]domain.RewardStatus
}

func NewWechatNativeRewardService(paymentCli pmtv1.WechatPaymentServiceClient,
	accountCli accountv1.AccountServiceClient,
	repo repository.RewardRepository,
	l accesslog.Logger) RewardService {
	return &WechatNativeRewardService{
		paymentCli: paymentCli,
		accountCli: accountCli,
		repo:       repo,
		l:          l,
		paymentStatusToStatus: map[pmtv1.PaymentStatus]domain.RewardStatus{
			pmtv1.PaymentStatus_PaymentStatusInit:    domain.RewardStatusInit,
			pmtv1.PaymentStatus_PaymentStatusUnknown: domain.RewardStatusUnknown,
			pmtv1.PaymentStatus_PaymentStatusSuccess: domain.RewardStatusPayed,
			pmtv1.PaymentStatus_PaymentStatusFailed:  domain.RewardStatusFailed,
			pmtv1.PaymentStatus_PaymentStatusRefund:  domain.RewardStatusFailed,
		},
	}
}

// PreReward 预打赏
func (w *WechatNativeRewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	ctx, span := otel.Tracer("github.com/dadaxiaoxiao/reward/internal/service").
		Start(ctx, "PreReward")
	defer func() {
		span.End()
	}()
	// 缓存二维码，一旦发现支付成功，就清除打赏人的二维码
	cu, err := w.repo.GetCachedCodeURL(ctx, r)
	if err == nil {
		return cu, nil
	}

	r.Status = domain.RewardStatusInit
	// 先创建打赏记录
	rid, err := w.repo.CreateReward(ctx, r)
	if err != nil {
		return domain.CodeURL{}, nil
	}

	//创建预支付
	resp, err := w.paymentCli.NativePrePay(ctx, &pmtv1.PrePayRequest{
		Amt: &pmtv1.Amount{
			Total: r.Amt,
			// 这里统一默认人民币
			Currency: "CNY",
		},
		BizTradeNo:  w.bizTradeNO(rid),
		Description: fmt.Sprintf("打赏-%s", r.Traget.BizName),
	})
	if err != nil {
		return domain.CodeURL{}, nil
	}
	// 组合返回
	res := domain.CodeURL{
		Rid: rid,
		URL: resp.GetCodeUrl(),
	}

	go func() {
		err1 := w.repo.CachedCodeURL(ctx, cu, r)
		if err1 != nil {
			w.l.Error("缓存打赏二维码失败", accesslog.Int64("Uid", r.Uid),
				accesslog.String("Biz", r.Traget.Biz),
				accesslog.Int64("BizId", r.Traget.BizId),
				accesslog.String("BizName", r.Traget.BizName),
				accesslog.Error(err1))
		}
	}()
	return res, nil
}

// GetReward 获取打赏记录
func (w *WechatNativeRewardService) GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error) {
	ctx, span := otel.Tracer("github.com/dadaxiaoxiao/reward/internal/service").Start(ctx, "GetReward")
	defer func() {
		span.End()
	}()
	// 快路径查询
	r, err := w.repo.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, nil
	}
	// 查询人是否打赏人
	if rid != r.Uid {
		// 说明是非法查询
		return domain.Reward{}, errors.New("查询的打赏记录和打赏人对不上")
	}
	if r.Completed() || ctx.Value("limited") == "true" {
		// 支付完成或者处于降级
		return r, nil
	}
	// 打赏没有完成的原因 1.用户没有完成支付  2. 用户支付了，但是业务方还没有收到通知
	// 这里可以执行慢路径 直接查询 payment 支付微服务
	// 只能解决，支付收到微信收到，但是 消息队列没通知到 reward 服务
	// 降级状态，限流状态，熔断状态，不要走慢路径
	resp, err := w.paymentCli.GetPayment(ctx, &pmtv1.GetPaymentRequest{
		BizTradeNo: w.bizTradeNO(rid),
	})
	if err != nil {
		// 这里直接返回从数据库查询的结果，不是从微信服务查询的结果
		// 如果支付服务查询支付，也有类似的慢路径，那么最后的结果就是 微信服务的查询结果
		w.l.Error("慢路径查询支付结果失败",
			accesslog.Int64("rid", r.Id), accesslog.Error(err))
		return r, nil
	}

	// 更新本地状态
	r.Status = w.paymentStatusToStatus[resp.GetStatus()]
	err = w.repo.UpdateStatus(ctx, rid, r.Status)
	if err != nil {
		w.l.Error("更新本地打赏状态失败",
			accesslog.Int64("rid", r.Id), accesslog.Error(err))
		return r, nil
	}
	go func() {
		err1 := w.repo.DeleteCachedCodeURL(ctx, r)
		w.l.Error("删除本地打赏缓存codeurl 失败",
			accesslog.Int64("rid", r.Id), accesslog.Error(err1))
	}()
	return r, nil
}

// UpdateReward 更新打赏
//
// 暴露给消息队列使用
func (w *WechatNativeRewardService) UpdateReward(ctx context.Context,
	bizTradeNO string, status domain.RewardStatus) error {
	rid := w.toRid(bizTradeNO)
	err := w.repo.UpdateStatus(ctx, rid, status)

	// 支付完成，请求账号服务，入账
	if status == domain.RewardStatusPayed {
		r, err := w.repo.GetReward(ctx, rid)
		if err != nil {
			return err
		}

		weAmt := int64(float64(r.Amt) * 0.1)
		_, err = w.accountCli.Credit(ctx, &accountv1.CreditRequest{
			Biz:   "reward",
			BizId: rid,
			Items: []*accountv1.CreditItem{
				{
					AccountType: accountv1.AccountType_AccountTypeSystem,
					// 虽然可能为 0，但是也要记录出来
					Amt:      weAmt,
					Currency: "CNY",
				},
				{
					Account:     r.Uid,
					Uid:         r.Uid,
					AccountType: accountv1.AccountType_AccountTypeReward,
					Amt:         r.Amt - weAmt,
					Currency:    "CNY",
				},
			},
		})
	}
	return err
}

// bizTradeNO 生成商户订单号
//
// 业务方生成商品订单号
func (w *WechatNativeRewardService) bizTradeNO(rid int64) string {
	return fmt.Sprintf("reward-%d", rid)
}

// toRid 商户订单号 解析 打赏id
func (s *WechatNativeRewardService) toRid(tradeNO string) int64 {
	ridStr := strings.Split(tradeNO, "-")
	val, _ := strconv.ParseInt(ridStr[1], 10, 64)
	return val
}
