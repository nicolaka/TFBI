package collector

import (
	"context"

	"github.com/nicolaka/tfbi/internal/setup"

	"github.com/prometheus/client_golang/prometheus"
)

// Scraper is minimal interface that let's you add new prometheus metrics to tf_exporter.
type Scraper interface {
	// Name of the Scraper. Should be unique.
	Name() string

	// Help describes the role of the Scraper.
	// Example: "Collect Workspaces"
	Help() string

	// Version of Terraform Cloud/Enterprise API from which scraper is available.
	Version() string

	// Scrape collects data from a particular terraform cloud/enterprise API and sends it over channel as prometheus metric.
	Scrape(ctx context.Context, config *setup.Config, ch chan<- prometheus.Metric) error
}
