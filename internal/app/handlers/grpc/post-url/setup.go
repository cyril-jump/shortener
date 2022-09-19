package posturl

import (
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/post-url/adapters"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/post-url/requests"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/post-url/usecase"
	"github.com/cyril-jump/shortener/internal/app/storage"
)

func Setup(provider storage.DB, users storage.Users, cfg storage.Cfg,
) *requests.Requests {
	repo := adapters.New(provider)

	uc := usecase.New(repo, cfg)

	reqs := requests.New(uc, users)

	return reqs
}
