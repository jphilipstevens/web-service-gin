package middleware

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var setupOnce sync.Once

// SetupLogger configures the global logrus logger for structured JSON output.
// It is safe to call multiple times; configuration only occurs once.
func SetupLogger() {
	setupOnce.Do(func() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stdout)
		logrus.SetReportCaller(true)
	})
}
