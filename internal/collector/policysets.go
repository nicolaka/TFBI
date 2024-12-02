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
	// policysets is the Metric subsystem we use.
	policysetsSubsystem = "policysets"
)

// Metric descriptors.
var (
	PolicySetsInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, policysetsSubsystem, "info"),
		"Information about existing policysets",
		[]string{"id", "name", "description", "kind", "global", "policy_count", "workspace_count", "project_count", "created_at", "updated_at", "organization"}, nil,
	)
)

// ScrapePolicySets scrapes metrics about the policysets.
type ScrapePolicySets struct{}

func init() {
	Scrapers = append(Scrapers, ScrapePolicySets{})
}

// Name of the Scraper. Should be unique.
func (ScrapePolicySets) Name() string {
	return policysetsSubsystem
}

// Help describes the role of the Scraper.
func (ScrapePolicySets) Help() string {
	return "Scrape information from the PolicySets API: https://www.terraform.io/docs/cloud/api/policysets.html"
}

// Version of Terraform Cloud/Enterprise API from which scraper is available.
func (ScrapePolicySets) Version() string {
	return "v2"
}

// []string{"id", "name","description","kind","global","policy_count","workspace_count","project_count","created_at","updated_at"}, nil,

func (ScrapePolicySets) Scrape(ctx context.Context, config *setup.Config, ch chan<- prometheus.Metric) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, name := range config.Organizations {
		g.Go(func() error {
			policysetList, err := config.Client.PolicySets.List(ctx, name, nil)
			for _, p := range policysetList.Items {

				if err != nil {
					return fmt.Errorf("%v, organization=%s", err, name)
				}

				select {
				case ch <- prometheus.MustNewConstMetric(
					PolicySetsInfo,
					prometheus.GaugeValue,
					1,
					p.ID,
					p.Name,
					p.Description,
					string(p.Kind),
					strconv.FormatBool(p.Global),
					strconv.Itoa(p.PolicyCount),
					strconv.Itoa(p.WorkspaceCount),
					strconv.Itoa(p.ProjectCount),
					p.CreatedAt.String(),
					p.UpdatedAt.String(),
					p.Organization.Name,
				):
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			return nil
		})

	}
	return g.Wait()
}
