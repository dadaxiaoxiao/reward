package grpc

import (
	"context"
	rewardv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/reward/v1"
	"github.com/dadaxiaoxiao/reward/internal/domain"
	"github.com/dadaxiaoxiao/reward/internal/service"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RewardServiceServer struct {
	rewardv1.UnimplementedRewardServiceServer
	svc service.RewardService
}

func NewRewardServiceServer(svc service.RewardService) *RewardServiceServer {
	return &RewardServiceServer{
		svc: svc,
	}
}

// 注册grpc
func (s *RewardServiceServer) Registor(server *grpc.Server) {
	rewardv1.RegisterRewardServiceServer(server, s)
}

func (s *RewardServiceServer) PreReward(ctx context.Context, request *rewardv1.PreRewardRequest) (*rewardv1.PreRewardResponse, error) {
	ctx, span := otel.Tracer("github.com/dadaxiaoxiao/reward/internal/grpc").Start(ctx, "PreReward")
	defer func() {
		span.End()
	}()
	codeURL, err := s.svc.PreReward(ctx, domain.Reward{
		Uid: request.GetUid(),
		Traget: domain.Target{
			Biz:     request.GetBiz(),
			BizId:   request.GetBizId(),
			BizName: request.GetBizName(),
			Uid:     request.GetTargetUid(),
		},
		Amt: request.GetAmt(),
	})

	return &rewardv1.PreRewardResponse{
		CodeUrl: codeURL.URL,
		Rid:     codeURL.Rid,
	}, err
}
func (s *RewardServiceServer) GetReward(ctx context.Context, request *rewardv1.GetRewardRequest) (*rewardv1.GetRewardResponse, error) {
	_, span := otel.Tracer("github.com/dadaxiaoxiao/reward/internal/grpc").Start(ctx, "GetReward")
	defer func() {
		span.End()
	}()
	return nil, status.Errorf(codes.Unimplemented, "method GetReward not implemented")
}
