package middleware

import (
	"bytes"
	"example/web-service-gin/app/clientContext"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogEntry struct {
	Timestamp     time.Time                   `json:"timestamp"`
	ClientContext clientContext.ClientContext `json:"client_context"`
}

func JsonLogger() gin.HandlerFunc {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	return func(c *gin.Context) {

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the original body to avoid affecting subsequent handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process the users request
		c.Next()

		writer := c.Writer

		currentContext := c.Request.Context().Value(clientContext.ClientContextKey).(*clientContext.ClientContext)
		entry := LogEntry{
			Timestamp:     time.Now(),
			ClientContext: *currentContext,
		}

		level := logrus.InfoLevel
		if writer.Status() >= http.StatusInternalServerError {
			level = logrus.ErrorLevel
		}

		// Log the entry as JSON
		logrus.WithFields(logrus.Fields{
			"entry": entry,
		}).Log(level, "Request logged")
	}

}
