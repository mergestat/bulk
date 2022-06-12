package run

import (
	"context"
	"os/exec"
	"runtime"

	"golang.org/x/sync/errgroup"
)

func (r *run) clone(ctx context.Context) error {
	concurrency := r.config.Concurrency
	if concurrency == 0 {
		concurrency = runtime.NumCPU()
	}
	r.logger.Debug().Msgf("concurrency: %d", concurrency)

	repos := make(chan string, concurrency)
	var g errgroup.Group
	for i := 0; i < concurrency; i++ {
		g.Go(func() error {
			for repo := range repos {
				r.logger.Info().Msgf("cloning %s", repo)
				cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "2", repo)
				cmd.Dir = r.workDir

				if err := cmd.Run(); err != nil {
					r.logger.Err(err).Msg("run git command")
					continue
				}
			}

			return nil
		})
	}

	for _, repo := range r.repos {
		repos <- repo
	}
	close(repos)

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
