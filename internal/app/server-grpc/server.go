package servergrpc

import (
	"google.golang.org/grpc"

	"github.com/cyril-jump/shortener/internal/app/Interceptors"
	deleteurlbatch "github.com/cyril-jump/shortener/internal/app/handlers/grpc/delete-url-batch"
	getstats "github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-stats"
	geturl "github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-url"
	geturlsbyuserid "github.com/cyril-jump/shortener/internal/app/handlers/grpc/get-urls-by-user-id"
	"github.com/cyril-jump/shortener/internal/app/handlers/grpc/ping"
	posturl "github.com/cyril-jump/shortener/internal/app/handlers/grpc/post-url"
	posturlbatch "github.com/cyril-jump/shortener/internal/app/handlers/grpc/post-url-batch"
	"github.com/cyril-jump/shortener/internal/app/storage"
	protoapi "github.com/cyril-jump/shortener/pkg/api/proto"
)

func InitSrv(db storage.DB, cfg storage.Cfg, usr storage.Users, inWorker storage.InWorker) *grpc.Server {

	// handlers
	pingReqs := ping.Setup(db)
	getStatsReqs := getstats.Setup(db)
	getURL := geturl.Setup(db)
	postURL := posturl.Setup(db, usr, cfg)
	getURLsByUserID := geturlsbyuserid.Setup(db, usr, cfg)
	postURLBatch := posturlbatch.Setup(db, usr, cfg)
	delURLBatch := deleteurlbatch.Setup(usr, inWorker)

	joinedRequests := joinRequests(
		pingReqs,
		getStatsReqs,
		getURL,
		postURL,
		getURLsByUserID,
		postURLBatch,
		delURLBatch,
	)

	// Interceptors
	mw := interceptors.New(usr)
	s := grpc.NewServer(grpc.UnaryInterceptor(mw.UserIDInterceptor))

	protoapi.RegisterShortenerServer(s, joinedRequests)

	// SRV GRPC

	return s
}
