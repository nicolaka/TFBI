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

func TestScrapeOrganizations(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"id":"test-org",
				"type":"organizations",
				"attributes": {
					"external-id":"test-external-id",
					"created-at":"1010-10-10T10:10:10.101Z",
					"email":"test-email",
					"owners-team-saml-role-id":"test-role-id",
					"saml-enabled":true,
					"two-factor-conformant":false
				}
			}
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
		if err = (ScrapeOrganizations{}).Scrape(context.Background(), config, ch); err != nil {
			t.Errorf("error calling function on test: %s", err)
		}
	}()

	counterExpected := []MetricResult{
		{labels: labelMap{"created_at": "1010-10-10 10:10:10.101 +0000 UTC", "email": "test-email", "external_id": "test-external-id", "name": "test-org", "owners_team_saml_role_id": "test-role-id", "saml_enabled": "true", "two_factor_conformant": "false", "assessment_enforced": "false"}, value: 1, metricType: dto.MetricType_GAUGE},
	}
	convey.Convey("Metrics comparison", t, func() {
		for _, expect := range counterExpected {
			got := readMetric(<-ch)
			convey.So(got, convey.ShouldResemble, expect)
		}
	})
}
