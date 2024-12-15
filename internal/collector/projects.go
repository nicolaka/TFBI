package collector

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-tfe"
	"github.com/nicolaka/tfbi/internal/setup"
	"golang.org/x/sync/errgroup"

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
		[]string{"id", "name", "organization", "description"}, nil,
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

func getProjectsListPage(ctx context.Context, page int, organization string, config *setup.Config, ch chan<- prometheus.Metric) error {
	projectsList, err := config.Client.Projects.List(ctx, organization, &tfe.ProjectListOptions{
		ListOptions: tfe.ListOptions{
			PageSize:   pageSize,
			PageNumber: page,
		},
	})

	if err != nil {
		return fmt.Errorf("%v, (organization=%s, page=%d)", err, organization, page)
	}

	for _, p := range projectsList.Items {
		select {
		case ch <- prometheus.MustNewConstMetric(
			ProjectsInfo,
			prometheus.GaugeValue,
			1,
			p.ID,
			p.Name,
			p.Organization.Name,
			p.Description,
		):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (ScrapeProjects) Scrape(ctx context.Context, config *setup.Config, ch chan<- prometheus.Metric) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, name := range config.Organizations {
		name := name
		g.Go(func() error {
			projectsList, err := config.Client.Projects.List(ctx, name, &tfe.ProjectListOptions{
				ListOptions: tfe.ListOptions{
					PageSize: pageSize,
				}})

			if err != nil {
				return fmt.Errorf("%v, organization=%s", err, name)
			}

			for i := 1; i <= projectsList.Pagination.TotalPages; i++ {
				if err := getProjectsListPage(ctx, i, name, config, ch); err != nil {
					return err
				}
			}

			return nil
		})
	}

	return g.Wait()
}
