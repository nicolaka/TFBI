package collector

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/nicolaka/tfbi/internal/setup"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// organizations is the Metric subsystem we use.
	organizationsSubsystem = "organizations"
)

// Metric descriptors.
var (
	OrganizationsInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, organizationsSubsystem, "info"),
		"Information about existing organizations",
		[]string{"name", "created_at", "email", "external_id", "owners_team_saml_role_id", "saml_enabled", "two_factor_conformant","assessment_enforced"}, nil,
	)
)

// ScrapeOrganizations scrapes metrics about the organizations.
type ScrapeOrganizations struct{}

func init() {
	Scrapers = append(Scrapers, ScrapeOrganizations{})
}

// Name of the Scraper. Should be unique.
func (ScrapeOrganizations) Name() string {
	return organizationsSubsystem
}

// Help describes the role of the Scraper.
func (ScrapeOrganizations) Help() string {
	return "Scrape information from the Organizations API: https://www.terraform.io/docs/cloud/api/organizations.html"
}

// Version of Terraform Cloud/Enterprise API from which scraper is available.
func (ScrapeOrganizations) Version() string {
	return "v2"
}

func getOrganization(ctx context.Context, name string, config *setup.Config, ch chan<- prometheus.Metric) error {
	o, err := config.Client.Organizations.Read(ctx, name)
	if err != nil {
		return fmt.Errorf("%v, organization=%s", err, name)
	}

	select {
	case ch <- prometheus.MustNewConstMetric(
		OrganizationsInfo,
		prometheus.GaugeValue,
		1,
		o.Name,
		o.CreatedAt.String(),
		o.Email,
		o.ExternalID,
		o.OwnersTeamSAMLRoleID,
		strconv.FormatBool(o.SAMLEnabled),
		strconv.FormatBool(o.TwoFactorConformant),
		strconv.FormatBool(o.AssessmentsEnforced),
	):
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// Scrape collects data from Terraform API and sends it over channel as prometheus metric.
func (ScrapeOrganizations) Scrape(ctx context.Context, config *setup.Config, ch chan<- prometheus.Metric) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, name := range config.Organizations {
		name := name
		g.Go(func() error {
			return getOrganization(ctx, name, config, ch)
		})
	}

	return g.Wait()
}
