package main

import (
	"bytes"
	"context"
	"go-drive/api"
	"go-drive/utils"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"

	"github.com/urfave/cli/v3"
)

func initLogger(debug bool, bufferize bool) *bytes.Buffer {
	var logBuffer bytes.Buffer
	var logger *log.Logger
	if debug {
		logger = log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller: true,
			Level:        log.DebugLevel,
		})
	} else {
		logger = log.NewWithOptions(os.Stderr, log.Options{
			Level: log.InfoLevel,
		})
	}
	log.SetDefault(logger)
	if bufferize {
		log.SetOutput(&logBuffer)
		return &logBuffer
	}
	return nil
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get user home directory: %v", err)
	}

	common_flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "credentials",
			Aliases: []string{"c"},
			Value:   filepath.Join(home, ".config", "godrive", "credentials.json"),
			Usage:   "path to the Google API credentials file",
		},
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t"},
			Value:   filepath.Join(home, ".config", "godrive", "token.json"),
			Usage:   "path to the Google Drive token file",
		},
		&cli.StringFlag{
			Name:    "local",
			Aliases: []string{"l"},
			Value:   filepath.Join(home, ".config", "godrive", "storage"),
			Usage:   "path to the local mirror storage",
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "enable debug mode",
		},
	}

	cmd := &cli.Command{
		Name:    "godrive",
		Version: "1.0.0",
		Usage:   "access Google Drive from the command line",
		Commands: []*cli.Command{
			{
				Name:    "auth",
				Aliases: []string{"a"},
				Usage:   "authenticate with Google Drive and save the token",
				Flags:   common_flags,
				Action: func(ctx context.Context, c *cli.Command) error {
					debug_mode := c.Bool("debug")
					initLogger(debug_mode, false)
					_, err := api.GetGoogleDriveService(ctx, c.String("credentials"), c.String("token"), drive.DriveScope, true)
					return utils.ToHumanReadableError(err, debug_mode)
				},
			},
			{
				Name:  "ls",
				Usage: "list content of the remote directory",
				Flags: common_flags,
				Action: func(ctx context.Context, c *cli.Command) error {
					debug_mode := c.Bool("debug")
					initLogger(debug_mode, false)
					dir := ""
					if c.NArg() > 0 {
						dir = c.Args().Get(0)
					}
					if srv, err := api.GetGoogleDriveService(ctx, c.String("credentials"), c.String("token"), drive.DriveScope, false); err != nil {
						return utils.ToHumanReadableError(err, debug_mode)
					} else {
						return api.DualLs(srv, c.String("local"), dir, debug_mode)
					}
				},
			},
			{
				Name:  "cd",
				Usage: "open file picker with a dual view (remote and local)",
				Flags: common_flags,
				Action: func(ctx context.Context, c *cli.Command) error {
					debug_mode := c.Bool("debug")
					debugBuffer := initLogger(debug_mode, true)
					dir := ""
					if c.NArg() > 0 {
						dir = c.Args().Get(0)
					}
					if srv, err := api.GetGoogleDriveService(ctx, c.String("credentials"), c.String("token"), drive.DriveScope, false); err != nil {
						return utils.ToHumanReadableError(err, debug_mode)
					} else {
						return api.DualFileBrowser(srv, c.String("local"), dir, debug_mode, debugBuffer)
					}
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
