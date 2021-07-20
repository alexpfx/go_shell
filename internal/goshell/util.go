package goshell

import (
	"encoding/json"
	"log"
	"os"
)

func ToJsonStr(obj interface{}) []byte{
	bs, err := json.MarshalIndent(obj, "", "    ")
	if err != nil{
		log.Fatal(err)
	}
	return bs
}

func CheckFileExists(path string) bool{
	stat, err := os.Stat(path)
	if err != nil{
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("Check file error: ", path)
	}
	return !stat.IsDir()

}