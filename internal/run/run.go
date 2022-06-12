package run

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/mergestat/bulk/internal/config"
	"github.com/rs/zerolog"
	"github.com/shurcooL/githubv4"
)

type run struct {
	config             *config.Config
	logger             zerolog.Logger
	cachedGitHubClient *githubv4.Client
	repos              []string
	workDir            string
}

type Option func(*run)

func WithLogger(logger zerolog.Logger) Option {
	return func(r *run) {
		r.logger = logger
	}
}

// New returns a new run instance, which represents an execution of bulk actions against a set of repos
func New(config *config.Config, options ...Option) *run {
	r := &run{
		config:             config,
		logger:             zerolog.Nop(),
		cachedGitHubClient: nil,
		repos:              make([]string, 0),
	}

	for _, opt := range options {
		opt(r)
	}

	return r
}

// materializeRepos populates the `repos` slice on the run
func (r *run) materializeRepos(ctx context.Context) error {
	var repos []string

	for _, orgRepo := range r.config.Repos.GitHubOrgRepos {
		if r, err := r.materializeGitHubOrgRepos(ctx, orgRepo); err != nil {
			return err
		} else {
			repos = append(repos, r...)
		}
	}

	for _, orgRepo := range r.config.Repos.GitHubUserRepos {
		if r, err := r.materializeGitHubUserRepos(ctx, orgRepo); err != nil {
			return err
		} else {
			repos = append(repos, r...)
		}
	}

	// de-duplicate repos
	repoMap := make(map[string]struct{})
	for _, inlineRepo := range r.config.Repos.Inline {
		repos = append(repos, inlineRepo.Repo)
	}

	for _, repo := range repos {
		repoMap[repo] = struct{}{}
	}

	repos = make([]string, 0, len(repoMap))
	for r := range repoMap {
		repos = append(repos, r)
	}

	r.repos = repos

	return nil
}

// Exec executes the run
func (r *run) Exec(ctx context.Context) error {
	if err := r.materializeRepos(ctx); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := r.config.WorkDir
	if dir == "" {
		dir = ".bulk"
	}

	p := path.Join(cwd, dir)

	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		return err
	}

	if r.workDir, err = ioutil.TempDir(p, ""); err != nil {
		return err
	}

	if err := r.clone(ctx); err != nil {
		return err
	}

	if err := r.exec(ctx); err != nil {
		return err
	}

	return nil
}
