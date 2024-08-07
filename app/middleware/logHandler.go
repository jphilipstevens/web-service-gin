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

		responseTime := time.Since(startTime)
		writer := c.Writer

		clientContext.AddResponseInfo(c.Request.Context(), clientContext.ResponseInfo{
			Status: writer.Status(),
		})
		clientContext.AddResponseTime(c.Request.Context(), responseTime)
		currentContext := clientContext.GetClientContext(c.Request.Context())

		level := logrus.InfoLevel
		if currentContext.Response.Status >= http.StatusInternalServerError {
			level = logrus.ErrorLevel
		}

		// Log the entry as JSON
		logrus.WithFields(logrus.Fields{
			"clientContext": *currentContext,
		}).Log(level, "Request logged")
	}

}
