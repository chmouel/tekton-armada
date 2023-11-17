package armada

import (
	"context"
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
		Name:    "armada",
		Usage:   "Armada the mother of all jobs",
		Version: strings.TrimSpace(string(Version)),
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Armada Server",
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
				Usage: "Armada Client",
				Action: func(c *cli.Context) error {
					logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
					var targetURL string
					if os.Getenv("ARMADA_TARGET_URL") != "" {
						targetURL = os.Getenv("ARMADA_TARGET_URL")
					} else {
						if c.NArg() != 1 {
							return fmt.Errorf("need at least a target-url")
						}
						targetURL = c.Args().First()
					}
					if _, err := url.Parse(targetURL); err != nil {
						return fmt.Errorf("target url %s is not a valid url %w", targetURL, err)
					}
					armada, err := NewArmada(logger, targetURL)
					if err != nil {
						return err
					}
					ctx := context.Background()
					return armada.clientSetup(ctx)
				},
			},
		},
	}
	return app
}

func Run(args []string) error {
	return makeapp().Run(args)
}
