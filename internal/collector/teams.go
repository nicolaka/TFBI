package collector

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/hashicorp/go-tfe"
	"github.com/nicolaka/tfbi/internal/setup"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// teams is the Metric subsystem we use.
	teamsSubsystem = "teams"
)

// Metric descriptors.
var (
	TeamsInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, teamsSubsystem, "info"),
		"Information about existing teams",
		[]string{"id", "name", "sso_team_id", "users_count"}, nil,
	)
)

// ScrapeTeams scrapes metrics about the teams.
type ScrapeTeams struct{}

func init() {
	Scrapers = append(Scrapers, ScrapeTeams{})
}

// Name of the Scraper. Should be unique.
func (ScrapeTeams) Name() string {
	return teamsSubsystem
}

// Help describes the role of the Scraper.
func (ScrapeTeams) Help() string {
	return "Scrape information from the Teams API: https://www.terraform.io/docs/cloud/api/teams.html"
}

// Version of Terraform Cloud/Enterprise API from which scraper is available.
func (ScrapeTeams) Version() string {
	return "v2"
}

func getTeamsListPage(ctx context.Context, page int, organization string, config *setup.Config, ch chan<- prometheus.Metric) error {
	teamsList, err := config.Client.Teams.List(ctx, organization, &tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{
			PageSize:   pageSize,
			PageNumber: page,
		},
		Include: []tfe.TeamIncludeOpt{
			"organization-memberships",
		},
	})

	if err != nil {
		return fmt.Errorf("%v, (organization=%s, page=%d)", err, organization, page)
	}

	for _, t := range teamsList.Items {
		select {
		case ch <- prometheus.MustNewConstMetric(
			TeamsInfo,
			prometheus.GaugeValue,
			1,
			t.ID,
			t.Name,
			t.SSOTeamID,
			strconv.Itoa(t.UserCount),
		):
		case <-ctx.Done():
			return ctx.Err()
		}

	}

	return nil
}

func (ScrapeTeams) Scrape(ctx context.Context, config *setup.Config, ch chan<- prometheus.Metric) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, name := range config.Organizations {
		name := name
		g.Go(func() error {
			teamsList, err := config.Client.Teams.List(ctx, name, &tfe.TeamListOptions{
				ListOptions: tfe.ListOptions{
					PageSize: pageSize,
				}})

			if err != nil {
				return fmt.Errorf("%v, organization=%s", err, name)
			}

			for i := 1; i <= teamsList.Pagination.TotalPages; i++ {
				if err := getTeamsListPage(ctx, i, name, config, ch); err != nil {
					return err
				}
			}

			return nil
		})
	}

	return g.Wait()
}
