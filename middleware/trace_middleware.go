package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"
)

const CorrelationIDHeader = "X-Correlation-ID"

func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Add the correlation ID to the context using zerolog context
		ctx := log.NewLoggingContext(c.Request.Context(), "correlation_id", correlationID, "ip_address", c.ClientIP())
		c.Request = c.Request.WithContext(ctx)

		// Set the correlation ID in the response header
		c.Writer.Header().Set(CorrelationIDHeader, correlationID)

		// Capture start time
		startTime := time.Now()

		// Read and parse the request body
		var bodyMap map[string]interface{}
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body for the next handler
				json.Unmarshal(bodyBytes, &bodyMap)
			}
		}

		// Log request details
		log.FromCtx(c.Request.Context()).Info(fmt.Sprintf("Method: %s, Path: %s, Headers: %v, Query: %v, Body: %v", c.Request.Method, c.Request.URL.Path, c.Request.Header, c.Request.URL.Query(), bodyMap))

		// Process request
		c.Next()

		// Log response details
		log.FromCtx(c.Request.Context()).Info(fmt.Sprintf("Method: %s, Path: %s, Status: %d, Duration: %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(startTime)))
	}
}
