package collector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicolaka/tfbi/internal/setup"

	tfe "github.com/hashicorp/go-tfe"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/smartystreets/goconvey/convey"
)

func TestScrapeWorkspaces(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"meta":{
				"pagination":{"current-page":1,"prev-page":null,"next-page":null,"total-pages":1,"total-count":2}
			},
			"data":[{
				"id":"test-id-1",
				"type":"workspaces",
				"attributes":{
					"name":"dev",
					"created-at":"1010-10-10T10:10:10.101Z",
					"environment":"test-environment",
					"terraform-version":"0.14.3",
					"latest-change-at":"2020-10-10T10:10:10.101Z"
				},
				"relationships":{
					"organization":{"data":{"id":"test-org","type":"organizations"}},
					"current-run":{
						"data":{
							"id":"run-id-1",
							"type":"runs",
							"attributes": {
								"created-at":"1010-10-10T10:10:10.101Z",
								"status": "applied"
							}
						}
					}
				}
			}, {
				"id":"test-id-2",
				"type":"workspaces",
				"attributes":{
					"name":"stg",
					"created-at":"1010-10-10T10:10:10.101Z",
					"environment":"test-environment",
					"terraform-version":"0.14.2",
					"latest-change-at":"2020-10-10T10:10:10.101Z"
				},
				"relationships":{
					"organization":{"data":{"id":"test-org","type":"organizations"}}
				}
			}]
		}`))
	}))
	defer mockAPI.Close()

	client, err := tfe.NewClient(&tfe.Config{
		Address: mockAPI.URL,
		Token:   "test",
	})
	if err != nil {
		t.Fatalf("error creating a stub api client: %s", err)
	}

	config := &setup.Config{
		Client: *client,
		CLI:    setup.CLI{Organizations: []string{"test-org"}},
	}

	ch := make(chan prometheus.Metric)
	go func() {
		defer close(ch)
		if err = (ScrapeWorkspaces{}).Scrape(context.Background(), config, ch); err != nil {
			t.Errorf("error calling function on test: %s", err)
		}
	}()

	counterExpected := []MetricResult{
		{labels: labelMap{"created_at": "1010-10-10 10:10:10.101 +0000 UTC", "current_run": "run-id-1", "current_run_status": "applied", "current_run_created_at": "1010-10-10 10:10:10.101 +0000 UTC", "environment": "test-environment", "id": "test-id-1", "name": "dev", "organization": "test-org", "terraform_version": "0.14.3"}, value: 1, metricType: dto.MetricType_GAUGE},
		{labels: labelMap{"created_at": "1010-10-10 10:10:10.101 +0000 UTC", "current_run": "na", "current_run_status": "na", "current_run_created_at": "na", "environment": "test-environment", "id": "test-id-2", "name": "stg", "organization": "test-org", "terraform_version": "0.14.2"}, value: 1, metricType: dto.MetricType_GAUGE},
	}
	convey.Convey("Metrics comparison", t, func() {
		for _, expect := range counterExpected {
			got := readMetric(<-ch)
			convey.So(got, convey.ShouldResemble, expect)
		}
	})
}
