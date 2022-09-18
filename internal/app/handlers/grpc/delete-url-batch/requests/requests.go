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
	inWorker InWorker
	users    Users
}

type InWorker interface {
	Do(t dto.Task)
}

type Users interface {
	CheckToken(tokenString string) (string, bool)
}

func New(inWorker InWorker, users Users) *Requests {
	return &Requests{
		inWorker: inWorker,
		users:    users,
	}
}

func (r *Requests) DeleteURLBatch(ctx context.Context, req *protoapi.DeleteURLBatchRequest) (*protoapi.DeleteURLBatchResponse, error) {
	token := ctx.Value(config.UserID).(string)
	userID, ok := r.users.CheckToken(token)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	var task dto.Task
	task.ID = userID

	for _, delURL := range req.RequestUrls.Urls {
		task.ShortURL = delURL
		r.inWorker.Do(task)
	}

	res := protoapi.DeleteURLBatchResponse{}

	return &res, nil
}
