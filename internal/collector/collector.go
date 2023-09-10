// Package collector includes all individual collectors to gather and export system metrics.
package collector

import (
	"context"
	"sync"
	"time"

	"github.com/kaizendorks/terraform-cloud-exporter/internal/setup"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Namespace defines the common namespace to be used by all metrics.
	namespace = "tf"
	// Subsystem(s).
	exporter = "exporter"
)

// Exporter collects TF metrics. It implements the prometheus.Collector interface.
type Exporter struct {
	ctx      context.Context
	logger   log.Logger
	config   setup.Config
	scrapers []Scraper
	metrics  Metrics
}

// Metrics represents exporter metrics which values can be carried between http requests.
type Metrics struct {
	TotalScrapes prometheus.Counter
	ScrapeErrors *prometheus.CounterVec
	Error        prometheus.Gauge
}

var (
	// scrapers lists all possible collection methods.
	Scrapers = []Scraper{}
	// Metric descriptors.
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, exporter, "collector_duration_seconds"),
		"Collector time duration.",
		[]string{"collector"}, nil,
	)
)

// New returns a new Terraform API exporter for the provided Config.
func New(ctx context.Context, config setup.Config, metrics Metrics) *Exporter {
	return &Exporter{
		ctx:      ctx,
		logger:   config.Logger,
		config:   config,
		scrapers: Scrapers,
		metrics:  metrics,
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metrics.TotalScrapes.Desc()
	ch <- e.metrics.Error.Desc()
	e.metrics.ScrapeErrors.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(e.ctx, ch)

	ch <- e.metrics.TotalScrapes
	ch <- e.metrics.Error
	e.metrics.ScrapeErrors.Collect(ch)
}

func (e *Exporter) scrape(ctx context.Context, ch chan<- prometheus.Metric) {
	e.metrics.TotalScrapes.Inc()
	if len(e.config.Organizations) == 0 {
		// Note: At some point this will return a paginated response.
		oo, err := e.config.Client.Organizations.List(ctx, nil)
		if err != nil {
			e.metrics.Error.Set(1)
			level.Error(e.logger).Log("msg", "Unable to List Organizations", "err", err)
			return
		}

		for _, o := range oo.Items {
			e.config.Organizations = append(e.config.Organizations, o.Name)
		}
	}

	e.metrics.Error.Set(0)

	var wg sync.WaitGroup
	defer wg.Wait()
	for _, scraper := range e.scrapers {
		wg.Add(1)
		go func(scraper Scraper) {
			defer wg.Done()
			label := "collect." + scraper.Name()
			scrapeTime := time.Now()
			if err := scraper.Scrape(ctx, &e.config, ch); err != nil {
				level.Error(e.logger).Log("msg", "Error from scraper", "scraper", scraper.Name(), "err", err)
				e.metrics.ScrapeErrors.WithLabelValues(label).Inc()
				e.metrics.Error.Set(1)
			}
			ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), label)
		}(scraper)
	}
}

// NewMetrics creates new Metrics instance.
func NewMetrics() Metrics {
	return Metrics{
		TotalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "scrapes_total",
			Help:      "Total number of times the Terraform API was scraped for metrics.",
		}),
		ScrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "scrape_errors_total",
			Help:      "Total number of times an error occurred scraping the Terraform API.",
		}, []string{"collector"}),
		Error: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from Terraform API resulted in an error (1 for error, 0 for success).",
		}),
	}
}
