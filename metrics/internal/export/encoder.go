package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mozilla-services/rapid-release-model/pkg/github"
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
	var err error

	csvw := csv.NewWriter(w)

	switch v := v.(type) {
	case []github.PullRequest:
		records = PullRequestsToCSVRecords(v)
	case []github.Release:
		records = ReleasesToCSVRecords(v)
	case []github.ReleaseWithPRs:
		records, err = ReleasesWithPRsToCSVRecords(v)
		if err != nil {
			return err
		}
	case []github.Deployment:
		records = DeploymentsToCSVRecords(v)
	case map[string][]*github.DeploymentWithCommits:
		records = DeploymentsWithCommitsToCSVRecords(v)
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

// ReleasesWithPRsToCSVRecords is identical to ReleasesToCSVRecords with the
// addition of an extra column on the right for a JSON array of PR numbers.
func ReleasesWithPRsToCSVRecords(rs []github.ReleaseWithPRs) ([][]string, error) {
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
		"prs",
	})

	// Add a record for each release
	for _, r := range rs {
		prs, err := json.Marshal(r.PRs)
		if err != nil {
			return nil, fmt.Errorf("error encoding PRs as JSON: %w", err)
		}
		record := []string{
			r.Release.Name,
			r.Release.TagName,
			strconv.FormatBool(r.Release.IsDraft),
			strconv.FormatBool(r.Release.IsLatest),
			strconv.FormatBool(r.Release.IsPrerelease),
			r.Release.Description,
			r.Release.CreatedAt.Format(time.RFC3339),
			r.Release.PublishedAt.Format(time.RFC3339),
			string(prs),
		}
		records = append(records, record)
	}

	return records, nil
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
		"abbreviatedCommitSHA",
		"commitSHA",
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
			string(d.State),
			string(d.Commit.AbbreviatedSHA),
			string(d.Commit.SHA),
		}
		records = append(records, record)
	}
	return records
}

func DeploymentsWithCommitsToCSVRecords(dByEnv map[string][]*github.DeploymentWithCommits) [][]string {
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
		"abbreviatedCommitSHA",
		"commitSHA",
		"deployedCommits",
	})

	// Add a record for each deployment
	for _, deploys := range dByEnv {
		for _, d := range deploys {

			var deployedSHAs []string
			for _, commit := range d.DeployedCommits {
				deployedSHAs = append(deployedSHAs, commit.SHA)
			}

			record := []string{
				d.Description,
				d.CreatedAt.Format(time.RFC3339),
				d.UpdatedAt.Format(time.RFC3339),
				d.OriginalEnvironment,
				d.LatestEnvironment,
				d.Task,
				string(d.State),
				string(d.Commit.AbbreviatedSHA),
				string(d.Commit.SHA),
				strings.Join(deployedSHAs, ","),
			}
			records = append(records, record)
		}
	}
	return records
}
