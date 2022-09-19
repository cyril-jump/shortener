package posturlbatch

import (
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/post-url-batch/adapters"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/post-url-batch/requests"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/post-url-batch/usecase"
	"github.com/cyril-jump/shortener/internal/app/storage"
)

func Setup(provider storage.DB, users storage.Users, cfg storage.Cfg,
) *requests.Requests {
	repo := adapters.New(provider)

	uc := usecase.New(repo, cfg)

	reqs := requests.New(uc, users)

	return reqs
}
