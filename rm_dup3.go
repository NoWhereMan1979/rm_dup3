package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var wd = "/home/happy/hdd460/hae/solopianoradio/Whisperings- Solo Piano Radio"
var match = regexp.MustCompile(`\s\([^\d]*(\d+)[^\d]*\)`)

func main() {
	err := os.Chdir(wd)
	if err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
	}
	groupMap := make(map[string][]os.FileInfo)
	for _, f := range files {
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
	fmt.Println(len(groupMap))
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
	var dirName = "./pldnufhsd"
	e := os.Mkdir(dirName, 0777)
	if e != nil {
		os.Exit(14)
	}
	for _, x := range lastList {
		rer := os.Rename(x.Name(), dirName+"/"+trim(x.Name()))
		if rer != nil {
			fmt.Println("error moving file ", x.Name(), "to", dirName+"/"+trim(x.Name()))
		}
	}
	fmt.Println("moved files", len(lastList))
}
func trim(n string) string {
	var loc = match.FindStringIndex(n)
	if loc != nil {
		n = n[:loc[0]] + n[loc[1]:]
	}
	return n
}
