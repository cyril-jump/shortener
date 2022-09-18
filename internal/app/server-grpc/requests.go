package servergrpc

import (
	"context"

	protoapi "github.com/cyril-jump/shortener/pkg/api/proto"
)

type PingMethod interface {
	Ping(
		context.Context, *protoapi.PingRequest,
	) (
		*protoapi.PingResponse, error,
	)
}

type GetStatsMethod interface {
	GetStats(
		context.Context, *protoapi.GetStatsRequest,
	) (
		*protoapi.GetStatsResponse, error,
	)
}

type GetURLMethod interface {
	GetURL(
		context.Context, *protoapi.GetURLRequest,
	) (
		*protoapi.GetURLResponse, error,
	)
}

type PostURLMethod interface {
	PostURL(
		context.Context, *protoapi.PostURLRequest,
	) (
		*protoapi.PostURLResponse, error,
	)
}

type GetURLsByUserIDMethod interface {
	GetURLsByUserID(
		context.Context, *protoapi.GetURLsByUserIDRequest,
	) (
		*protoapi.GetURLsByUserIDResponse, error,
	)
}

type PostURLBatchMethod interface {
	PostURLBatch(
		context.Context, *protoapi.PostURLBatchRequest,
	) (
		*protoapi.PostURLBatchResponse, error,
	)
}

type DeleteURLBatchMethod interface {
	DeleteURLBatch(
		context.Context, *protoapi.DeleteURLBatchRequest,
	) (
		*protoapi.DeleteURLBatchResponse, error,
	)
}

func joinRequests(
	pingMethod PingMethod,
	getStatsMethod GetStatsMethod,
	getURLMethod GetURLMethod,
	postURLMethod PostURLMethod,
	getURLsByUserIDMethod GetURLsByUserIDMethod,
	postURLBatchMethod PostURLBatchMethod,
	deleteURLBatchMethod DeleteURLBatchMethod,
) protoapi.ShortenerServer {
	return &struct {
		PingMethod
		GetStatsMethod
		GetURLMethod
		PostURLMethod
		GetURLsByUserIDMethod
		PostURLBatchMethod
		DeleteURLBatchMethod
	}{
		PingMethod:            pingMethod,
		GetStatsMethod:        getStatsMethod,
		GetURLMethod:          getURLMethod,
		PostURLMethod:         postURLMethod,
		GetURLsByUserIDMethod: getURLsByUserIDMethod,
		PostURLBatchMethod:    postURLBatchMethod,
		DeleteURLBatchMethod:  deleteURLBatchMethod,
	}
}
