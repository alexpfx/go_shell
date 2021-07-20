package goshell

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

const pass = "pass"

type GpgPassInfo struct {
	Password string `json:"password,omitempty"`
	PassName string `json:"pass_name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
}

func EncryptFile(file, password string) (string, error) {
	cmdArgs := []string{"gpg", "--pinentry-mode=loopback", "--passphrase", password, "-c", file}
	out, err := callCmd(cmdArgs)
	if err != nil {
		if password == "" {
			fmt.Println("Output file password: ")
			pp := promptPass()
			fmt.Println("Confirm: ")
			rp := promptPass()
			if (pp != rp){
				log.Fatal("Typed passwords doesn't match")
			}

			return EncryptFile(file, pp)
		}
		return out, err
	}
	return out, err
}

func OpenPass(storePath, file string) GpgPassInfo {
	return CallGpg(storePath, file, "")
}

func CallPass(args []string) string {
	cmd := exec.Command(pass, args...)
	password, err := cmd.CombinedOutput()
	if err != nil {

		log.Fatal(err)
	}
	return string(password)
}

func callCmd(cmdArgs []string) (string, error) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func promptPass() string {
	bs, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	return string(bs)
}

func CallGpg(basedir, file, pass string) GpgPassInfo {
	cmdArgs := []string{"gpg2", "-d", "--quiet", "--yes",
		"--compress-algo=none", "--pinentry-mode=loopback",
		"--passphrase", pass, file}
	out, err := callCmd(cmdArgs)
	fmt.Println(pass)
	if err != nil {
		if pass == "" {
			fmt.Println("Input files password: ")
			pp := promptPass()
			return CallGpg(basedir, file, pp)
		}
		
		log.Fatal("Não foi possível logar no arquivo de origem.")				
	}

	return GpgPassInfo{
		Password: out,
		FilePath: file,
		PassName: getPassName(basedir, file),
	}
}

func GetPassKeys(storePath string) []string {
	return GetGpgFilePaths(storePath)
}

func getPassName(baseDir, fullPath string) string {
	return strings.TrimSuffix(strings.TrimPrefix(fullPath, baseDir+string(os.PathSeparator)), filepath.Ext(fullPath))
}

func GetGpgFilePaths(baseDir string) []string {
	files := make([]string, 0)

	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if d.Type().IsDir() {
			return nil
		}

		if !strings.EqualFold(filepath.Ext(path), ".gpg") {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return files
}
