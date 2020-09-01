package app

import (
	"air-sync/handlers"
	"air-sync/services"
	"air-sync/subscribers"
	"air-sync/util/pubsub"
	"context"

	"github.com/gorilla/mux"
)

type MonolithicApplication struct {
	Addr       string
	MongoUri   string
	EnableCORS bool
}

var _ Application = (*MonolithicApplication)(nil)

func (s *MonolithicApplication) Start(ctx context.Context) error {
	repoSrv := services.NewMongoRepositoryService(s.MongoUri, "airsync")
	if err := repoSrv.Initialize(); err != nil {
		return err
	}
	defer repoSrv.Deinitialize()

	stream := pubsub.NewStream()
	defer stream.Shutdown()
	subscribers.SubscribeSession(stream)

	router := mux.NewRouter()
	handlers.NewApiHandler(
		handlers.NewSessionRestHandler(repoSrv.SessionRepository(), stream),
		handlers.QrRestHandler(0),
	).RegisterRoutes(router)
	handlers.NewWebSocketHandler(handlers.WebSocketHandlerOptions{
		Repository: repoSrv.SessionRepository(),
		Stream:     stream,
		EnableCORS: s.EnableCORS,
	}).RegisterRoutes(router)
	handlers.NewStreamingHandler(repoSrv.SessionRepository(), stream).RegisterRoutes(router)
	handlers.NewLongPollHandler(repoSrv.SessionRepository(), stream).RegisterRoutes(router)
	handlers.NewWebHandler(handlers.WebHandlerOptions{
		PublicPath:   "public",
		IndexFile:    "index.html",
		NotFoundFile: "404.html",
	}).RegisterRoutes(router)

	srv := &WebApplication{
		Router:     router,
		Addr:       s.Addr,
		EnableCORS: s.EnableCORS,
	}
	return srv.Start(ctx)
}
