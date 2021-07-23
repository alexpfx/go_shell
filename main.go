package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/alexpfx/go_shell/internal/goshell"
	"github.com/urfave/cli/v2"
)

const defaultBackupFile = "password-store.plain"

func main() {
	homeDir, _ := os.UserHomeDir()
	var defaultPasswordStore = filepath.Join(homeDir, ".password-store/")
	var debugMode bool
	app := &cli.App{
		Name:  "go_shell",
		Usage: "scripts de linux",
		Commands: []*cli.Command{
			{
				Name:  "pass",
				Usage: "comandos relativos ao comando pass",
				Subcommands: []*cli.Command{
					{
						Name: "restore",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "backup-file",
								Aliases: []string{"t"},
								Value:   defaultBackupFile + ".gpg",
							},
							&cli.StringFlag{
								Name:    "password-store",
								Aliases: []string{"d"},
								EnvVars: []string{"PASSWORD_STORE_DIR"},
								Value:   defaultPasswordStore,
							},
							&cli.StringFlag{
								Name:    "public-key",
								Aliases: []string{"k"},
								Value:   "",
							},
						},
						Action: func(c *cli.Context) error {
							target := c.String("backup-file")
							goshell.Restore(target)

							return nil
						},
					},
					{
						Name: "backup",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "password-store",
								Aliases: []string{"d"},
								EnvVars: []string{"PASSWORD_STORE_DIR"},
								Value:   defaultPasswordStore,
							},
							&cli.StringFlag{
								Name:    "backup-file",
								Aliases: []string{"t"},
								Value:   defaultBackupFile,
							},
						},
						Action: func(c *cli.Context) error {
							passwordStore := c.String("password-store")
							target := c.String("backup-file")
							goshell.Backup(target, passwordStore)
							return nil
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"D"},
				Destination: &debugMode,
			},
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
