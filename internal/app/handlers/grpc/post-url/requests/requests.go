package requests

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	protoapi "github.com/cyril-jump/shortener/pkg/api/proto"
)

type Requests struct {
	usecase Usecase
	users   Users
}

type Usecase interface {
	ProcessPostURL(userID string, baseURL string) (string, error)
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

func (r *Requests) PostURL(ctx context.Context, req *protoapi.PostURLRequest) (*protoapi.PostURLResponse, error) {
	token := ctx.Value(config.UserID).(string)
	userID, ok := r.users.CheckToken(token)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}
	baseURL := req.BaseUrl
	shortURL, err := r.usecase.ProcessPostURL(userID, baseURL)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			res := protoapi.PostURLResponse{
				ShortUrl: shortURL,
			}
			return &res, status.Error(codes.AlreadyExists, `Entry already exists and was returned in response body`)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := protoapi.PostURLResponse{
		ShortUrl: shortURL,
	}

	return &res, nil
}
