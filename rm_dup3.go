package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var path = flag.String("p", wd, "execution path")
var mvDest = flag.String("m", movingDestination, "moving destination directory")

var wd = "./"
var movingDestination = "mvd"
var regexStr = `\s\([^\d]*(\d+)[^\d]*\)`
var match = regexp.MustCompile(regexStr)
var spash = `################### File name dupplication remover ###################
Chooses most complete files from streamripper dupplicate ripped files.
Then renames to original file name and moves them to a directory
Use with -h for usage help
**********************************************************************`

func main() {
	flag.Parse()
	if flag.Arg(1) != "h" {
		fmt.Println(spash)
	}
	if err := loadRegex(); err != nil {
		log.Fatal("error loading regex.bytes", err)
	}
	wd = *path
	match = regexp.MustCompile(regexStr)
	chDir(wd)
	makeDir(movingDestination)
	groupMap := readNGroup()
	lastList := chooseSuits(groupMap)
	renameNMove(lastList)
}

func readNGroup() map[string][]os.FileInfo {
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
	}
	groupMap := make(map[string][]os.FileInfo)
	for _, f := range files {
		if !f.IsDir() {
			var nameF = f.Name()
			nameF = trim(nameF)
			if gr, ok := groupMap[nameF]; ok {
				gr = append(gr, f)
				groupMap[nameF] = gr
			} else {
				var g []os.FileInfo
				g = append(g, f)
				groupMap[nameF] = g
			}
		}
	}
	fmt.Println(len(groupMap))
	return groupMap
}

func chooseSuits(groupMap map[string][]os.FileInfo) []os.FileInfo {
	var lastList []os.FileInfo
	for _, g := range groupMap {
		var selected os.FileInfo
		var size = int64(0)
		for _, fi := range g {
			if fi.Size() > size {
				selected = fi
				size = fi.Size()
			}
		}
		lastList = append(lastList, selected)
	}
	return lastList
}

func renameNMove(lastList []os.FileInfo) {
	var count = 0
	for _, x := range lastList {
		rer := os.Rename(x.Name(), movingDestination+"/"+trim(x.Name()))
		if rer != nil {
			fmt.Println("error moving file ", x.Name(), "to", movingDestination+"/"+trim(x.Name()))
			continue
		}
		count++
	}
	fmt.Println("moved files", count)
}

func makeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		e := os.Mkdir(dir, 0777)
		if e != nil {
			fmt.Println(e.Error())
			os.Exit(1)
		}
	} else {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		if len(files) != 0 {
			fmt.Println("directory", dir, "exists and is not empty")
			os.Exit(3)
		}
	}
}

func chDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err := os.Chdir(dir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

func trim(n string) string {
	var loc = match.FindStringIndex(n)
	if loc != nil {
		n = n[:loc[0]] + n[loc[1]:]
	}
	return n
}

func writeBytes2File() error {
	file, err := os.OpenFile(
		"regex.bytes",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write bytes to file
	byteSlice := []byte(regexStr)
	bytesWritten, err := file.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
		// return err
	}
	log.Printf("Default regex.bytes file created. %d bytes.\n", bytesWritten)
	return nil
}

func loadRegex() error {
	b, err := ioutil.ReadFile("regex.bytes")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("regex.bytes file doesn't exist!")
			return writeBytes2File()
		}
		return err
	}
	regexStr = string(b)
	return nil
}
