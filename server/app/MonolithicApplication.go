package app

import (
	"air-sync/handlers"
	"air-sync/services"
	"context"
	"net/url"

	"github.com/gorilla/mux"
)

type MonolithicApplication struct {
	Addr          string
	MongoUrl      *url.URL
	MongoDatabase string
	RedisAddr     string
	RedisPassword string
	EnableCORS    bool
}

var _ Application = (*MonolithicApplication)(nil)

func (s *MonolithicApplication) Start(ctx context.Context) error {
	repos := services.NewMongoRepositoryService(ctx, services.MongoRepositoryOptions{
		URL:      s.MongoUrl,
		Database: s.MongoDatabase,
	})
	if err := repos.Initialize(); err != nil {
		return err
	}
	defer repos.Deinitialize()

	eventBroker := services.NewEventBrokerService(ctx)
	eventBroker.Initialize()
	defer eventBroker.Deinitialize()

	redisBroker := services.NewRedisBrokerService(ctx, services.RedisBrokerOptions{
		Publisher: eventBroker.Publisher(),
		Addr:      s.RedisAddr,
		Password:  s.RedisPassword,
	})
	if err := redisBroker.Initialize(); err != nil {
		return err
	}
	defer redisBroker.Deinitialize()

	router := mux.NewRouter()
	handlers.NewApiHandler(
		handlers.NewSessionRestHandler(
			repos.SessionRepository(),
			eventBroker.Publisher(),
		),
		handlers.QrRestHandler(0),
	).RegisterRoutes(router)
	handlers.NewWebSocketHandler(handlers.WebSocketOptions{
		Repository: repos.SessionRepository(),
		Publisher:  eventBroker.Publisher(),
		EnableCORS: s.EnableCORS,
	}).RegisterRoutes(router)
	handlers.NewStreamingHandler(
		repos.SessionRepository(),
		eventBroker.Publisher(),
	).RegisterRoutes(router)
	handlers.NewLongPollHandler(
		repos.SessionRepository(),
		eventBroker.Publisher(),
	).RegisterRoutes(router)
	handlers.NewWebHandler(handlers.WebOptions{
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
