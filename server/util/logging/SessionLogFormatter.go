package logging

import (
	"air-sync/models"

	log "github.com/sirupsen/logrus"
)

type SessionLogFormatter struct {
	log.Formatter
	session *models.Session
}

var _ log.Formatter = (*SessionLogFormatter)(nil)

func NewSessionLogFormatter(fmt log.Formatter, session *models.Session) *SessionLogFormatter {
	return &SessionLogFormatter{
		Formatter: fmt,
		session:   session,
	}
}

func (f *SessionLogFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Data["session_id"] = f.session.Id
	return f.Formatter.Format(e)
}
