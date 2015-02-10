package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"

	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type ReloadMux struct {
	port        int
	root        string
	dirListTmpl *template.Template
	docTmpl     *template.Template
	private     bool
}

var reloadCfg = ReloadMux{}

func showDoc(w http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(200)
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
}

func formatSize(file os.FileInfo) string {
	if file.IsDir() {
		return "-"
	}
	size := file.Size()
	switch {
	case size > 1024*1024:
		return fmt.Sprintf("%.1fM", float64(size)/1024/1024)
	case size > 1024:
		return fmt.Sprintf("%.1fk", float64(size)/1024)
	default:
		return strconv.Itoa(int(size))
	}
	return ""
}

func dirList(w http.ResponseWriter, f *os.File) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if dirs, err := f.Readdir(-1); err == nil {
		files := make([]map[string]string, len(dirs)+1)
		files[0] = map[string]string{
			"name": "..", "href": "..", "size": "-", "mtime": "-",
		}
		for i, d := range dirs {
			href := d.Name()
			if d.IsDir() {
				href += "/"
			}
			files[i+1] = map[string]string{
				"name":  d.Name(),
				"href":  href,
				"size":  formatSize(d),
				"mtime": d.ModTime().Format("2006-01-02 15:04:05"),
			}
		}
		reloadCfg.dirListTmpl.Execute(w, map[string]interface{}{
			"dir":   f.Name(),
			"files": files,
		})
	}
}

func fileHandler(w http.ResponseWriter, path string, req *http.Request) {
	if path == "" {
		path = "."
	}
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		showDoc(w, req, err)
		return
	}
	defer f.Close()

	d, err1 := f.Stat()
	if err1 != nil {
		log.Println(err)
		showDoc(w, req, err)
		return
	}

	if d.IsDir() {
		dirList(w, f)
	} else {
		ctype := mime.TypeByExtension(filepath.Ext(path))
		if ctype != "" {
			// go return charset=utf8 even if the charset is not utf8
			idx := strings.Index(ctype, "; ")
			if idx > 0 {
				// remove charset; anyway, browsers are very good at guessing it.
				ctype = ctype[0:idx]
			}
			w.Header().Set("Content-Type", ctype)
		}
		w.WriteHeader(200)
		io.Copy(w, f)
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	fileHandler(w, path[1:], req)
}

func main() {
	flag.IntVar(&(reloadCfg.port), "port", 8000, "Which port to listen")
	flag.StringVar(&(reloadCfg.root), "root", ".", "Watched root directory for filesystem events, also the HTTP File Server's root directory")
	flag.BoolVar(&(reloadCfg.private), "private", false, "Only listen on lookback interface, otherwise listen on all interface")
	flag.Parse()

	t, _ := template.New("dirlist").Parse(DIR_HTML)
	reloadCfg.dirListTmpl = t
	t, _ = template.New("doc").Parse(HELP_HTML)
	reloadCfg.docTmpl = t

	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := os.Chdir(reloadCfg.root); e != nil {
		log.Panic(e)
	}
	http.HandleFunc("/", handler)

	int := ":" + strconv.Itoa(reloadCfg.port)
	p := strconv.Itoa(reloadCfg.port)
	mesg := "; please visit http://127.0.0.1:" + p
	if reloadCfg.private {
		int = "localhost" + int
		log.Printf("listens on 127.0.0.1@" + p + mesg)
	} else {
		log.Printf("listens on 0.0.0.0@" + p + mesg)
	}
	if err := http.ListenAndServe(int, nil); err != nil {
		log.Fatal(err)
	}
}
