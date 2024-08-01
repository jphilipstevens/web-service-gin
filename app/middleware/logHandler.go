package middleware

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogEntry struct {
	Timestamp   time.Time     `json:"timestamp"`
	ClientIP    string        `json:"client_ip"`
	Method      string        `json:"method"`
	Path        string        `json:"path"`
	Protocol    string        `json:"protocol"`
	UserAgent   string        `json:"user_agent"`
	StatusCode  int           `json:"status_code"`
	Latency     time.Duration `json:"latency"`
	RequestBody string        `json:"request_body"`
}

func JsonLogger() gin.HandlerFunc {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the original body to avoid affecting subsequent handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process the users request
		c.Next()

		writer := c.Writer

		latency := time.Since(startTime)
		entry := LogEntry{
			Timestamp:   time.Now(),
			ClientIP:    c.ClientIP(),
			Method:      c.Request.Method,
			Path:        c.Request.URL.Path,
			Protocol:    c.Request.Proto,
			UserAgent:   c.Request.UserAgent(),
			StatusCode:  writer.Status(),
			Latency:     latency,
			RequestBody: string(requestBody),
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
