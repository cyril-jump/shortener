package geturlsbyuserid

import (
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-urls-by-user-id/adapters"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-urls-by-user-id/requests"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-urls-by-user-id/usecase"
	"github.com/cyril-jump/shortener/internal/app/storage"
)

func Setup(provider storage.DB, users storage.Users, cfg storage.Cfg,
) *requests.Requests {
	repo := adapters.New(provider)

	uc := usecase.New(repo, cfg)

	reqs := requests.New(uc, users)

	return reqs
}
