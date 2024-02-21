package collector

import (
	"context"
	"strconv"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/nicolaka/tfbi/internal/setup"

	tfe "github.com/hashicorp/go-tfe"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// workspaces is the Metric subsystem we use.
	workspacesSubsystem = "workspaces"
	// Selecting a page size of 100 to minimize number of requests for larger  workspace deployments
	pageSize = 100
)

// Metric descriptors.
var (
	WorkspacesInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, workspacesSubsystem, "info"),
		"Information about existing workspaces",
		[]string{"id", "name", "organization", "terraform_version", "created_at", "environment", "current_run", "current_run_status", "current_run_created_at", "project", "assessments_enabled", "description", "resource_count", "policy_check_failures",  "run_failures", "runs_count" }, nil,
	)
)

// ScrapeWorkspaces scrapes metrics about the workspaces.
type ScrapeWorkspaces struct{}

func init() {
	Scrapers = append(Scrapers, ScrapeWorkspaces{})
}

// Name of the Scraper. Should be unique.
func (ScrapeWorkspaces) Name() string {
	return workspacesSubsystem
}

// Help describes the role of the Scraper.
func (ScrapeWorkspaces) Help() string {
	return "Scrape information from the Workspaces API: https://www.terraform.io/docs/cloud/api/workspaces.html"
}

// Version of Terraform Cloud/Enterprise API from which scraper is available.
func (ScrapeWorkspaces) Version() string {
	return "v2"
}

func getWorkspacesListPage(ctx context.Context, page int, organization string, config *setup.Config, ch chan<- prometheus.Metric) error {
	workspacesList, err := config.Client.Workspaces.List(ctx, organization, &tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageSize:   pageSize,
			PageNumber: page,
		},
		Include: []tfe.WSIncludeOpt{
			// go-tfe bug 764 "project",
			"current_run",
			"organization",
		},
	})
	if err != nil {
		return fmt.Errorf("%v, (organization=%s, page=%d)", err, organization, page)
	}

	for _, w := range workspacesList.Items {
		select {
		case ch <- prometheus.MustNewConstMetric(
			WorkspacesInfo,
			prometheus.GaugeValue,
			1,
			w.ID,
			w.Name,
			w.Organization.Name,
			w.TerraformVersion,
			w.CreatedAt.String(),
			w.Environment,
			getCurrentRunID(w.CurrentRun),
			getCurrentRunStatus(w.CurrentRun),
			getCurrentRunCreatedAt(w.CurrentRun),
			w.Project.ID,
			strconv.FormatBool(w.AssessmentsEnabled),
			w.Description,
			strconv.Itoa(w.ResourceCount),
			strconv.Itoa(w.PolicyCheckFailures),
			strconv.Itoa(w.RunFailures),
			strconv.Itoa(w.RunsCount),
		):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// Scrape collects data from Terraform API and sends it over channel as prometheus metric.
func (ScrapeWorkspaces) Scrape(ctx context.Context, config *setup.Config, ch chan<- prometheus.Metric) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, name := range config.Organizations {
		name := name
		g.Go(func() error {
			workspacesList, err := config.Client.Workspaces.List(ctx, name,&tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageSize:   pageSize,
				},})
			
			if err != nil {
				return fmt.Errorf("%v, organization=%s", err, name)
			}
			
			for i := 1; i <= workspacesList.Pagination.TotalPages; i++ {
				if err := getWorkspacesListPage(ctx, i, name, config, ch); err != nil {
					return err
				}
			}

			return nil
		})
	}

	return g.Wait()
}

func getCurrentRunID(r *tfe.Run) string {
	if r == nil {
		return "na"
	}

	return r.ID
}

func getCurrentRunStatus(r *tfe.Run) string {
	if r == nil {
		return "na"
	}

	return string(r.Status)
}

func getCurrentRunCreatedAt(r *tfe.Run) string {
	if r == nil {
		return "na"
	}

	return r.CreatedAt.String()
}
