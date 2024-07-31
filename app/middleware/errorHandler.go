package middleware

import (
	"errors"
	"example/web-service-gin/app/apiErrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error apiErrors.APIError `json:"error"`
}

func ErrorHandler(c *gin.Context) {
	c.Next()
	if err := c.Errors.ByType(gin.ErrorTypePrivate).Last(); err != nil {
		var appError *apiErrors.APIError
		if ok := errors.As(err, &appError); ok {
			c.JSON(appError.Status, ErrorResponse{Error: *appError})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		c.Abort()
	}
}
