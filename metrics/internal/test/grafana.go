package test

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-cmp/cmp"
)

// GrafanaReqParams holds expected values for outgoing requests to Grafana.
type GrafanaReqParams struct {
	Path   string
	Params url.Values
}

// FakeGrafanaClient returns canned responses (fixtures) rather than
// sending queries to the live Grafana REST API.
type FakeGrafanaClient struct {
	reqParams *GrafanaReqParams
}

func (c *FakeGrafanaClient) Get(ctx context.Context, p string, params url.Values) ([]byte, error) {
	if c.reqParams != nil && !cmp.Equal(p, c.reqParams.Path) {
		return nil, fmt.Errorf("unexpected path for HTTP query\n%v", cmp.Diff(p, c.reqParams.Path))
	}

	if c.reqParams != nil && !cmp.Equal(params, c.reqParams.Params) {
		return nil, fmt.Errorf("unexpected URL values for HTTP query\n%v", cmp.Diff(params, c.reqParams.Params))
	}

	return LoadFixture("grafana", p, "response.json")
}
