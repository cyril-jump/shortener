package geturl

import (
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-url/adapters"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-url/requests"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-url/usecase"
	"github.com/cyril-jump/shortener/internal/app/storage"
)

func Setup(provider storage.DB,
) *requests.Requests {
	repo := adapters.New(provider)

	uc := usecase.New(repo)

	reqs := requests.New(uc)

	return reqs
}
