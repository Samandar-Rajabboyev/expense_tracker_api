package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		event := log.Info().
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency_ms", latency).
			Str("client_ip", clientIP)

		if status >= 500 {
			event = log.Error().
				Str("method", method).
				Str("path", path).
				Int("status", status).
				Dur("latency_ms", latency).
				Str("client_ip", clientIP)
		} else if status >= 400 {
			event = log.Warn().
				Str("method", method).
				Str("path", path).
				Int("status", status).
				Dur("latency_ms", latency).
				Str("client_ip", clientIP)
		}
		event.Msg("HTTP request")
	}
}
