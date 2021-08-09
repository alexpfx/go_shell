package util

import (
	"encoding/json"

	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"golang.org/x/term"
)

func ToJsonStr(obj interface{}) []byte {
	bs, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return bs
}

func CheckFileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {		
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("Check file error: ", path)
	}
	return !stat.IsDir()

}

func MoveFile(targetPath string, originalPath string) {
	
	originalFile, err := os.Open(originalPath)
	if err != nil {
		log.Fatal(err)
	}

	defer originalFile.Close()
	defer os.Remove(originalFile.Name())

	copyFile, err := os.Create(targetPath)
	if err != nil {
		log.Fatal(err)
	}
	defer copyFile.Close()

	_, err = io.Copy(copyFile, originalFile)

	if err != nil {
		log.Fatal(err)
	}
}

func ExecCmd(cmdArgs []string) ([]byte, error) {
	log.Println(cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	out, err := cmd.CombinedOutput()
	return out, err
}

func PromptPass(prompt string) string {
	log.Println(prompt)
	bs, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	return string(bs)
}

func TempFile() *os.File {
	f, err := ioutil.TempFile("", "gsh*")
	if err != nil {
		log.Fatal(f)
	}
	return f
}
