package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var NOTES_PATH string = userHomeDir() + "/switchdrive/Notes"

func main() {
	log.Println("Hello Notes")
	log.Println(NOTES_PATH)

	if len(os.Args) == 1 {
		content := ReadFileToString(filepath.Join(NOTES_PATH, todayFilename()))
		lines := readLineByLine(content)
		printNote(lines)

	} else {
		if os.Args[1] == "todo" {
			ParseAllFiles(NOTES_PATH, "todo", 0)
		}

		if os.Args[1] == "search" {
			ParseAllFiles(NOTES_PATH, os.Args[2], 0)
		}
	}
}

func listTodos() {

}

func readLineByLine(s string) []string {
	lines := strings.Split(s, "\n")

	return lines
}

func printNote(as []string) {

	for i, _ := range as {
		if strings.HasPrefix(as[i], "##") {
			fmt.Println(strings.ToUpper(as[i]))
		} else if strings.HasPrefix(as[i], "*") {
			fmt.Println("  " + as[i])
		} else {
			fmt.Println(as[i])
		}

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

func ParseAllFiles(path string, filter string, of int) {

	file_list, _ := ListDir(NOTES_PATH)

	for _, i := range file_list {
		if !(i.IsDir()) {
			s := ReadFileToString(filepath.Join(NOTES_PATH, i.Name()))
			sarr := readLineByLine(s)
			for j, _ := range sarr {
				if strings.Contains(sarr[j], filter) {
					fmt.Println(i.Name() + "  " + sarr[j])
				}
			}
		}
	}

}
