package collector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-tfe"
	"github.com/nicolaka/tfbi/internal/setup"
	"golang.org/x/sync/errgroup"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// modules is the Metric subsystem we use.
	registrymodulesSubsystem = "registrymodules"
)

// Metric descriptors.
var (
	RegistryModulesInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, registrymodulesSubsystem, "info"),
		"Information about existing registrymodules",
		[]string{"id", "name", "provider", "registry_name", "no_code", "status", "created_at", "updated_at", "organization"}, nil,
	)
)

// ScrapeRegistryModules scrapes metrics about the registrymodules.
type ScrapeRegistryModules struct{}

func init() {
	Scrapers = append(Scrapers, ScrapeRegistryModules{})
}

// Name of the Scraper. Should be unique.
func (ScrapeRegistryModules) Name() string {
	return registrymodulesSubsystem
}

// Help describes the role of the Scraper.
func (ScrapeRegistryModules) Help() string {
	return "Scrape information from the Registry Modules API: https://www.terraform.io/docs/cloud/api/modules.html"
}

// Version of Terraform Cloud/Enterprise API from which scraper is available.
func (ScrapeRegistryModules) Version() string {
	return "v2"
}

// []string{"id", "name", "provider", "registry-name","no-code", "status", "created-at","updated-at"}, nil,

func getModulesListPage(ctx context.Context, page int, organization string, config *setup.Config, ch chan<- prometheus.Metric) error {
	registrymodulesList, err := config.Client.RegistryModules.List(ctx, organization, &tfe.RegistryModuleListOptions{
		ListOptions: tfe.ListOptions{
			PageSize:   pageSize,
			PageNumber: page,
		},
	})

	if err != nil {
		return fmt.Errorf("%v, (organization=%s, page=%d)", err, organization, page)
	}

	for _, m := range registrymodulesList.Items {
		select {
		case ch <- prometheus.MustNewConstMetric(
			RegistryModulesInfo,
			prometheus.GaugeValue,
			1,
			m.ID,
			m.Name,
			m.Provider,
			string(m.RegistryName),
			strconv.FormatBool(m.NoCode),
			string(m.Status),
			m.CreatedAt,
			m.UpdatedAt,
			m.Organization.Name,
		):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (ScrapeRegistryModules) Scrape(ctx context.Context, config *setup.Config, ch chan<- prometheus.Metric) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, name := range config.Organizations {
		name := name
		g.Go(func() error {
			registrymodulesList, err := config.Client.RegistryModules.List(ctx, name, &tfe.RegistryModuleListOptions{
				ListOptions: tfe.ListOptions{
					PageSize: pageSize,
				}})

			if err != nil {
				return fmt.Errorf("%v, organization=%s", err, name)
			}

			for i := 1; i <= registrymodulesList.Pagination.TotalPages; i++ {
				if err := getModulesListPage(ctx, i, name, config, ch); err != nil {
					return err
				}
			}

			return nil
		})
	}

	return g.Wait()
}
