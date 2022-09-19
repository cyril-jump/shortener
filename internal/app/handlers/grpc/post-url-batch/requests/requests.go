package requests

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	protoapi "github.com/cyril-jump/shortener/pkg/api/proto"
)

type Requests struct {
	usecase Usecase
	users   Users
}

type Usecase interface {
	ProcessPostURLBatch(userID string, urls []dto.ModelURLBatchRequest) ([]dto.ModelURLBatchResponse, error)
}

type Users interface {
	CheckToken(tokenString string) (string, bool)
}

func New(usecase Usecase, users Users) *Requests {
	return &Requests{
		usecase: usecase,
		users:   users,
	}
}

func (r *Requests) PostURLBatch(ctx context.Context, req *protoapi.PostURLBatchRequest) (*protoapi.PostURLBatchResponse, error) {
	token := ctx.Value(config.UserID).(string)
	userID, ok := r.users.CheckToken(token)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	var urls []dto.ModelURLBatchRequest
	var url dto.ModelURLBatchRequest

	for _, requestBatchURL := range req.RequestUrls {
		url.BaseURL = requestBatchURL.Url
		url.CorID = requestBatchURL.CorrelationId
		urls = append(urls, url)
	}

	resURLs, err := r.usecase.ProcessPostURLBatch(userID, urls)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := protoapi.PostURLBatchResponse{}

	for n, val := range resURLs {
		res.ResponseUrls[n].Url = val.ShortURL
		res.ResponseUrls[n].CorrelationId = val.CorID
	}

	return &res, nil
}
