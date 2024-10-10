package ginkit

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/web"

	shortuuid "github.com/lithammer/shortuuid/v3"
)

// Middlewares
type MiddlewareFunc func(c *gin.Context)

// RequestIDLoggerMiddleware - adds request ID for incoming http request.
// it also set request with new context with logger
// it's useful when you want to using log package with requestID.
func RequestIDLoggerMiddleware(cfg *RuntimeConfig, o ...Option) MiddlewareFunc {
	opts := options{
		rIDKey: web.HTTPKeyRequestID,
		tIDKey: web.HTTPKeyTraceID,
	}

	for _, o := range o {
		o(&opts)
	}

	m := &rIDLoggerMiddleware{
		rIDKey: opts.rIDKey,
		tIDKey: opts.tIDKey,
		cfg:    cfg,
	}

	return m.handle
}

type rIDLoggerMiddleware struct {
	rIDKey string
	tIDKey string
	cfg    *RuntimeConfig
}

func (m *rIDLoggerMiddleware) handle(ctx *gin.Context) {
	pathName := ctx.Request.URL.Path

	if pathName == m.cfg.HealthCheckPath ||
		pathName == m.cfg.InfoCheckPath ||
		pathName == "/metrics" {
		ctx.Next()
		return
	}

	rCtx := log.NewLoggingContext(ctx.Request.Context())
	logger := log.FromCtx(ctx.Request.Context())

	// trace ID
	tID := ctx.Request.Header.Get(m.tIDKey)

	if tID == "" {
		tID = shortuuid.New()
		ctx.Request.Header.Add(web.HTTPKeyTraceID, tID)

		if web.HTTPKeyTraceID != m.tIDKey {
			ctx.Request.Header.Add(m.tIDKey, tID)
		}
	}

	ctx.Writer.Header().Set(web.HTTPKeyTraceID, tID)
	logger.AddField("trace_id", tID)

	// request ID
	rID := ctx.Request.Header.Get(m.rIDKey)

	if rID == "" {
		rID = shortuuid.New()

		ctx.Request.Header.Add(web.HTTPKeyRequestID, rID)

		if web.HTTPKeyRequestID != m.rIDKey {
			ctx.Request.Header.Add(m.rIDKey, rID)
		}
	}

	ctx.Writer.Header().Set(web.HTTPKeyRequestID, rID)
	logger.AddField("request_id", rID)

	rCtx = log.AddToContext(rCtx, logger)

	ctx.Request = ctx.Request.WithContext(rCtx)

	ctx.Next()
}

type options struct {
	rIDKey string
	tIDKey string
}

// Option sets options for request middleware.
type Option func(*options)

// WithRequestIDKey returns an Option which sets `key` as request-ID lookup
// to use for logging server requests.
func WithRequestIDKey(key string) Option {
	return func(o *options) {
		if key != "" {
			o.rIDKey = key
		}
	}
}

// WithTraceIDKey returns an Option which sets `key` as trace-ID lookup
// to use for logging server requests.
func WithTraceIDKey(key string) Option {
	return func(o *options) {
		if key != "" {
			o.tIDKey = key
		}
	}
}

// BodyDumpHandler logs incoming request & outgoing response body.
func BodyDumpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		reqBody, _ := c.Get(gin.BodyBytesKey)
		respBody := c.Writer.Status()

		l := log.FromCtx(c.Request.Context())
		respStr := strings.Replace(fmt.Sprint(respBody), "}\n", "}", 1)
		reqStr := strings.Replace(string(reqBody.([]byte)), "}\n", "}", 1)

		l.Debug("REST request completed",
			"http.response", respStr,
			"http.request", reqStr,
			"http.status_code", c.Writer.Status(),
			"http.header", c.Request.Header,
		)
	}
}

// func bodyDumpHandlerFunc(c *gin.Context, reqBody []byte, respBody []byte) {
// 	l := log.FromCtx(c.Request.Context())
// 	respStr := strings.Replace(string(respBody), "}\n", "}", 1)
// 	reqStr := strings.Replace(string(reqBody), "}\n", "}", 1)

// 	l.Debug("REST request completed",
// 		"http.response", respStr,
// 		"http.request", reqStr,
// 		"http.status_code", c.Writer.Status(),
// 		"http.header", c.Request.Header,
// 	)
// }

// // TimeoutConfig request timeout configuration
// // default value:
// //   - timeout: 7 seconds
// //   - middleware.DefaultSkipper / apply to all url
// type TimeoutConfig struct {
// 	Timeout time.Duration `json:"timeout,omitempty"`
// 	Skipper func(*gin.Context) bool
// }

// // TimeoutMiddleware sets upstream request context's timeout.
// func TimeoutMiddleware(cfg *TimeoutConfig) gin.HandlerFunc {
// 	// setup default value
// 	if cfg.Timeout <= 0 {
// 		cfg.Timeout = defaultReqTimeout
// 	}

// 	return func(c *gin.Context) {
// 		// If Skipper function is defined and returns true, skip the middleware
// 		if cfg.Skipper != nil && cfg.Skipper(c) {
// 			c.Next()
// 			return
// 		}

// 		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.Timeout)
// 		defer cancel()

// 		c.Request = c.Request.WithContext(ctx)
// 		c.Next()
// 	}
// }
