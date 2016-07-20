package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/macaron.v1"
)

func findLength(str string) int64 {
	re, _ := regexp.Compile(`\nLength: (.+?) `)
	subs := re.FindAllStringSubmatch(str, 1)
	if len(subs) != 0 {
		i, _ := strconv.ParseInt(subs[0][1], 10, 64)
		return i
	}
	return -1
}

func findName(str string) string {
	re, _ := regexp.Compile("Saving to: [`‘](.*)['’]")
	subs := re.FindAllStringSubmatch(str, 1)
	if len(subs) != 0 {
		return subs[0][1]
	}

	re, _ = regexp.Compile("=> [`‘](.*)['’]")
	subs = re.FindAllStringSubmatch(str, 1)
	if len(subs) != 0 {
		return subs[0][1]
	}
	return ""
}

func getFileSize(fp string) int64 {
	finfo, err := os.Stat(fp)

	if err != nil {
		return -1
	}

	return finfo.Size()
}

func isWgetExit(str string) bool {
	re, _ := regexp.Compile(`failed: .*\n`)
	subs := re.FindAllStringSubmatch(str, -1)
	sl := len(subs)
	if sl != 0 {
		if strings.Contains(subs[sl-1][0], "No route") {

		} else {
			return true
		}
	}
	return false
}

var fileManager = struct {
	sync.RWMutex
	m map[string]int64
}{m: make(map[string]int64)}

var obfEncoding = base64.NewEncoding("-+ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func decode(str string) string {
	tmp, _ := obfEncoding.DecodeString(str)
	return string(tmp)
}

func WgetHandler(req *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	url := req.URL.Path
	url = decode(url[7:])

	args := strings.Split(url, " ")
	dir := "downloads"
	root, _ := filepath.Abs(gcfg.root)

	fspath := filepath.Join(root, dir)
	os.MkdirAll(fspath, os.ModePerm)

	cmd := exec.Command("wget", "--content-disposition", "-P", fspath)
	cmd.Args = append(cmd.Args, args...)
	log.Println("exec:", cmd.Args)

	var serr bytes.Buffer
	cmd.Stderr = &serr
	err = cmd.Start()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var size int64
	var name string

	var goon bool = false
	for i := 0; i < 300; i++ {
		if isWgetExit(serr.String()) {
			http.Error(w, serr.String(), http.StatusInternalServerError)
			return
		}

		time.Sleep(time.Millisecond * 100)
		size = findLength(serr.String())
		if size == -1 {
			continue
		}
		name = findName(serr.String())
		if name == "" {
			continue
		}

		goon = true
		break
	}

	if !goon {
		http.Error(w, serr.String(), http.StatusInternalServerError)
		return
	}

	name = path.Base(name)

	fileManager.Lock()
	fileManager.m[name] = size
	fileManager.Unlock()

	json := fmt.Sprintf(`{"fname":"%s","fsize":%d}`, name, size)
	fmt.Println(json)

	w.Write([]byte(json))
}

func WstatHandler(req *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
	url := req.URL.Path
	fname := url[8:]

	dir := "downloads"
	root, _ := filepath.Abs(gcfg.root)
	fspath := filepath.Join(root, dir, fname)

	downloaded := getFileSize(fspath)

	json := fmt.Sprintf(`{"fname":"%s","downloaded":%d}`, fname, downloaded)
	fmt.Println(json)

	w.Write([]byte(json))
}
