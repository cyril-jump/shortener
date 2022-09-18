package requests

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cyril-jump/shortener/internal/app/dto"
	protoapi "github.com/cyril-jump/shortener/pkg/api/proto"
)

type Requests struct {
	usecase Usecase
}

type Usecase interface {
	ProcessGetStats() (dto.Stat, error)
}

func New(usecase Usecase) *Requests {
	return &Requests{
		usecase: usecase,
	}
}

func (r *Requests) GetStats(_ context.Context, _ *protoapi.GetStatsRequest) (*protoapi.GetStatsResponse, error) {
	stats, err := r.usecase.ProcessGetStats()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := protoapi.GetStatsResponse{
		Users: int64(stats.Users),
		Urls:  int64(stats.URLs),
	}

	return &res, nil
}
