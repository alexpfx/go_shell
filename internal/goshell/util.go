package goshell

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func ToJsonStr(obj interface{}) []byte {
	bs, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return bs
}

func checkFileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("Check file error: ", path)
	}
	return !stat.IsDir()

}

func promptPass(prompt string) string {
	fmt.Println(prompt)
	bs, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	return string(bs)
}

func GetPublicKey(storePath string) string {
	return ""

}
