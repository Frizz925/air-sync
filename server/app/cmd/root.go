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

		gracePeriod, err := util.ParseTimeDuration(util.GetEnvDefault("CLEANUP_GRACE_PERIOD", "0s"))
		if err != nil {
			log.Fatal(err)
			return
		}

		err = (&app.MonolithicApplication{
			Addr: ":" + util.GetEnvDefault("PORT", "8080"),
			Mongo: app.MongoOptions{
				URL:      mongoUrl,
				Database: util.GetEnvDefault("MONGODB_DATABASE", "airsync"),
				Recreate: util.GetEnvBoolDefault("MONGODB_RECREATE", false),
			},
			StorageMode: util.GetEnvDefault("STORAGE_MODE", "local"),
			BucketName:  util.GetEnvDefault("BUCKET_NAME", "airsync"),
			UploadsDir:  util.GetEnvDefault("UPLOADS_DIR", "uploads"),
			Redis: services.RedisOptions{
				Addr:     util.GetEnvDefault("REDIS_ADDR", "localhost:6379"),
				Password: util.GetEnvDefault("REDIS_PASSWORD", ""),
			},
			GooglePubSub: services.GooglePubSubOptions{
				ProjectID:      gcp.EnvProjectID(),
				TopicID:        gcp.EnvPubSubTopicID(),
				SubscriptionID: gcp.EnvPubSubSubscriptionID(),
			},
			EventService:    util.GetEnvDefault("EVENT_SERVICE", ""),
			CronEnvironment: util.GetEnvDefault("CRON_ENVIRONMENT", ""),
			GracePeriod:     gracePeriod,
			EnableCORS:      enableCORS,
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
