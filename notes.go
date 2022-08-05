package main


import (
	"io/ioutil"
	"io/fs"
	"path/filepath"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var NOTES_PATH string = userHomeDir() + "/switchdrive/Notes"

func main(){
	log.Println("Hello Notes")
	log.Println(NOTES_PATH)

	if len(os.Args) == 1 {
		content := ReadFileToString(filepath.Join(NOTES_PATH, todayFilename()))
		fmt.Println(content)
		todayFilename()
		ParseAllFiles(NOTES_PATH)
	}
}

func todayFilename() string {
	dt := time.Now()
	day := dt.Format("2006-01-02")
	return day + ".md"
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func ReadFileToString(path string) string {
	var s string

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	b, e := ioutil.ReadAll(file)
	if e != nil {
		log.Println(e)
	}
	s = string(b[:])

	return s
}

func ListDir(path string) ([]fs.FileInfo, error) {

	filesinfo, err := ioutil.ReadDir(path)
	if err != nil {
		log.SetPrefix("function listDir: ")
		log.Println(err)
	}

	return filesinfo, err

}

func ParseAllFiles(path string) {

	file_list, _ := ListDir(NOTES_PATH)

	for _, i := range file_list {
		if !(i.IsDir()) {

			fmt.Println(i.Name())

		}

	}
}