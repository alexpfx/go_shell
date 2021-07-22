package goshell

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const pass = "pass"

type GpgPassInfo struct {
	Password string `json:"password,omitempty"`
	PassName string `json:"pass_name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
}

func EncryptWithKey(targetFile, key string) (string, error) {
	return "nil", nil
}

func encryptFile(file, password string) error {
	cmdArgs := []string{"gpg", "--pinentry-mode=loopback", "--passphrase", password, "-c", file}
	_, err := callCmd(cmdArgs)
	return err
}

func DecryptFile(filePath string) string {
	d := decrypt(filePath, "")
	return d
}

func AllPassInfos() {

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

func decrypt(file, pass string) string {
	cmdArgs := []string{"gpg2", "-d", "--quiet", "--yes",
		"--compress-algo=none", "--pinentry-mode=loopback",
		"--passphrase", pass, file}
	out, err := callCmd(cmdArgs)
	if err != nil {
		log.Fatalf("Cannot decrypt: %s", file)
	}

	return out
}

func getPassKeys(storePath string) []string {
	return getGpgFiles(storePath)
}

func getPassName(baseDir, fullPath string) string {
	return strings.TrimSuffix(strings.TrimPrefix(fullPath, baseDir+string(os.PathSeparator)), filepath.Ext(fullPath))
}

func Backup(target, passwordStore string) {
	if CheckFileExists(target) {
		if target != passwordStore {
			log.Fatal("File exists: ", target)
		}
	}
	pass := promptPass("Passphrase: ")
	allPassInfos := make([]GpgPassInfo, 0)
	paths := getGpgFiles(passwordStore)
	for _, path := range paths {
		dec := decrypt(path, pass)

		allPassInfos = append(allPassInfos, GpgPassInfo{
			Password: strings.TrimRight(dec, "\n"),
			PassName: getPassName(passwordStore, path),
			FilePath: path,
		})
	}
	createBackupFile(allPassInfos, target, pass)
}

func createBackupFile(passiInfos []GpgPassInfo, target, pass string) {
	defer os.Remove(target)

	bs, err := json.MarshalIndent(passiInfos, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(target, bs, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = encryptFile(target, pass)
	if err != nil {
		log.Fatal("Cannot encrypt file: ", target)
	}

}

func getGpgFiles(baseDir string) []string {
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
