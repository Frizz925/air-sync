package logging

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type RequestLogFormatter struct {
	log.Formatter
	request *http.Request
}

var _ log.Formatter = (*RequestLogFormatter)(nil)

func NewRequestLogFormatter(fmt log.Formatter, req *http.Request) *RequestLogFormatter {
	return &RequestLogFormatter{
		Formatter: fmt,
		request:   req,
	}
}

func (f *RequestLogFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Data["client"] = f.request.RemoteAddr
	return f.Formatter.Format(e)
}
