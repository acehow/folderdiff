package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

var filemap = make(map[string]string)
var destmap = make(map[string]string)
var filemapmd5 = make(map[string]string)
var destmapmd5 = make(map[string]string)
var wg = sync.WaitGroup{}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("need 2 parameters ")
		os.Exit(1)
	}
	
	//src := "c:/wuhao/tmp/test1"
	//desc := "c:/wuhao/tmp/test2"
	src := strings.TrimRight(os.Args[1],"/")
	desc := strings.TrimRight(os.Args[2],"/")
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	GetAllFile(src, len(src)+1)
	GetAllFile2(desc, len(desc)+1)

	for k := range destmap {
		if _, ok := filemap[k]; ok {
			filemap[k] = "1"
		} else {
			fmt.Println("In dest:" + k)
		}
	}

	for k, v := range filemap {
		if v != "1" {
			fmt.Println("In source:" + k)
		} else {
			filemapmd5[k] = ""
			destmapmd5[k] = ""
		}
	}

	/*
		for k, v := range filemap {
			if v == "1" {
				if md5f(src+"/"+k) != md5f(desc+"/"+k) {
					fmt.Println("Diff:" + k)
				}
			}
		}
	*/
	wg.Add(2)
	go setmap(filemapmd5, src+"/")
	go setmap(destmapmd5, desc+"/")
	wg.Wait()
	for k := range filemapmd5 {
		if filemapmd5[k] != destmapmd5[k] {
			fmt.Println("Diff:" + k)
		}
	}
}

func GetAllFile(pathname string, strip int) error {
	rd, err := os.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFile(pathname+"/"+fi.Name(), strip)
		} else {
			fn := pathname + "/" + fi.Name()
			filemap[fn[strip:]] = ""
			//fmt.Println(fn[strip:])
		}
	}
	return err
}

func GetAllFile2(pathname string, strip int) error {
	rd, err := os.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFile(pathname+"/"+fi.Name(), strip)
		} else {
			fn := pathname + "/" + fi.Name()
			destmap[fn[strip:]] = ""
			//fmt.Println(fn[strip:])
		}
	}
	return err
}

func setmap(tmap map[string]string, prefix string) {
	for k := range tmap {
		tmap[k] = md5f(prefix + k)
	}
	wg.Done()
}

func md5f(fName string) string {
	f, e := os.Open(fName)
	defer fmt.Println(f.Close().Error())
	if e != nil {
		fmt.Println(e.Error())
	}
	h := md5.New()
	_, e = io.Copy(h, f)
	if e != nil {
		fmt.Println(e.Error())
	}
	return hex.EncodeToString(h.Sum(nil))
}
