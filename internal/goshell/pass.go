package goshell

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexpfx/go_shell/internal/util"
)

const GPG_EXT = ".gpg"

func NewBackup(passwordStore, target string) backup {
	return backup{
		passwordStore: passwordStore,
		target:        target,
	}
}

func NewRestore(prefix, passwordStoreDir, backupFilePath string, update bool) restore {
	return restore{
		targetPasswordStore: passwordStoreDir,
		backupFilePath:      backupFilePath,
		prefix:              prefix,
		update:              update,
	}
}

type Backup interface {
	Do()
}

type Restore interface {
	Do()
}

type backup struct {
	passwordStore string
	target        string
}

type restore struct {
	targetPasswordStore string
	backupFilePath      string
	prefix              string
	update              bool
}

func (b backup) Do() {
	v("iniciando backup...")
	v("store: ", b.passwordStore)
	v("target: ", b.target)

	if util.CheckFileExists(b.target) {
		f("arquivo de destino existe")
	}

	gpgFiles := getGpgFiles(b.passwordStore)

	allPassInfos := make([]GpgPassInfo, 0)
	for _, path := range gpgFiles {
		dec := decrypt(path)
		allPassInfos = append(allPassInfos, GpgPassInfo{
			Password: strings.TrimRight(string(dec), "\n"),
			PassName: getPassName(b.passwordStore, path),
			FilePath: path,
		})
	}

	doBackup(allPassInfos, b.target)
}

type GpgPassInfo struct {
	Password string `json:"password,omitempty"`
	PassName string `json:"pass_name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
}

func List() {
	l, err := util.ExecCmd([]string{"pass", "list"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(l))
}

func (r restore) Do() {
	v("iniciando restore...")

	str := decrypt(r.backupFilePath)
	v("descriptografando arquivo... ", r.backupFilePath)
	passInfo := make([]GpgPassInfo, 0)

	err := json.Unmarshal(str, &passInfo)
	if err != nil {
		f(err.Error())
	}
	for _, p := range passInfo {
		v("extraindo... ", p.PassName)
		targetGpgFile := filepath.Join(r.targetPasswordStore, r.prefix+p.PassName+GPG_EXT)

		if util.CheckFileExists(targetGpgFile) {
			if !r.update {
				v("Arquivo existe e não será sobreescrito: ", targetGpgFile)
				continue
			}
		}

		tmpFile := util.TempFile()
		v("criou arquivo temporario... ", tmpFile.Name())
		defer tmpFile.Close()

		_, err := fmt.Fprintln(tmpFile, p.Password)
		if err != nil {
			f(err.Error())
		}

		out := encryptFile(tmpFile.Name())
		if err != nil {
			e(out)
			f(err.Error())
		}

		v("movendo arquivo...", targetGpgFile, " ... ", tmpFile.Name()+GPG_EXT)
		util.MoveFile(targetGpgFile, tmpFile.Name()+GPG_EXT)

		os.Remove(tmpFile.Name())

	}
}

func encryptFile(file string) []byte {
	//cmdArgs := []string{"gpg", "--pinentry-mode=loopback", "--passphrase", password, "-c", file}
	v("codificando arquivo...", file)
	cmdArgs := []string{"gpg", "--pinentry-mode=loopback", "-c", file}
	out, err := util.ExecCmd(cmdArgs)
	if err != nil {
		e("erro ao codificar arquivo: ")
		e(err)
	}
	return out
}

func decrypt(file string) []byte {
	cmdArgs := []string{"gpg2", "-d", "--quiet", "--yes",
		"--compress-algo=none", "--pinentry-mode=loopback",
		file}
	out, err := util.ExecCmd(cmdArgs)

	if err != nil {
		e(string(out))
		f("não pode descriptografar arquivo.")
	}

	return out
}

func getPassName(baseDir, fullPath string) string {
	return strings.TrimSuffix(strings.TrimPrefix(fullPath, baseDir+string(os.PathSeparator)), filepath.Ext(fullPath))
}

func doBackup(passiInfos []GpgPassInfo, target string) {
	tmpFile := util.TempFile()
	v("arquivo temporario criado...", tmpFile.Name())

	bs, err := json.MarshalIndent(passiInfos, "", "   ")
	if err != nil {
		e("erro ao serializar: ")
		f(err.Error())
	}

	err = ioutil.WriteFile(tmpFile.Name(), bs, 0644)
	if err != nil {
		e("erro ao gravar arquivo temporario: ")
		f(err.Error())
	}

	encryptFile(tmpFile.Name())

	createdFile := tmpFile.Name() + GPG_EXT

	v("movendo arquivo criptografado")
	backupFile := strings.TrimSuffix(target, GPG_EXT) + GPG_EXT

	v("movendo arquivo...", createdFile, " ... ", backupFile)
	util.MoveFile(backupFile, createdFile)
	i("backup criado em ", backupFile)

	v("removendo arquivo temporário")
	err = os.Remove(tmpFile.Name())
	if err != nil {
		e("erro ao gravar remover arquivo temporário: ")
		f(err.Error())
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

		if !strings.EqualFold(filepath.Ext(path), GPG_EXT) {
			return nil
		}
		files = append(files, path)
		v("obteve arquivo gpg: \n", path)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return files
}

func i(msg ...interface{}) {
	for _, m := range msg {
		fmt.Println(m)
	}

}

func v(msg ...interface{}) {
	sb := strings.Builder{}

	for _, m := range msg {
		sb.WriteString(fmt.Sprintf("%s ", m))
	}
	log.Println(sb.String())
}
func f(err string) {
	log.Fatalln(err)

}
func e(msg ...interface{}) {
	sb := strings.Builder{}

	for _, m := range msg {
		sb.WriteString(fmt.Sprintf("%s ", m))
	}

	log.Println(sb.String())
}
