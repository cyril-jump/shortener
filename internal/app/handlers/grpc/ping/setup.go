package ping

import (
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/ping/adapters"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/ping/requests"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/ping/usecase"
	"github.com/cyril-jump/shortener/internal/app/storage"
)

func Setup(provider storage.DB,
) *requests.Requests {
	repo := adapters.New(provider)

	uc := usecase.New(repo)

	reqs := requests.New(uc)

	return reqs
}
