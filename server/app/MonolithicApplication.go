package app

import (
	"air-sync/handlers"
	"air-sync/services"
	"context"
	"net/url"
	"time"

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

	StorageMode string
	BucketName  string
	UploadsDir  string

	CronEnvironment string
	GracePeriod     time.Duration

	EnableCORS bool
}

var _ Application = (*MonolithicApplication)(nil)

func (a *MonolithicApplication) Start(ctx context.Context) error {
	repos := services.NewMongoRepositoryService(ctx, services.MongoRepositoryOptions{
		URL:      a.Mongo.URL,
		Database: a.Mongo.Database,
	})
	if err := repos.Initialize(); err != nil {
		return err
	}
	defer repos.Deinitialize()

	storageService := services.NewStorageService(ctx, services.StorageOptions{
		StorageMode: services.StorageMode(a.StorageMode),
		BucketName:  a.BucketName,
		UploadsDir:  a.UploadsDir,
	})
	if err := storageService.Initialize(); err != nil {
		return err
	}
	defer storageService.Deinitialize()

	eventBroker := services.NewEventBrokerService(ctx, services.EventBrokerOptions{
		Service:      services.EventService(a.EventService),
		Redis:        a.Redis,
		GooglePubSub: a.GooglePubSub,
	})
	if err := eventBroker.Initialize(); err != nil {
		return err
	}
	defer eventBroker.Deinitialize()

	cronJobService := services.NewCronJobService(services.CronJobOptions{
		SessionRepository:    repos.SessionRepository(),
		AttachmentRepository: repos.AttachmentRepository(),
		Publisher:            eventBroker.Publisher(),
		Storage:              storageService.Storage(),
		GracePeriod:          a.GracePeriod,
	})
	if err := cronJobService.Initialize(); err != nil {
		return err
	}
	defer cronJobService.Deinitialize()

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
		EnableCORS: a.EnableCORS,
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
		storageService.Storage(),
	).RegisterRoutes(router)

	handlers.NewCronHandler(
		handlers.CronEnvironment(a.CronEnvironment),
		cronJobService,
	).RegisterRoutes(router)

	handlers.NewWebHandler(handlers.WebOptions{
		PublicPath:   "public",
		IndexFile:    "index.html",
		NotFoundFile: "404.html",
	}).RegisterRoutes(router)

	srv := &WebApplication{
		Router:     router,
		Addr:       a.Addr,
		EnableCORS: a.EnableCORS,
	}
	return srv.Start(ctx)
}
