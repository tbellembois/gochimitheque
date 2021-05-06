package logger

import "github.com/sirupsen/logrus"

var (
	// Log is the general application logger
	Log *logrus.Logger
	// LogInternal is the application logger used to log fatal errors
	LogInternal *logrus.Logger
)

func init() {

	Log = logrus.New()
	LogInternal = logrus.New()

}
