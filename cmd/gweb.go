package main

import (
	"os"
	"time"

	"github.com/imssyang/gweb/internal/cmd"
	"github.com/imssyang/gweb/internal/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "Gweb",
		Version:  "v1.0.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Quincy Yang",
				Email: "imssyang@gmail.com",
			},
		},
		Usage:  "A painless self-hosted tool service",
		Flags:  cmd.Flags,
		Action: cmd.Action,
	}

	if err := app.Run(os.Args); err != nil {
		log.Zap.Fatal(err)
	}
}
