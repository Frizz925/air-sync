package cmd

import (
	"air-sync/handlers"
	repos "air-sync/repositories"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

type ServerOptions struct {
	Addr       string
	PublicPath string
}

var (
	publicPath string
	indexFile  string
)

var rootCmd = &cobra.Command{
	Use:   "air-sync",
	Short: "Air Sync is a small web application to quickly send messages over the internet",
	Long: `Small and lightweight, Air Sync is aimed to send various messages 
		from one device to the other securely over the internet.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			sig := <-ch
			log.Infof("Received signal: %+v", sig)
			cancel()
		}()

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
			log.Infof("Defaulting to port %s", port)
		}

		return serve(ctx, ServerOptions{
			Addr:       ":" + port,
			PublicPath: publicPath,
		})
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&publicPath, "public-path", "public", "The public directory for serving the web page")
}

func Execute() error {
	return rootCmd.Execute()
}

func serve(ctx context.Context, options ServerOptions) error {
	repo := repos.NewSessionRepository()
	router := mux.NewRouter()
	handlers.NewApiHandler(repo).RegisterRoutes(router)
	handlers.NewWebSocketHandler(repo).RegisterRoutes(router)
	handlers.NewWebHandler(handlers.WebHandlerOptions{
		PublicPath:   options.PublicPath,
		IndexFile:    "index.html",
		NotFoundFile: "404.html",
	}).RegisterRoutes(router)

	l, err := net.Listen("tcp4", options.Addr)
	if err != nil {
		return err
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	server := &http.Server{
		Handler:      c.Handler(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := server.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Infof("Server listening at %s", options.Addr)
	<-ctx.Done()
	log.Info("Shutting down server")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		return err
	}

	log.Info("Server shutdown properly")
	return nil
}
