package main

import (
	"fmt"
	"os"

	"github.com/sladenkiy77/worker-service/cmd"
	msentry "github.com/sladenkiy77/worker-service/internal/app/monitor/sentry"
	"github.com/sladenkiy77/worker-service/internal/pkg/constant"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    constant.AppName,
		Version: constant.Version,
		Usage:   "Worker management service",
		Commands: []*cli.Command{
			cmd.WebServer(),
			cmd.SubscribeOrderCreated(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-json",
				Usage: "Human-readable log format",
			},
		},
	}

	defer msentry.Flush()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
