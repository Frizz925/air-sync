package gcp

import (
	"os"
)

func EnvProjectID() string {
	return os.Getenv("GOOGLE_CLOUD_PROJECT")
}

func EnvPubSubTopicID() string {
	return os.Getenv("GOOGLE_PUBSUB_TOPIC_ID")
}

func EnvPubSubSubscriptionID() string {
	return os.Getenv("GOOGLE_PUBSUB_SUBSCRIPTION_ID")
}
