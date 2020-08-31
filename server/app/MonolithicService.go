package app

import (
	"air-sync/handlers"
	repos "air-sync/repositories"
	"air-sync/subscribers"
	"air-sync/util/pubsub"
	"context"

	"github.com/gorilla/mux"
)

type MonolithicService struct {
	Addr       string
	EnableCORS bool
}

var _ Service = (*MonolithicService)(nil)

func (s *MonolithicService) Start(ctx context.Context) error {
	repo := repos.NewSessionLocalRepository()
	stream := pubsub.NewStream()
	defer stream.Shutdown()
	subscribers.SubscribeSession(stream)

	router := mux.NewRouter()
	handlers.NewApiHandler(
		handlers.NewSessionRestHandler(repo, stream),
		handlers.QrRestHandler(0),
	).RegisterRoutes(router)
	handlers.NewWebSocketHandler(handlers.WebSocketHandlerOptions{
		Repository: repo,
		Stream:     stream,
		EnableCORS: s.EnableCORS,
	}).RegisterRoutes(router)
	handlers.NewStreamingHandler(repo, stream).RegisterRoutes(router)
	handlers.NewLongPollHandler(repo, stream).RegisterRoutes(router)
	handlers.NewWebHandler(handlers.WebHandlerOptions{
		PublicPath:   "public",
		IndexFile:    "index.html",
		NotFoundFile: "404.html",
	}).RegisterRoutes(router)

	srv := &WebService{
		Router:     router,
		Addr:       s.Addr,
		EnableCORS: s.EnableCORS,
	}
	return srv.Start(ctx)
}
