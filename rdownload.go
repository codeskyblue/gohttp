package main

import (
	"bufio"
	"bytes"
	"errors"
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
	"time"

	"gopkg.in/macaron.v1"
)

func findLength(str string) int64 {
	re, _ := regexp.Compile(`\nLength: (\d+) `)
	subs := re.FindAllStringSubmatch(str, 1)
	if len(subs) != 0 {
		i, _ := strconv.ParseInt(subs[0][1], 10, 64)
		return i
	}
	return -1
}

func findName(str string) string {
	re, _ := regexp.Compile("Saving to: `(.*)'")
	subs := re.FindAllStringSubmatch(str, 1)
	if len(subs) != 0 {
		return subs[0][1]
	}
	return ""
}

func getHeaderLen(str string) int {
	re, _ := regexp.Compile(`\n\n`)
	subs := re.FindAllStringSubmatchIndex(str, 1)
	if len(subs) != 0 {
		return subs[0][0] + 2
	}
	return -1
}

func findPercent(str string) int {
	re, _ := regexp.Compile(`(\d+)%`)
	subs := re.FindAllStringSubmatch(str, 1)
	if len(subs) != 0 {
		i, _ := strconv.Atoi(subs[0][1])
		return i
	}
	return -1
}

func WgetHandler(req *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	url := req.URL.Path
	url = url[7:]
	url = strings.Replace(url, "http:/", "http://", -1)
	url = strings.Replace(url, "https:/", "https://", -1)
	args := strings.Split(url, " ")
	dir := "downloads"
	root, _ := filepath.Abs(gcfg.root)

	fspath := filepath.Join(root, dir)
	os.MkdirAll(fspath, os.ModePerm)

	cmd := exec.Command("wget", "--content-disposition", "-P", fspath)
	cmd.Args = append(cmd.Args, args...)
	log.Println("exec: ", cmd.Args)

	var serr bytes.Buffer
	cmd.Stderr = &serr
	err = cmd.Start()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var size int64
	var name string
	var pos int
	var lastcount int = -1
	var errcount int = 0
	var goon bool = false
	for i := 0; i < 100; i++ {
		if i%10 == 0 {
			outcount := serr.Len()
			if outcount == lastcount {
				errcount++
			} else {
				lastcount = outcount
				errcount = 0
			}

			if errcount >= 5 {
				err = errors.New("exit while wget has no response")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
		pos = getHeaderLen(serr.String())
		if pos == -1 {
			continue
		}
		goon = true
		break
	}

	if !goon {
		err = errors.New("exit while can not get fileinfos ")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name = path.Base(name)
	//log.Printf("name:[%s], size:[%d], headerlen:[%d]\n", name, size, pos)
	var per int = 0

	//w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"fname":"%s","fsize":%d,"per":%d}`, name, size, per)
	w.(http.Flusher).Flush()

	serr.Next(pos)

	reader := bufio.NewReader(&serr)
	var lastline, tmpline []byte

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				tmpline = line
				break
			} else {
				lastline = append(tmpline, line...)
				tmpline = tmpline[:0]
			}
		}
		if len(tmpline) == 0 {
			per = 100
			//log.Printf("%d%%\n", per)
			fmt.Fprintf(w, `{"fname":"%s","fsize":%d,"per":%d}`, name, size, per)
			w.(http.Flusher).Flush()
			break
		}
		per = findPercent(string(lastline))
		if per != -1 {
			//log.Printf("%d%%\n", per)
			fmt.Fprintf(w, `{"fname":"%s","fsize":%d,"per":%d}`, name, size, per)
			w.(http.Flusher).Flush()
			if per == 100 {
				break
			}
		}
	}
	err = cmd.Wait()
	return
}
