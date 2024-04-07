package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func yaml() {
	flags := []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "test"}),
		&cli.StringFlag{Name: "load"},
	}

	app := &cli.App{
		Action: func(cCtx *cli.Context) error {
			fmt.Println("--test value.*default: 0")
			return nil
		},
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load")),
		Flags:  flags,
	}

	app.Run(os.Args)
}

func general() {
	var language string
	var count int
	tasks := []string{"cook", "clean", "laundry", "eat", "sleep", "code"}

	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Name:                   "greet",
		Version:                "v19.99.0",
		Usage:                  "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lang",
				Aliases:     []string{"l"},
				EnvVars:     []string{"LEGACY_COMPAT_LANG", "APP_LANG", "LANG"},
				Value:       "english",
				Usage:       "language for the greeting",
				Required:    true,
				Destination: &language,
			},
			&cli.BoolFlag{
				Name:  "foo",
				Usage: "foo greeting",
				Count: &count,
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "password for the mysql database",
				FilePath: "/etc/mysql/password", // set default value from file
			},
			&cli.StringSliceFlag{
				Name:  "greeting",
				Usage: "Pass multiple greetings",
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Use a randomized port",
				Value:       0,
				DefaultText: "random",
				Action: func(ctx *cli.Context, v int) error {
					if v >= 65536 {
						return fmt.Errorf("Flag port value %v out of range[0-65535]", v)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:     "silent",
				Aliases:  []string{"s"},
				Usage:    "no messages",
				Category: "Miscellaneous:",
			},
			&cli.BoolFlag{
				Name:     "perl-regexp",
				Aliases:  []string{"P"},
				Usage:    "PATTERNS are Perl regular expressions",
				Category: "Pattern selection:",
			},
			&cli.BoolFlag{
				Name:  "ginger-crouton",
				Usage: "is it in the soup?",
			},
			&cli.TimestampFlag{
				Name:     "meeting",
				Layout:   "2006-01-02T15:04:05",
				Timezone: time.Local,
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.Bool("ginger-crouton") {
				return cli.Exit("Ginger croutons are not in the soup", 86)
			}

			fmt.Println("count", count)
			fmt.Println(strings.Join(cCtx.StringSlice("greeting"), `, `))
			//fmt.Printf("%s", cCtx.Timestamp("meeting").String())

			name := "someone"
			if cCtx.NArg() > 0 {
				name = cCtx.Args().Get(0)
			}
			if language == "spanish" {
				fmt.Println("Hola", name)
			} else {
				fmt.Println("Hello", name)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:     "complete",
				Aliases:  []string{"c"},
				Usage:    "complete a task on the list",
				Category: "Misc",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
				BashComplete: func(cCtx *cli.Context) {
					// This will complete if no args are passed
					if cCtx.NArg() > 0 {
						return
					}
					for _, t := range tasks {
						fmt.Println(t)
					}
				},
			},
			{
				Name:     "add",
				Aliases:  []string{"a"},
				Usage:    "add a task to the list",
				Category: "Misc",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("added task: ", cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "options for task templates",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a new template",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("new task template: ", cCtx.Args().First())
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove an existing template",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("removed task template: ", cCtx.Args().First())
							return nil
						},
					},
				},
			},
			{
				Name:  "short",
				Usage: "complete a task on the list",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "serve", Aliases: []string{"s"}},
					&cli.BoolFlag{Name: "option", Aliases: []string{"o"}},
					&cli.StringFlag{Name: "message", Aliases: []string{"m"}},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("serve:", cCtx.Bool("serve"))
					fmt.Println("option:", cCtx.Bool("option"))
					fmt.Println("message:", cCtx.String("message"))
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
WEBSITE: http://awesometown.example.com
SUPPORT: support@awesometown.example.com
`, cli.AppHelpTemplate)

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"v"},
		Usage:   "print only the version",
	}

	//yaml()
	general()
}
