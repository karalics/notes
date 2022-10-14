package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var NOTES_PATH string = userHomeDir() + "/switchdrive/Notes"

var colorBlue string = "\033[34m"
var colorReset string = "\033[0m"
var colorPurple string = "\033[35m"
var colorRed string = "\033[31m"

func main() {
	log.Println(NOTES_PATH)

	today_filename := filepath.Join(NOTES_PATH, todayFilename())

	_, file_err := os.Stat(today_filename)

	if file_err != nil {
		log.Println("File does not exist, creating a new one")
		emptyFile, err := os.Create(today_filename)
		if err != nil {
			log.Fatal(err)
		}
		emptyFile.Close()
	}

	if len(os.Args) == 1 {
		content := ReadFileToString(today_filename)
		lines := readLineByLine(content)
		printNote(lines)

	} else {
		if os.Args[1] == "todo" {
			ParseAllFiles(NOTES_PATH, "TODO", 0)
		}

		if os.Args[1] == "search" {
			ParseAllFiles(NOTES_PATH, os.Args[2], 0)
		}

		if os.Args[1] == "add" {
			appendToNote(today_filename, mergeStringArray(os.Args[2:]))
		}

		if os.Args[1] == "done" {
			indexnumber, _ := strconv.Atoi(os.Args[2])

			line := lineByNumber("tindex", indexnumber)
			args := strings.Split(line, ":::")
			filename := args[1]
			linenumber, err := strconv.Atoi(args[2])
			if err != nil {
				log.Fatal(err)
			}
			writeDone(filename, linenumber)

			ParseAllFiles(NOTES_PATH, "TODO", 0)
		}
	}
}

func listTodos() {

}

func mergeStringArray(sarr []string) string {
	var mystring string
	for _, v := range sarr {
		mystring = mystring + v + " "
	}

	return mystring
}

func readLineByLine(s string) []string {
	lines := strings.Split(s, "\n")

	return lines
}

func printNote(as []string) {

	//colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	//colorYellow := "\033[33m"
	//colorBlue := "\033[34m"
	//colorPurple := "\033[35m"
	//colorCyan := "\033[36m"
	//colorWhite := "\033[37m"

	for i, _ := range as {
		if strings.HasPrefix(as[i], "##") {
			fmt.Println(colorGreen, strings.ToUpper(as[i]), colorReset)
		} else if strings.HasPrefix(as[i], "*") {
			fmt.Println("  " + as[i])
		} else if strings.Contains(as[i], "TODO") {
			fmt.Println(colorRed, as[i], colorReset)
		} else {
			fmt.Println(as[i])
		}

	}

}

func writeDone(filename string, linenumber int) {

	file := filepath.Join(NOTES_PATH, filename)

	input, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(input), "\n")

	pos := linenumber - 1

	lines[pos] = strings.Replace(lines[pos], "TODO", "DONE", -1)

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
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

func lineByNumber(filename string, linenumber int) (line string) {

	// Not Used at the moment

	readFile, err := os.Open(filepath.Join(NOTES_PATH, filename))

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	fileindex := 1
	for fileScanner.Scan() {
		if fileindex == linenumber {
			return fileScanner.Text()
		}
		fileindex++
	}

	readFile.Close()

	return "No Line Found"
}

func ParseAllFiles(path string, filter string, of int) {

	file_list, _ := ListDir(NOTES_PATH)

	var tindex int = 1
	_, err := os.Create(filepath.Join(NOTES_PATH, "tindex"))
	f, err := os.OpenFile(filepath.Join(NOTES_PATH, "tindex"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	var re = regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)

	for _, i := range file_list {
		if !(i.IsDir()) && (i.Name() != "tindex") {
			s := ReadFileToString(filepath.Join(NOTES_PATH, i.Name()))
			sarr := readLineByLine(s)
			for j, _ := range sarr {

				if strings.Contains(sarr[j], filter) {
					output := strconv.Itoa(tindex) + ":::" + i.Name() + ":::" + strconv.Itoa(j+1)
					if re.MatchString(sarr[j]) {
						fmt.Println(tindex, colorBlue, i.Name(), colorReset, colorPurple, j+1, colorReset, colorRed, sarr[j], colorReset)
					} else {
						fmt.Println(tindex, colorBlue, i.Name(), colorReset, colorPurple, j+1, colorReset, sarr[j])
					}
					tindex++
					if _, err := f.Write([]byte(output + "\n")); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func appendToNote(path string, s string) {

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(s + "\n\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
