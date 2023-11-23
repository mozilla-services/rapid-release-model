package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/mozilla-services/rapid-release-model/metrics/internal/github"
)

// Interface for CSV, JSON and other encoders
type Encoder interface {
	Encode(w io.Writer, v interface{}) error
}

type JSONEcoder struct{}

func (j *JSONEcoder) Encode(w io.Writer, v interface{}) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	return e.Encode(v)
}

func NewJSONEncoder() (*JSONEcoder, error) {
	return &JSONEcoder{}, nil
}

type PlainEncoder struct{}

func (p *PlainEncoder) Encode(w io.Writer, v interface{}) error {
	_, err := fmt.Fprint(w, v)
	return err
}

func NewPlainEncoder() (*PlainEncoder, error) {
	return &PlainEncoder{}, nil
}

type CSVEncoder struct{}

func (c *CSVEncoder) Encode(w io.Writer, v interface{}) error {
	var records [][]string

	csvw := csv.NewWriter(w)

	switch v := v.(type) {
	case []github.PullRequest:
		records = PullRequestsToCSVRecords(v)
	case []github.Release:
		records = ReleasesToCSVRecords(v)
	case []github.Deployment:
		records = DeploymentsToCSVRecords(v)
	default:
		return fmt.Errorf("unable to export type %T to CSV", v)
	}

	return csvw.WriteAll(records)
}

func NewCSVEncoder() (*CSVEncoder, error) {
	return &CSVEncoder{}, nil
}

func PullRequestsToCSVRecords(prs []github.PullRequest) [][]string {
	var records [][]string

	// Add column headers to records
	records = append(records, []string{
		"id",
		"number",
		"title",
		"createdAt",
		"updatedAt",
		"closedAt",
		"mergedAt",
	})

	// Add a record for each pull request
	for _, pr := range prs {
		record := []string{
			pr.ID,
			strconv.Itoa(pr.Number),
			pr.Title,
			pr.CreatedAt.Format(time.RFC3339),
			pr.UpdatedAt.Format(time.RFC3339),
			pr.ClosedAt.Format(time.RFC3339),
			pr.MergedAt.Format(time.RFC3339),
		}
		records = append(records, record)
	}
	return records
}

func ReleasesToCSVRecords(rs []github.Release) [][]string {
	var records [][]string

	// Add column headers to records
	records = append(records, []string{
		"name",
		"tagName",
		"isDraft",
		"isLatest",
		"isPrerelease",
		"description",
		"createdAt",
		"publishedAt",
	})

	// Add a record for each release
	for _, r := range rs {
		record := []string{
			r.Name,
			r.TagName,
			strconv.FormatBool(r.IsDraft),
			strconv.FormatBool(r.IsLatest),
			strconv.FormatBool(r.IsPrerelease),
			r.Description,
			r.CreatedAt.Format(time.RFC3339),
			r.PublishedAt.Format(time.RFC3339),
		}
		records = append(records, record)
	}
	return records
}

func DeploymentsToCSVRecords(ds []github.Deployment) [][]string {
	var records [][]string

	// Add column headers to records
	records = append(records, []string{
		"description",
		"createdAt",
		"updatedAt",
		"originalEnvironment",
		"latestEnvironment",
		"task",
		"state",
		"commitOid",
	})

	// Add a record for each deployment
	for _, d := range ds {
		record := []string{
			d.Description,
			d.CreatedAt.Format(time.RFC3339),
			d.UpdatedAt.Format(time.RFC3339),
			d.OriginalEnvironment,
			d.LatestEnvironment,
			d.Task,
			d.State,
			d.Commit.AbbreviatedOid,
		}
		records = append(records, record)
	}
	return records
}
