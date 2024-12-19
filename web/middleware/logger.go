package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"liewell.fun/alioth/core"
)

func Logger(skipPathArr []string) gin.HandlerFunc {

	skipPaths := make(map[string]bool, len(skipPathArr))
	for _, path := range skipPathArr {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					core.Logger.Error(e)
				}
			} else {
				if c.Request.Method != http.MethodOptions {
					core.Logger.Debugf("%s %s status: %d cost: %v", strings.ToUpper(c.Request.Method), c.Request.RequestURI, c.Writer.Status(), latency)
				}
			}
		}
	}
}
