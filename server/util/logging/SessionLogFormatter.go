package logging

import (
	"air-sync/repositories/entities"

	log "github.com/sirupsen/logrus"
)

type SessionLogFormatter struct {
	log.Formatter
	session entities.Session
}

var _ log.Formatter = (*SessionLogFormatter)(nil)

func NewSessionLogFormatter(fmt log.Formatter, session entities.Session) *SessionLogFormatter {
	return &SessionLogFormatter{
		Formatter: fmt,
		session:   session,
	}
}

func (f *SessionLogFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Data["session_id"] = f.session.ID
	return f.Formatter.Format(e)
}
