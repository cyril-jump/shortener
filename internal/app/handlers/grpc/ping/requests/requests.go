package requests

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protoapi "github.com/cyril-jump/shortener/pkg/api/proto"
)

type Requests struct {
	usecase Usecase
}

type Usecase interface {
	ProcessPing() error
}

func New(usecase Usecase) *Requests {
	return &Requests{
		usecase: usecase,
	}
}

func (r *Requests) Ping(_ context.Context, _ *protoapi.PingRequest) (*protoapi.PingResponse, error) {
	err := r.usecase.ProcessPing()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var res protoapi.PingResponse

	return &res, nil
}
