package util

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type RequestLogFormatter struct {
	log.Formatter
	*http.Request
}

var _ log.Formatter = (*RequestLogFormatter)(nil)

var DefaultTextFormatter = &log.TextFormatter{
	FullTimestamp: true,
}

func NewRequestLogFormatter(fmt log.Formatter, req *http.Request) *RequestLogFormatter {
	return &RequestLogFormatter{
		Formatter: fmt,
		Request:   req,
	}
}

func (f *RequestLogFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Data["client"] = f.Request.RemoteAddr
	return f.Formatter.Format(e)
}
