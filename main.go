package main

import (
	"fmt"
	"io/ioutil"
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
						Action: func(c *cli.Context) error {
							target := c.String("backup_file")
							dFile := goshell.DecryptFile(target + ".gpg")
							if (debugMode){
								fmt.Println(dFile)
							}

							return nil
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "backup_file",
								Aliases: []string{"t"},
								Value:   defaultBackupFile,
							},
							&cli.StringFlag{
								Name:    "password-store",
								Aliases: []string{"d"},
								EnvVars: []string{"PASSWORD_STORE_DIR"},
								Value:   defaultPasswordStore,
							},
							&cli.StringFlag{
								Name: "public-key",
								Aliases: []string{"k"},
								Value: "",
							},
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
								Name:    "backup_file",
								Aliases: []string{"t"},
								Value:   defaultBackupFile,
							},
						},
						Action: func(c *cli.Context) error {

							passwordStore := c.String("password-store")
							target := c.String("backup_file")

							list := goshell.GetGpgFilePaths(passwordStore)

							allPassInfos := make([]goshell.GpgPassInfo, 0)

							for _, pn := range list {
								password := goshell.OpenPass(passwordStore, pn)
								allPassInfos = append(allPassInfos, password)
							}

							if goshell.CheckFileExists(target) {
								if target != defaultPasswordStore {
									log.Fatal("File exists: ", target)
								}
							}

							ioutil.WriteFile(target, goshell.ToJsonStr(allPassInfos), 0644)
							_, err := goshell.EncryptFile(target, "")
							if err != nil {
								log.Fatal(err)
							}
							if !debugMode {
								err = os.Remove(target)
							}
							if err != nil {
								log.Fatal(err)
							}

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
