package restlogger

import (
	log "log/slog"
	"net/http"
	"time"

	"github.com/Azure/aks-middleware/logging"
)

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
	Logger  *log.Logger
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := lrt.Proxied.RoundTrip(req)
	methodInfo := logging.GetMethodInfo(req.Method, req.URL.Path)
	if err != nil {
		lrt.Logger.With(
			"code", "na",
			"component", "client",
			"time_ms", "na",
			"method", methodInfo,
			"service", req.Host,
			"source", "ApiAutoLog",
			"protocol", "REST",
			"method_type", "unary",
			"url", req.URL.Path,
			"error", err.Error(),
		).Error("error finishing call")
		return resp, err
	}

	latency := time.Since(start).Milliseconds()
	lrt.Logger.With(
		"code", resp.StatusCode,
		"component", "client",
		"time_ms", latency,
		"method", methodInfo,
		"service", req.Host,
		"source", "ApiAutoLog",
		"protocol", "REST",
		"method_type", "unary",
		"url", req.URL.Path,
	).Info("finished call")

	return resp, err
}

func NewLoggingClient(logger *log.Logger) *http.Client {
	return &http.Client{
		Transport: &LoggingRoundTripper{
			Proxied: http.DefaultTransport,
			Logger:  logger,
		},
	}
}
