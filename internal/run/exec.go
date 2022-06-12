package run

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
)

func (r *run) exec(ctx context.Context) error {
	concurrency := r.config.Concurrency
	if concurrency == 0 {
		concurrency = runtime.NumCPU()
	}

	var mut sync.Mutex

	repos := make(chan string, concurrency)
	var g errgroup.Group
	for i := 0; i < concurrency; i++ {
		g.Go(func() error {
			for repo := range repos {
				func() {
					mut.Lock()
					r.logger.Info().Msgf("running actions on: %s", repo)
					defer mut.Unlock()

					for _, action := range r.config.Actions {
						cmd := exec.CommandContext(ctx, "bash", "-c", action.Command)
						cmd.Dir = path.Join(r.workDir, repo)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr

						if err := cmd.Run(); err != nil {
							r.logger.Err(err).Msg("running action")
							continue
						}
					}
				}()
			}

			return nil
		})
	}

	repoDirs, err := ioutil.ReadDir(r.workDir)
	if err != nil {
		return err
	}

	for _, d := range repoDirs {
		if d.IsDir() {
			repos <- d.Name()
		}
	}
	close(repos)

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
