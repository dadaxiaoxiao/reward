package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/go-pkg/saramax"
	"github.com/dadaxiaoxiao/reward/internal/domain"
	"github.com/dadaxiaoxiao/reward/internal/service"
	"strings"
	"time"
)

// PaymentEvent 支付事件
type PaymentEvent struct {
	BizTradeNo string
	Status     uint8
}

func (p PaymentEvent) toDomainStatus() domain.RewardStatus {
	//PaymentStatusUnknown 0
	//PaymentStatusInit
	//PaymentStatusSuccess
	//PaymentStatusFailed
	//PaymentStatusRefund
	switch p.Status {
	case 1:
		return domain.RewardStatusInit
	case 2:
		return domain.RewardStatusPayed
	case 3, 4:
		return domain.RewardStatusFailed
	default:
		return domain.RewardStatusUnknown
	}
}

// PaymentEventConsumer 支付消费者
//
// 这里reward作为业务方消费 支付生产者消息
type PaymentEventConsumer struct {
	client sarama.Client
	l      accesslog.Logger
	svc    service.RewardService
}

func NewPaymentEventConsumer(client sarama.Client, l accesslog.Logger, svc service.RewardService) *PaymentEventConsumer {
	return &PaymentEventConsumer{
		client: client,
		l:      l,
		svc:    svc,
	}
}

func (p *PaymentEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("reward", p.client)
	if err != nil {
		return err
	}
	go func() {
		// 一条消息可以被多个不同的ConsumerGroup消费；但是一个ConsumerGroup中只能有一个Consumer能够消费该消息。
		// 不同ConsumerGroup中的Consumer可以各自独立地消费同一份消息
		// 属于同一个ConsumerGroup内的多个Consumer会竞争消费该Partition中的消息
		err := cg.Consume(context.Background(), []string{"payment_events"},
			saramax.NewHandler[PaymentEvent](p.l, p.Consume))
		if err != nil {
			p.l.Error("退出了消费循环异常", accesslog.Error(err))
		}
	}()
	return err
}

// Consume 消费逻辑
func (p *PaymentEventConsumer) Consume(msg *sarama.ConsumerMessage, evt PaymentEvent) error {
	// reward 服务，商户订单号生成逻辑 "reward-rid"
	if !strings.HasPrefix(evt.BizTradeNo, "reward") {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return p.svc.UpdateReward(ctx, evt.BizTradeNo, evt.toDomainStatus())
}
