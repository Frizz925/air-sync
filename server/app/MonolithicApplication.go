package app

import (
	"air-sync/handlers"
	"air-sync/services"
	"air-sync/storages"
	"context"
	"net/url"

	"github.com/gorilla/mux"
)

type MongoOptions struct {
	URL      *url.URL
	Database string
}

type MonolithicApplication struct {
	Addr string

	Mongo MongoOptions

	Redis        services.RedisOptions
	GooglePubSub services.GooglePubSubOptions
	EventService string

	BucketName string
	UploadsDir string

	EnableCORS bool
}

var _ Application = (*MonolithicApplication)(nil)

func (s *MonolithicApplication) Start(ctx context.Context) error {
	repos := services.NewMongoRepositoryService(ctx, services.MongoRepositoryOptions{
		URL:      s.Mongo.URL,
		Database: s.Mongo.Database,
	})
	if err := repos.Initialize(); err != nil {
		return err
	}
	defer repos.Deinitialize()

	storageService := services.NewStorageService(ctx, services.StorageOptions{
		BucketName: s.BucketName,
		UploadsDir: s.UploadsDir,
	})
	if err := storageService.Initialize(); err != nil {
		return err
	}
	defer storageService.Deinitialize()
	// HACK: Not calling initialize for our cache storage
	// because the underlying storages' lifecycle have
	// already been managed by the storage service itself
	storage := storages.NewCacheStorage(
		storageService.FileStorage(),
		storageService.CloudStorage(),
	)

	eventBroker := services.NewEventBrokerService(ctx, services.EventBrokerOptions{
		Service:      s.EventService,
		Redis:        s.Redis,
		GooglePubSub: s.GooglePubSub,
	})
	if err := eventBroker.Initialize(); err != nil {
		return err
	}
	defer eventBroker.Deinitialize()

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

	handlers.NewAttachmentHandler(
		repos.AttachmentRepository(),
		storage,
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
