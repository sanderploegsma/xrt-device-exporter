package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

func newHandler(logger *slog.Logger) http.Handler {
	collector := NewCollector(logger)
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:            slog.NewLogLogger(logger.Handler(), slog.LevelError),
		ErrorHandling:       promhttp.ContinueOnError,
		MaxRequestsInFlight: 1,
	})
}

func main() {
	var (
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()

		toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":9101")
	)
	logconfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, logconfig)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.Parse()
	logger := promslog.New(logconfig)

	logger.Info("Starting xrt-device-exporter")

	http.Handle(*metricsPath, newHandler(logger))

	server := &http.Server{}
	if err := web.ListenAndServe(server, toolkitFlags, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
