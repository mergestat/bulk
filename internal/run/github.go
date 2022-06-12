package run

import (
	"context"
	"os"

	"github.com/mergestat/bulk/internal/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// githubClient retrieves a GitHub API client, using the cached one if it exists
func (r *run) githubClient() *githubv4.Client {
	if r.cachedGitHubClient == nil {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		httpClient := oauth2.NewClient(context.Background(), src)
		r.cachedGitHubClient = githubv4.NewClient(httpClient)
	}

	return r.cachedGitHubClient
}

type githubReposGraphQL struct {
	Nodes []struct {
		Name        githubv4.String
		Description githubv4.String
		URL         githubv4.String `graphql:"url"`
		Languages   struct {
			Nodes []struct {
				Name githubv4.String
			}
		} `graphql:"languages(first: 100)"`
		RepositoryTopics struct {
			Nodes []struct {
				Topic struct {
					Name githubv4.String
				}
			}
		} `graphql:"repositoryTopics(first: 100)"`
	}
	PageInfo struct {
		EndCursor   githubv4.String
		HasNextPage githubv4.Boolean
	}
}

func (r *run) materializeGitHubOrgRepos(ctx context.Context, orgRepos *config.GitHubOrgReposSpec) ([]string, error) {
	if orgRepos == nil {
		return []string{}, nil
	}

	client := r.githubClient()

	var query struct {
		Organizations struct {
			Repositories *githubReposGraphQL `graphql:"repositories(first: $perPage, after: $pageCursor)"`
		} `graphql:"organization(login: $login)"`
	}

	var cursor *githubv4.String
	var repos []string

	for {
		if err := client.Query(ctx, &query, map[string]interface{}{
			"login":      githubv4.String(orgRepos.OrgName),
			"perPage":    githubv4.Int(10),
			"pageCursor": cursor,
		}); err != nil {
			return nil, err
		}

		cursor = &query.Organizations.Repositories.PageInfo.EndCursor

		if !query.Organizations.Repositories.PageInfo.HasNextPage {
			break
		}

		for _, repo := range query.Organizations.Repositories.Nodes {
			repos = append(repos, string(repo.URL))
		}
	}

	return repos, nil
}

func (r *run) materializeGitHubUserRepos(ctx context.Context, userRepos *config.GithubUserReposSpec) ([]string, error) {
	if userRepos == nil {
		return []string{}, nil
	}

	client := r.githubClient()

	var query struct {
		User struct {
			Repositories *githubReposGraphQL `graphql:"repositories(first: $perPage, after: $pageCursor)"`
		} `graphql:"user(login: $login)"`
	}

	var cursor *githubv4.String
	var repos []string

	for {
		if err := client.Query(ctx, &query, map[string]interface{}{
			"login":      githubv4.String(userRepos.UserLogin),
			"perPage":    githubv4.Int(10),
			"pageCursor": cursor,
		}); err != nil {
			return nil, err
		}

		cursor = &query.User.Repositories.PageInfo.EndCursor

		if !query.User.Repositories.PageInfo.HasNextPage {
			break
		}

		for _, repo := range query.User.Repositories.Nodes {
			repos = append(repos, string(repo.URL))
		}
	}

	return repos, nil
}
