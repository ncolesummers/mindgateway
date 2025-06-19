package logging

import (
	"os"
	
	"github.com/sirupsen/logrus"
)

// Logger extends logrus.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new configured logger
func NewLogger(level string) *Logger {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logLevel)
	
	// Use JSON formatter for structured logging
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	
	return &Logger{Logger: l}
}

// WithComponent adds component information to the logger
func (l *Logger) WithComponent(component string) *logrus.Entry {
	return l.WithField("component", component)
}

// WithRequest adds request information to the logger
func (l *Logger) WithRequest(requestID string, method string, path string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"request_id": requestID,
		"method":     method,
		"path":       path,
	})
}

// WithUser adds user information to the logger
func (l *Logger) WithUser(userID string) *logrus.Entry {
	return l.WithField("user_id", userID)
}

// WithWorker adds worker information to the logger
func (l *Logger) WithWorker(workerID string) *logrus.Entry {
	return l.WithField("worker_id", workerID)
}