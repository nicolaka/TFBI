package collector

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/kaizendorks/terraform-cloud-exporter/internal/setup"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// projects is the Metric subsystem we use.
	projectsSubsystem = "projects"
)

// Metric descriptors.
var (
	ProjectsInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, projectsSubsystem, "info"),
		"Information about existing projects",
		[]string{"id", "name"}, nil,
	)
)

// ScrapeProjects scrapes metrics about the projects.
type ScrapeProjects struct{}

func init() {
	Scrapers = append(Scrapers, ScrapeProjects{})
}

// Name of the Scraper. Should be unique.
func (ScrapeProjects) Name() string {
	return projectsSubsystem
}

// Help describes the role of the Scraper.
func (ScrapeProjects) Help() string {
	return "Scrape information from the Projects API: https://www.terraform.io/docs/cloud/api/projects.html"
}

// Version of Terraform Cloud/Enterprise API from which scraper is available.
func (ScrapeProjects) Version() string {
	return "v2"
}

func (ScrapeProjects) Scrape(ctx context.Context, config *setup.Config, ch chan<- prometheus.Metric) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, name := range config.Organizations {
		g.Go(func() error {
			projectList, err := config.Client.Projects.List(ctx, name, nil)
			for _, p := range projectList.Items {
					fmt.Println(p.Name)
					fmt.Println(p.ID)
					select {
					case ch <- prometheus.MustNewConstMetric(
						ProjectsInfo,
						prometheus.GaugeValue,
						1,
						p.ID,
						p.Name,
					):
					case <-ctx.Done():
						return ctx.Err()
					}
			}

			if err != nil {
				return fmt.Errorf("%v, organization=%s", err, name)
			}

			return nil
		})
		
	}
	return g.Wait()	
}