package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alexpfx/go_shell/internal/goshell"
	"github.com/urfave/cli/v2"
)

const defaultBackupFile = ".passbackup"

func main() {
	homeDir, _ := os.UserHomeDir()
	var defaultPasswordStore = filepath.Join(homeDir, ".password-store/")

	app := &cli.App{
		Name:  "go_shell",
		Usage: "scripts de linux",
		Commands: []*cli.Command{
			{
				Name:  "pass",
				Usage: "comandos relativos ao comando pass",
				Subcommands: []*cli.Command{
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
								Name:    "target",
								Aliases: []string{"t"},
								Value:   fmt.Sprintf("%s_%d", defaultBackupFile, time.Now().Local().Unix()),
							},
						},
						Action: func(c *cli.Context) error {

							passwordStore := c.String("password-store")
							target := c.String("target")

							list := goshell.GetGpgFilePaths(passwordStore)

							allPassInfos := make([]goshell.GpgPassInfo, 0)

							for _, pn := range list {
								password := goshell.OpenPass(passwordStore, pn)
								allPassInfos = append(allPassInfos, password)
							}
							
							
							if goshell.CheckFileExists(target){
								log.Fatal("File exists: ", target)
							}
							ioutil.WriteFile(target, goshell.ToJsonStr(allPassInfos), 0644)
							_, err := goshell.EncryptFile(target, "")
							if err != nil {								
								log.Fatal(err)
							}
														

							return nil
						},
					},
				},
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
