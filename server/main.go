package main

import (
	"air-sync/app/cmd"
	"air-sync/util"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(util.DefaultTextFormatter)
}

func main() {
	util.SetupCredentialsEnv()
	if err := util.LoadDotEnv(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
