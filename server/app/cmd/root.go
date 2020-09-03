package cmd

import (
	"air-sync/app"
	"air-sync/services"
	"air-sync/util"
	"air-sync/util/gcp"
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

		mongoUrl, err := util.EnvMongoUrl()
		if err != nil {
			log.Fatal(err)
			return
		}

		err = (&app.MonolithicApplication{
			Addr: ":" + util.EnvPort(),
			Mongo: app.MongoOptions{
				URL:      mongoUrl,
				Database: util.EnvMongoDatabase(),
			},
			Redis: services.RedisOptions{
				Addr:     util.EnvRedisAddr(),
				Password: util.EnvRedisPassword(),
			},
			GooglePubSub: services.GooglePubSubOptions{
				ProjectID:      gcp.EnvProjectID(),
				TopicID:        gcp.EnvPubSubTopicID(),
				SubscriptionID: gcp.EnvPubSubSubscriptionID(),
			},
			EventService: util.EnvEventService(),
			EnableCORS:   enableCORS,
		}).Start(ctx)
		if err != nil {
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
