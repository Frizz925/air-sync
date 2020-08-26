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
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
