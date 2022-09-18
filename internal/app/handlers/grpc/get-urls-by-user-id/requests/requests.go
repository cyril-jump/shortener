package requests

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/dto"
	protoapi "github.com/cyril-jump/shortener/pkg/api/proto"
)

type Requests struct {
	usecase Usecase
	users   Users
}

type Usecase interface {
	ProcessGetAllURLsByUserID(userID string) ([]dto.ModelURL, error)
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

func (r *Requests) GetURLsByUserID(ctx context.Context, _ *protoapi.GetURLsByUserIDRequest) (*protoapi.GetURLsByUserIDResponse, error) {
	token := ctx.Value(config.UserID).(string)
	userID, ok := r.users.CheckToken(token)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	urls, err := r.usecase.ProcessGetAllURLsByUserID(userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	res := protoapi.GetURLsByUserIDResponse{}
	for _, val := range urls {
		responseURL := protoapi.ResponseURLs{
			BaseUrl:  val.BaseURL,
			ShortUrl: val.ShortURL,
		}
		res.ResponseUrls = append(res.ResponseUrls, &responseURL)
	}
	return &res, nil
}
