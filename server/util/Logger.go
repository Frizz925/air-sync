package util

import (
	log "github.com/sirupsen/logrus"
)

var DefaultTextFormatter = &log.TextFormatter{
	FullTimestamp: true,
}
