# Adopting Go with Trusted GitHub API Clients

* Status: **accepted** 👍
* Deciders: [@hackebrot] [@dlactin]
* Date: 2025-03-31

Technical Story: This decision addresses issues discussed in Jira tickets
[CICD-45] and [CICD-439].

[CICD-45]: https://mozilla-hub.atlassian.net/browse/CICD-45
[CICD-439]: https://mozilla-hub.atlassian.net/browse/CICD-439

[@hackebrot]: https://github.com/hackebrot
[@dlactin]: https://github.com/dlactin

## Context and Problem Statement

The Four Keys BigQuery dataset currently lacks historical data because it only
ingests events (such as Git commit pushes, deployments, and incidents) as they
occur through GitHub webhooks. This real-time approach leaves a data gap that
limits the completeness of our retroactive deployment analysis.

Furthermore, the current method for measuring `Lead Time for Changes` links
deployments to pushes using a single commit SHA. This approach fails to capture
all relevant commits—particularly those from earlier pushes that were not
deployed—and assumes an immutable git history, which is not always the case. By
leveraging the GitHub APIs for retroactive data collection, we can fill this
historical data gap and capture all relevant commit information, thereby
enhancing the accuracy of our deployment-related metrics.

## Decision Drivers

1. **Retroactive Data Collection Capability:**
   The solution must facilitate the retrieval of historical data via the GitHub
   APIs, effectively addressing the gaps in the Four Keys dataset.
2. **Stability and Maturity:**
   The solution must leverage well-established libraries that have proven
   reliability in production environments.
3. **Maintenance and Support:**
   The selected tools must be actively maintained and well-documented, with
   strong community or official support to ensure long-term viability.
4. **Team Expertise and Ecosystem Fit:**
   It should align with the team's existing skills and integrate smoothly with
   our current infrastructure and technology stack (e.g., Kubernetes, Argo CD).

## Considered Options

* A. Python with packages like `PyGitHub` and `gql`
* B. Go with packages like `github.com/shurcooL/githubv4` and `github.com/google/go-github/v68/github`
* C. Custom-Built API Wrappers

## Decision Outcome

Chosen option:

**B. Go with packages like `github.com/shurcooL/githubv4` and `github.com/google/go-github/v68/github`**

Option B was selected because it best satisfies our key decision drivers. Go’s
strong compile-time type safety, combined with the maturity and active
maintenance of the selected packages, ensures production-ready support for the
GitHub GraphQL API. Additionally, Go aligns well with our team's expertise and
technology stack (e.g., Kubernetes, Argo CD).

While Python offers advantages in rapid prototyping, its available GitHub API
libraries lack the specialized stability and GitHub-specific support required
for this project.

Furthermore, developing custom API wrappers was ruled out due to the
significantly higher development and maintenance overhead, and the risk of
re-implementing functionality already well-supported by mature libraries.

## Links

* Go client library for accessing the GitHub GraphQL API v4: https://github.com/shurcooL/githubv4
* Go client library for accessing the GitHub REST API v3: https://github.com/google/go-github
* Python client library for accessing the GitHub REST API v3: https://pypi.org/project/PyGithub/
* General purpose Python GraphQL client: https://pypi.org/project/gql/
