package cmd

import (
	"air-sync/app"
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	enableCORS bool
)

var rootCmd = &cobra.Command{
	Use:   "air-sync",
	Short: "Air Sync is a small web application to quickly send messages over the internet",
	Long: `Small and lightweight, Air Sync is aimed to send various messages 
		from one device to the other securely over the internet.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		srv := &app.MonolithicService{
			Addr:       ":" + port,
			EnableCORS: enableCORS,
		}
		if err := srv.Start(ctx); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&enableCORS, "cors", "c", false, "Enable CORS headers on the server")
}

func Execute() error {
	return rootCmd.Execute()
}
