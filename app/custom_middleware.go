package app

import "github.com/gin-gonic/gin"

// userMiddleware holds middleware functions registered by client applications.
// These run after the built-in middleware defined in RunServer.
var userMiddleware []gin.HandlerFunc

// UseMiddleware registers a middleware to execute for every request.
// It must be called before RunServer to take effect.
func UseMiddleware(mw gin.HandlerFunc) {
	userMiddleware = append(userMiddleware, mw)
}
