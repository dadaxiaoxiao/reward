//go:build wireinject

package main

import (
	"github.com/dadaxiaoxiao/go-pkg/customserver"
	"github.com/dadaxiaoxiao/reward/internal/events"
	"github.com/dadaxiaoxiao/reward/internal/grpc"
	"github.com/dadaxiaoxiao/reward/internal/repository"
	"github.com/dadaxiaoxiao/reward/internal/repository/cache"
	"github.com/dadaxiaoxiao/reward/internal/repository/dao"
	"github.com/dadaxiaoxiao/reward/internal/service"
	"github.com/dadaxiaoxiao/reward/ioc"
	"github.com/google/wire"
)

var thirdPartyProvider = wire.NewSet(
	ioc.InitEtcdClient,
	ioc.InitOTEL,
	ioc.InitRedis,
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitPaymentGRPCClient,
	ioc.InitAccountGRPCClient,
	ioc.InitKafka,
	ioc.NewConsumers,
)

func InitApp() *customserver.App {
	wire.Build(
		thirdPartyProvider,
		dao.NewRewardGORMDAO,
		cache.NewRewardCache,
		repository.NewRewardRepository,
		service.NewWechatNativeRewardService,
		events.NewPaymentEventConsumer,
		grpc.NewRewardServiceServer,
		ioc.InitGRPCServer,
		wire.Struct(new(customserver.App), "GRPCServer", "Consumers"),
	)
	return new(customserver.App)
}
