package main

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/nicolaka/tfbi/internal/collector"
	"github.com/nicolaka/tfbi/internal/setup"

	"github.com/go-kit/kit/log/level"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Build information. Populated at build-time via ldflags.
var (
	Version   string
	Commit    string
	GoVersion = runtime.Version()
	BuildDate string
)

func newHandler(metrics collector.Metrics, config setup.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use request context for cancellation when connection gets closed.
		ctx := r.Context()
		// If a timeout is configured via the Prometheus header, add it to the context.
		if v := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds"); v != "" {
			timeoutSeconds, err := strconv.ParseFloat(v, 64)
			if err != nil {
				level.Error(config.Logger).Log("msg", "Failed to parse timeout from Prometheus header", "err", err)
			} else {
				// Create new timeout context with request context as parent.
				ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds*float64(time.Second)))
				defer cancel()
				// Overwrite request with timeout context.
				r = r.WithContext(ctx)
			}
		}

		registry := prometheus.NewRegistry()
		registry.MustRegister(collector.New(ctx, config, metrics))

		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}
		// Delegate http serving to Prometheus client library, which will call collector.Collect.
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func main() {
	config := setup.NewConfig()
	level.Info(config.Logger).Log("msg", "Starting tf_exporter", "version", Version)
	level.Debug(config.Logger).Log("msg", "Build Context", "go", GoVersion, "date", BuildDate)

	handlerFunc := newHandler(collector.NewMetrics(), config)
	http.Handle("/metrics", promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Terraform Cloud Business Insights </title></head>
			<body>
			<h1>Terraform Cloud Business Insights</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	level.Info(config.Logger).Log("msg", "Listening on address", "address", config.ListenAddress)
	if err := http.ListenAndServe(config.ListenAddress, nil); err != nil {
		level.Error(config.Logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
