package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/alexpfx/go_shell/internal/goshell"
	"github.com/urfave/cli/v2"
)

const defaultBackupFile = "pass-backup"

func main() {
	log.SetFlags(log.Llongfile)

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
						Name: "list",
						Action: func(c *cli.Context) error {
							goshell.List()
							return nil
						},
					},
					{
						Name: "restore",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "backup-file",
								Aliases: []string{"t"},
								Value:   defaultBackupFile + ".gpg",
							},
							&cli.BoolFlag{
								Name:    "update",
								Aliases: []string{"U"},
								Value:   false,
							},
							&cli.StringFlag{
								Name:  "prefix",
								Value: "restored/",
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
							passwordStore := c.String("password-store")
							prefix := c.String("prefix")
							forceUpdate := c.Bool("update")

							restore := goshell.NewRestore(prefix, passwordStore, target, forceUpdate)
							restore.Do()

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

							if debugMode {
								log.SetFlags(log.Llongfile)
								//log.SetFlags(0)
								//log.SetOutput(io.Discard)
							}

							backup := goshell.NewBackup(passwordStore, target)
							backup.Do()

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
