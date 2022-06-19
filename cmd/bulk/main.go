package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/mergestat/bulk/internal/config"
	"github.com/mergestat/bulk/internal/run"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

var (
	debug bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "include debug output")
	flag.Parse()
}

func main() {
	configFilePath := "bulk.yaml"

	if len(flag.Args()) > 1 {
		configFilePath = flag.Arg(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configContents, err := ioutil.ReadFile(path.Join(cwd, configFilePath))
	if err != nil {
		log.Fatal(err)
	}

	var c config.Config
	err = yaml.Unmarshal(configContents, &c)
	if err != nil {
		log.Fatal(err)
	}

	level := zerolog.InfoLevel

	if debug {
		level = zerolog.DebugLevel
	}

	l := zerolog.New(os.Stderr).
		Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp}).
		Level(level).
		With().
		Timestamp().Logger()

	r := run.New(&c, run.WithLogger(l))

	if err := r.Exec(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
