package requests

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	protoapi "github.com/cyril-jump/shortener/pkg/api/proto"
)

type Requests struct {
	usecase Usecase
}

type Usecase interface {
	ProcessGetBaseURL(shortURL string) (string, error)
}

func New(usecase Usecase) *Requests {
	return &Requests{
		usecase: usecase,
	}
}

func (r *Requests) GetURL(_ context.Context, req *protoapi.GetURLRequest) (*protoapi.GetURLResponse, error) {

	shortURL := req.ShortUrlId

	url, err := r.usecase.ProcessGetBaseURL(shortURL)
	if err != nil {
		if errors.Is(err, errs.ErrWasDeleted) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := protoapi.GetURLResponse{
		RedirectTo: url,
	}

	return &res, nil
}
