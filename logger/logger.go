package logger

import "github.com/sirupsen/logrus"

var (
	// Log is the general application logger.
	Log *logrus.Logger
)

func init() {
	Log = logrus.New()
}
