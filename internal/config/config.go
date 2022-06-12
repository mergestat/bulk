package config

type Config struct {
	Concurrency int           `yaml:"concurrency"`
	WorkDir     string        `yaml:"directory"`
	Repos       *RepoSpec     `yaml:"repos"`
	Actions     []*ActionSpec `yaml:"actions"`
}

type RepoSpec struct {
	Inline          []*InlineRepoSpec      `yaml:"inline"`
	GitHubOrgRepos  []*GitHubOrgReposSpec  `yaml:"githubOrgRepos"`
	GitHubUserRepos []*GithubUserReposSpec `yaml:"githubUserRepos"`
}

type InlineRepoSpec struct {
	Repo       string   `yaml:"repo"`
	Ref        *string  `yaml:"ref"`
	CloneFlags []string `yaml:"cloneFlags"`
}

type GitHubOrgReposSpec struct {
	OrgName       string `yaml:"org"`
	IncludeFilter string `yaml:"include"`
	ExcludeFilter string `yaml:"exlcude"`
}

type GithubUserReposSpec struct {
	UserLogin     string `yaml:"user"`
	IncludeFilter string `yaml:"include"`
	ExcludeFilter string `yaml:"exlcude"`
}

type ActionSpec struct {
	Command string `yaml:"command"`
}
