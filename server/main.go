package main

import (
	appHandlers "air-sync/handlers"
	repos "air-sync/repositories"
	"air-sync/util"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(util.DefaultTextFormatter)
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-ch
		log.Infof("Received signal: %+v", sig)
		cancel()
	}()

	serve(ctx, "127.0.0.1:8080")
}

func serve(ctx context.Context, addr string) {
	repo := repos.NewSessionRepository()
	router := mux.NewRouter()
	appHandlers.NewApiHandler(repo).RegisterRoutes(router)
	appHandlers.NewWebSocketHandler(repo).RegisterRoutes(router)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
		return
	}

	server := &http.Server{
		Handler:      handlers.CORS()(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	go func() {
		if err = server.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Infof("Server listening at %s", addr)
	<-ctx.Done()

	log.Info("Shutting down server")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Failed to shutdown: %+s\n", err)
	} else {
		log.Info("Server shutdown properly")
	}
}
