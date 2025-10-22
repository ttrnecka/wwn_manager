package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

const maxBodyLogSize = 256 // limit body log size

func RequestLogger(logger zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogLatency:       true,
		LogProtocol:      true,
		LogMethod:        true,
		LogRemoteIP:      true,
		LogUserAgent:     true,
		LogContentLength: true,
		LogError:         true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// --- Capture request body ---
			reqBody, _ := c.Get("reqBody").(string)

			// --- Capture response body ---
			respRecorder, ok := c.Response().Writer.(*responseRecorder)
			var respBody string
			if ok {
				respBody = respRecorder.body.String()
				if len(respBody) > maxBodyLogSize {
					respBody = respBody[:maxBodyLogSize] + "...(truncated)"
				}
				respBody = maskSensitive(respBody)
			}

			// --- Build log entry ---
			evt := logger.Info()
			if v.Error != nil {
				evt = logger.Error().Err(v.Error)
			}

			evt.
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("protocol", v.Protocol).
				Str("remote_ip", v.RemoteIP).
				Str("user_agent", v.UserAgent).
				Int("status", v.Status).
				Str("content_length", v.ContentLength).
				Dur("latency", v.Latency).
				Time("start_time", v.StartTime).
				Time("end_time", time.Now()).
				Str("req_body", reqBody).
				Str("resp_body", respBody).
				Msg("http_request")

			return nil
		},
	})
}

func RequestBodyCaptureMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			if req.Body != nil {
				// Read the full body
				bodyBytes, _ := io.ReadAll(req.Body)
				_ = req.Body.Close()

				// Restore the full body for the next handler
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Make a truncated copy for logging
				reqBody := string(bodyBytes)
				if len(reqBody) > maxBodyLogSize {
					reqBody = reqBody[:maxBodyLogSize] + "...(truncated)"
				}
				c.Set("reqBody", maskSensitive(reqBody))
			}
			return next(c)
		}
	}
}

// --- Response recorder middleware ---
// Wrap response writer before RequestLogger runs
func ResponseBodyCaptureMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Wrap the response writer so we can read the body later
			rec := newResponseRecorder(c.Response().Writer)
			c.Response().Writer = rec
			return next(c)
		}
	}
}

type responseRecorder struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// --- Naive sensitive data masker ---
func maskSensitive(s string) string {
	replacements := []string{"password", "secret", "token"}
	for _, key := range replacements {
		s = strings.ReplaceAll(s, key, "[REDACTED]")
	}
	return s
}
