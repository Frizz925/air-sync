package app

import (
	"air-sync/handlers"
	repos "air-sync/repositories"
	"context"

	"github.com/gorilla/mux"
)

type MonolithicService struct {
	Addr       string
	EnableCORS bool
}

var _ Service = (*MonolithicService)(nil)

func (s *MonolithicService) Start(ctx context.Context) error {
	repo := repos.NewSessionRepository()
	router := mux.NewRouter()
	handlers.NewApiHandler(repo).RegisterRoutes(router)
	handlers.NewWebSocketHandler(repo, s.EnableCORS).RegisterRoutes(router)
	handlers.NewStreamingHandler(repo).RegisterRoutes(router)
	handlers.NewLongPollHandler(repo).RegisterRoutes(router)
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
