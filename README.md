# bulk
 
`bulk` is a CLI for running bulk actions (ad-hoc scripts) on sets of Git repositories.
Often, it can be useful to run some action(s) against a grouping of repositories (all the repos in an org, all repos of a certain language, all rails repos, all golang repos etc.)

## Examples

The `GITHUB_TOKEN` env variable must be set to fetch repositories from a GitHub org or user.

### Check if repo has a `LICENSE` file

```yaml
# bulk.yaml
repos:
  githubOrgRepos:
    - org: mergestat

actions:
  # - command: cat LICENSE
  - command: >-
      if [[ -f "LICENSE" ]]; then echo ✅; else echo ❌;fi
```



## Roadmap
- [ ] File based output mode - create an output file for the results of commands run in each repo
- [ ] Logging improvements
- [ ] Support repo caching, so that a new clone is not needed on every run
- [ ] Release with GoReleaser
