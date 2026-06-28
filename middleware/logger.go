package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		uid, _ := c.Get(CtxUserID)
		log.Printf("[%s] %s %s | uid=%v | %d | %s | %q",
			start.Format("15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			uid,
			c.Writer.Status(),
			time.Since(start),
			string(bodyBytes),
		)
	}
}
