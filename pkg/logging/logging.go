package logging

import (
	"github.com/sirupsen/logrus"
)

// Logger - структура для представления логгера
type Logger struct {
	log *logrus.Logger
}

// NewLogger - создает новый экземпляр логгера
func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return &Logger{
		log: logger,
	}
}

// GetLogger - возвращает текущий логгер
func (l *Logger) GetLogger() *logrus.Logger {
	return l.log
}
