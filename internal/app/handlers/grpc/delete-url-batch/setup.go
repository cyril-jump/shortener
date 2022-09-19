package deleteurlbatch

import (
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/delete-url-batch/requests"
	"github.com/cyril-jump/shortener/internal/app/storage"
)

func Setup(users storage.Users, inWorker storage.InWorker,
) *requests.Requests {

	reqs := requests.New(inWorker, users)

	return reqs
}
