package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/jphilipstevens/web-service-gin/v2/pkg/clientContext"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// JsonLogger records structured details about each HTTP request and response.
// It assumes the global logger has been initialized elsewhere via SetupLogger.
func JsonLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logrus.WithError(err).Warn("Failed to read request body")
			} else {
				requestBody = body
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if len(requestBody) > 1024 {
				requestBody = []byte("<<truncated>>")
			}
		}

		c.Next()

		responseTime := time.Since(startTime)
		writer := c.Writer

		clientContext.AddResponseInfo(c.Request.Context(), clientContext.ResponseInfo{
			Status: writer.Status(),
		})
		clientContext.AddResponseTime(c.Request.Context(), responseTime)
		currentContext := clientContext.GetClientContext(c.Request.Context())

		status := writer.Status()
		level := logrus.InfoLevel
		switch {
		case status >= http.StatusInternalServerError:
			level = logrus.ErrorLevel
		case status >= http.StatusBadRequest:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		fields := logrus.Fields{
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"status":        status,
			"latency":       responseTime.String(),
			"ip":            c.ClientIP(),
			"userAgent":     c.Request.UserAgent(),
			"clientContext": *currentContext,
		}

		var skipBodyPaths = map[string]bool{
			"/auth/login":    true,
			"/auth/register": true,
		}
		if len(requestBody) > 0 && !skipBodyPaths[c.Request.URL.Path] {
			fields["requestBody"] = string(requestBody)
		}

		logrus.WithFields(fields).Log(level, "Request logged")
	}
}
