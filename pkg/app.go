package andromeda

import (
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func makeapp() *cli.App {
	app := &cli.App{
		Name:    "andromeda",
		Usage:   "Andromeda the mother of all jobs",
		Version: strings.TrimSpace(string(Version)),
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Andromeda Server",
				Action: func(c *cli.Context) error {
					if !isatty.IsTerminal(os.Stdout.Fd()) {
						ansi.DisableColors(true)
					}
					return serve(c)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "public-url",
						Usage: "Public URL to show to user, useful when you are behind a proxy.",
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   defaultServerPort,
						Usage:   "Port to listen on",
					},
					&cli.StringFlag{
						Name:    "address",
						Aliases: []string{"a"},
						Value:   defaultServerAddress,
						Usage:   "Address to listen on",
					},
				},
			},
			{
				Name:  "client",
				Usage: "Andromeda Client",
				Action: func(c *cli.Context) error {
					logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
					var targetURL string
					if os.Getenv("ANDROmeda_TARGET_URL") != "" {
						targetURL = os.Getenv("ANDROmeda_TARGET_URL")
					} else {
						if c.NArg() != 2 {
							return fmt.Errorf("need at least a target-url")
						}
						targetURL = c.Args().Get(1)
					}
					if _, err := url.Parse(targetURL); err != nil {
						return fmt.Errorf("target url %s is not a valid url %w", targetURL, err)
					}
					cfg := andromeda{
						targetURL: targetURL,
						logger:    logger,
					}
					err := cfg.clientSetup()
					return err
				},
			},
		},
	}
	return app
}

func Run(args []string) error {
	return makeapp().Run(args)
}
