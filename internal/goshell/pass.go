package goshell

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GpgPassInfo struct {
	Password string `json:"password,omitempty"`
	PassName string `json:"pass_name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
}

func Backup(target, passwordStore string) {
	if checkFileExists(target) {
		if target != passwordStore {
			log.Fatal("File exists: ", target)
		}
	}

	allPassInfos := make([]GpgPassInfo, 0)
	paths := getGpgFiles(passwordStore)
	for _, path := range paths {
		dec := decrypt(path)

		allPassInfos = append(allPassInfos, GpgPassInfo{
			Password: strings.TrimRight(string(dec), "\n"),
			PassName: getPassName(passwordStore, path),
			FilePath: path,
		})
	}
	createBackupFile(allPassInfos, target)
}

func Restore(backupPath, targetFolder, prefix string, override bool) {
	str := decrypt(backupPath)
	passInfo := make([]GpgPassInfo, 0)

	err := json.Unmarshal(str, &passInfo)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range passInfo {
		target := filepath.Join(prefix, targetFolder, p.PassName)
		if checkFileExists(target) {
			if !override {
				fmt.Println("File exists: ", p.PassName)
				continue
			}
			tmp, err := ioutil.TempFile("", "*")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(tmp.Name())

			out, err := encryptFile(tmp.Name(), target)
			if err != nil {
				fmt.Println(string(out))
				log.Fatal(err)
			}

		}

	}
}

func encryptFile(file string, output string) ([]byte, error) {
	//cmdArgs := []string{"gpg", "--pinentry-mode=loopback", "--passphrase", password, "-c", file}
	cmdArgs := []string{"gpg", "--pinentry-mode=loopback", "-c", file, "-o", output}
	out, err := callCmd(cmdArgs)
	return out, err
}

func callCmd(cmdArgs []string) ([]byte, error) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	out, err := cmd.CombinedOutput()
	return out, err
}

func decrypt(file string) []byte {
	cmdArgs := []string{"gpg2", "-d", "--quiet", "--yes",
		"--compress-algo=none", "--pinentry-mode=loopback",
		file}
	out, err := callCmd(cmdArgs)

	if err != nil {
		fmt.Println(string(out))
		log.Fatalf("Cannot decrypt: %s", file)
	}

	return out
}

func getPassName(baseDir, fullPath string) string {
	return strings.TrimSuffix(strings.TrimPrefix(fullPath, baseDir+string(os.PathSeparator)), filepath.Ext(fullPath))
}

func createBackupFile(passiInfos []GpgPassInfo, target string) {
	tmp, err := ioutil.TempFile("", "*")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmp.Name())

	bs, err := json.MarshalIndent(passiInfos, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(tmp.Name(), bs, 0644)
	if err != nil {
		log.Fatal(err)
	}

	out, err := encryptFile(tmp.Name(), target+".gpg")
	if err != nil {
		fmt.Println(string(out))
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
