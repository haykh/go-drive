package main

import (
	"context"
	"go-drive/api"
	"go-drive/browser"
	"go-drive/ui"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"

	"github.com/urfave/cli/v3"
)

func initLogger(debug bool) {
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
					initLogger(debug_mode)
					_, err := api.GetGoogleDriveService(ctx, c.String("credentials"), c.String("token"), drive.DriveScope, true)
					return api.ToHumanReadableError(err, debug_mode)
				},
			},
			{
				Name:  "ls",
				Usage: "list content of a Google Drive folder",
				Flags: common_flags,
				Action: func(ctx context.Context, c *cli.Command) error {
					debug_mode := c.Bool("debug")
					initLogger(debug_mode)
					dir := ""
					if c.NArg() > 0 {
						dir = c.Args().Get(0)
					}
					if srv, err := api.GetGoogleDriveService(ctx, c.String("credentials"), c.String("token"), drive.DriveScope, false); err != nil {
						return api.ToHumanReadableError(err, debug_mode)
					} else {
						if content, err := ui.RunWithSpinner(
							func() (any, error) {
								if ret, err := api.GetFolderContent(srv, dir); err != nil {
									return nil, api.ToHumanReadableError(err, debug_mode)
								} else {
									return ret, nil
								}
							},
							"loading",
							"unable to get directory content",
							"",
							debug_mode,
						); err != nil {
							return err
						} else {
							for _, str_item := range ui.StringizeItemList(content.(*drive.FileList), "%icon% %shared% %name% %size%") {
								log.Print(str_item)
							}
							return nil
						}
					}
				},
			},
			{
				Name:  "fs",
				Usage: "open file picker in the Google Drive folder",
				Flags: common_flags,
				Action: func(ctx context.Context, c *cli.Command) error {
					debug_mode := c.Bool("debug")
					initLogger(debug_mode)
					dir := ""
					if c.NArg() > 0 {
						dir = c.Args().Get(0)
					}
					if srv, err := api.GetGoogleDriveService(ctx, c.String("credentials"), c.String("token"), drive.DriveScope, false); err != nil {
						return api.ToHumanReadableError(err, debug_mode)
					} else {
						// @TODO
						// fmt.Printf("%v %s", srv, dir)
						// return nil
						return browser.FileBrowser(srv, dir, debug_mode)
					}
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// func main() {
// 	ctx := context.Background()
//
// 	scope := drive.DriveScope
// 	if client, err := api.GetGoogleClient(ctx, "./credentials/credentials.json", scope); err != nil {
// 		panic(err)
// 	} else {
// 		if srv, err := drive.NewService(ctx, option.WithHTTPClient(client)); err != nil {
// 			log.Fatalf("Unable to retrieve Drive client: %+v", err)
// 		} else {
// 			if newfile, err := api.UploadFile(srv, "./hello.txt", "Test", api.Overwrite); err != nil {
// 				panic(err)
// 			} else {
// 				fmt.Printf("%v\n", newfile)
// 			}
// 			// if folder_id, err := api.GetFolderId(srv, "Literature"); err != nil {
// 			// 	panic(err)
// 			// } else {
// 			// 	fmt.Printf("Folder id: %s\n", folder_id)
// 			// }
// 			// if content, err := api.GetFolderContent(srv, "Literature"); err != nil {
// 			// 	panic(err)
// 			// } else {
// 			// 	for _, c := range content {
// 			// 		fmt.Printf("%v\n", c)
// 			// 	}
// 			// }
// 		}
// 	}
//
// 	// r, err := srv.Files.List().
// 	// 	PageSize(10).
// 	// 	Fields("nextPageToken, files(id, name)").
// 	// 	Do()
// 	//
// 	// if err != nil {
// 	// 	log.Fatalf("Unable to retrieve files: %+v", err)
// 	// }
// 	// fmt.Println("Files:")
// 	// if len(r.Files) == 0 {
// 	// 	fmt.Println("No files found.")
// 	// } else {
// 	// 	for _, i := range r.Files {
// 	// 		fmt.Printf("%s (%s)\n", i.Name, i.Id)
// 	// 	}
// 	// }
// }
