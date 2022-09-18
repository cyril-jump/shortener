package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils"
)

type Interceptor struct {
	users storage.Users
}

func New(users storage.Users) *Interceptor {
	return &Interceptor{
		users: users,
	}
}

func (IS *Interceptor) UserIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var userID string
	var token string
	var err error
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(config.UserID.String())
		if len(values) > 0 {
			userID, ok = IS.users.CheckToken(values[0])
			if ok {
				return handler(context.WithValue(ctx, config.UserID, userID), req)
			}
		}
	}
	uid := utils.CreateID(16)
	token, err = IS.users.CreateToken(uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, `create token error: `+err.Error())
	}
	md := metadata.New(map[string]string{config.UserID.String(): token})
	err = grpc.SetTrailer(ctx, md)
	if err != nil {
		return nil, status.Errorf(codes.Internal, `set trailer err: `+err.Error())
	}
	return handler(context.WithValue(ctx, config.UserID, userID), req)
}
