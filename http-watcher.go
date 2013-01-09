package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

type FileEvent struct {
	File  string
	Event string
}

type Client struct {
	buf  *bufio.ReadWriter
	conn net.Conn
}

type ReloadMux struct {
	mu            sync.Mutex
	port          int
	ignores       string
	ignorePattens []*regexp.Regexp
	command       string
	root          string
	eventsCh      chan []FileEvent
	reloadJs      *template.Template
	dirListTmpl   *template.Template
	docTmpl       *template.Template
	clients       []Client
	private       bool
	proxy         int
}

var reloadCfg = ReloadMux{
	eventsCh: make(chan []FileEvent),
	clients:  make([]Client, 0),
}

func shouldIgnore(path string) bool {
	for _, p := range reloadCfg.ignorePattens {
		if p.Find([]byte(path)) != nil {
			return true
		}
	}
	return false
}

func showDoc(w http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(200)
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	abs, _ := filepath.Abs(".")
	reloadCfg.docTmpl.Execute(w, map[string]interface{}{
		"error":   err,
		"dir":     abs,
		"ignores": reloadCfg.ignorePattens,
		"hosts":   publicHosts(),
	})

}

func publicHosts() []string {
	ips := make([]string, 0)
	if reloadCfg.private {
		ips = append(ips, "127.0.0.1:"+strconv.Itoa(reloadCfg.port))
	} else {
		if addrs, err := net.InterfaceAddrs(); err == nil {
			r, _ := regexp.Compile(`(\d+\.){3}\d+`)
			for _, addr := range addrs {
				ip := addr.String()
				if strings.Contains(ip, "/") {
					ip = strings.Split(ip, "/")[0]
				}
				if r.Match([]byte(ip)) {
					ips = append(ips, ip+":"+strconv.Itoa(reloadCfg.port))
				}
			}
		}
	}
	return ips
}

func getAllFileMeta() map[string]time.Time {
	files := map[string]time.Time{}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil { // TODO permisstion denyed
		}
		ignore := shouldIgnore(path)
		if !info.IsDir() && !ignore {
			files[path] = info.ModTime()
		}

		if info.IsDir() && ignore {
			return filepath.SkipDir
		}
		return nil
	}

	if err := filepath.Walk(reloadCfg.root, walkFn); err != nil {
		log.Println(err)
	}

	return files
}

func formatSize(file os.FileInfo) string {
	if file.IsDir() {
		return "-"
	}
	size := int(file.Size())
	switch {
	case size > 1024*1024:
		return fmt.Sprintf("%.1fM", float64(size)/1024/1024)
	case size > 1024:
		return fmt.Sprintf("%.1fk", float64(size)/1024)
	default:
		return strconv.Itoa(size)
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

func reloadHandler(w http.ResponseWriter, path string, req *http.Request) {
	switch path {
	case "/js":
		w.Header().Add("Content-Type", "text/javascript")
		w.Header().Add("Cache-Control", "no-cache")
		reloadCfg.reloadJs.Execute(w, req.Host)
	case "/polling":
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reloadCfg.mu.Lock()
		reloadCfg.clients = append(reloadCfg.clients, Client{bufrw, conn})
		reloadCfg.mu.Unlock()
	default:
		showDoc(w, req, nil)
	}
}

func appendReloadHook(w http.ResponseWriter, ctype string, req *http.Request) {
	if strings.HasPrefix(ctype, "text/html") {
		w.Write([]byte("<script src=\"//" + req.Host + "/_d/js\"></script>"))
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
			w.Header().Set("Content-Type", ctype)
		}
		w.WriteHeader(200)
		io.Copy(w, f)
		appendReloadHook(w, ctype, req)
	}
}

// proxy dynamic website, add the reload hook if HTML
func proxyHandler(w http.ResponseWriter, req *http.Request) {
	host := "http://127.0.0.1:" + strconv.Itoa(reloadCfg.proxy)
	url := host + req.URL.String()
	client := &http.Client{}
	if request, err := http.NewRequest(req.Method, url, req.Body); err == nil {
		request.Header.Add("X-Forwarded-For", strings.Split(req.RemoteAddr, ":")[0])
		request.Header.Add("Host", host)
		for k, values := range req.Header {
			for _, v := range values {
				if k != "Host" {
					request.Header.Add(k, v)
				}
			}
		}
		if resp, err := client.Do(request); err == nil {
			for k, values := range resp.Header {
				for _, v := range values {
					// Transfer-Encoding:chunked, for append reload hook
					if k != "Content-Length" {
						w.Header().Add(k, v)
					}
				}
			}
			defer resp.Body.Close()
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			appendReloadHook(w, w.Header().Get("Content-Type"), req)
		} else {
			showDoc(w, req, err) // remote may refuse connection
		}
	} else {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if len(path) > 3 && path[0:3] == "/_d" {
		reloadHandler(w, path[3:], req)
	} else if reloadCfg.proxy == 0 {
		fileHandler(w, path[1:], req)
	} else {
		proxyHandler(w, req)
	}
}

func startMonitorFs() {
	files := getAllFileMeta()
	n := strconv.Itoa(len(files))
	if len(files) > 1000 {
		log.Println("WARN: directory has too many files:", n, "; if CPU usage is high, try tweak -ignores param")
	} else {
		log.Println(n, "files been watched for change, check interval 150ms")
	}
	for {
		time.Sleep(150 * time.Millisecond)
		events := make([]FileEvent, 0)
		tmp := getAllFileMeta()
		for file, mTime := range tmp {
			if oldTime, exits := files[file]; exits {
				if oldTime.Before(mTime) {
					events = append(events, FileEvent{file, MODIFY})
				}
				delete(files, file)
			} else {
				events = append(events, FileEvent{file, ADD})
			}
		}
		for file, _ := range files {
			events = append(events, FileEvent{file, REMOVE})
		}
		files = tmp
		if len(events) > 0 {
			reloadCfg.eventsCh <- events
		}
	}
}

func compilePattens() {
	reloadCfg.mu.Lock()
	defer reloadCfg.mu.Unlock()
	ignores := strings.Split(reloadCfg.ignores, ",")
	// ignore hidden file and emacs generated file
	ignores = append(ignores, `/\.[^/]+`, `/#\.[^/]+`)
	pattens := make([]*regexp.Regexp, 0)
	for _, s := range ignores {
		if len(s) > 0 {
			if p, e := regexp.Compile(s); e == nil {
				pattens = append(pattens, p)
			} else {
				log.Println("ERROR: can not compile to regex", s, e)
			}
		}
	}
	reloadCfg.ignorePattens = pattens
}

func notifyBrowsers() {
	reloadCfg.mu.Lock()
	defer reloadCfg.mu.Unlock()
	for _, c := range reloadCfg.clients {
		defer c.conn.Close()
		reload := "HTTP/1.1 200 OK\r\n"
		reload += "Cache-Control: no-cache\r\nContent-Type: text/javascript\r\n\r\n"
		reload += "location.reload(true);"
		c.buf.Write([]byte(reload))
		c.buf.Flush()
	}
	reloadCfg.clients = make([]Client, 0)
}

func processFsEvents() {
	for {
		events := <-reloadCfg.eventsCh
		command := reloadCfg.command
		if command != "" {
			args := make([]string, len(events)*2)
			for i, e := range events {
				args[2*i] = e.Event
				args[2*i+1] = e.File
			}
			sub := exec.Command(command, args...)
			var out bytes.Buffer
			sub.Stdout = &out
			err := sub.Run()
			if err == nil {
				log.Println("run "+command+" ok; output: ", out.String())
				notifyBrowsers()
			} else {
				log.Println("ERROR running "+command, err)
			}
		} else {
			notifyBrowsers()
		}
	}
}

func main() {
	flag.IntVar(&(reloadCfg.port), "port", 8000, "Which port to listen")
	flag.StringVar(&(reloadCfg.root), "root", ".", "Watched root directory for filesystem events, also the HTTP File Server's root directory")
	flag.StringVar(&(reloadCfg.command), "command", "", "Command to run before reload browser, useful for preprocess, like compile scss. The files been chaneged, along with event type are pass as arguments")
	flag.StringVar(&(reloadCfg.ignores), "ignores", "", "Ignored file pattens, seprated by ',', used to ignore the filesystem events of some files")
	flag.BoolVar(&(reloadCfg.private), "private", false, "Only listen on lookback interface, otherwise listen on all interface")
	flag.IntVar(&(reloadCfg.proxy), "proxy", 0, "Local dynamic site's port number, like 8080, HTTP watcher proxy it, automatically reload browsers when watched directory's file changed")
	flag.Parse()

	if _, e := os.Open(reloadCfg.command); e == nil {
		// turn to abs path if exits
		abs, _ := filepath.Abs(reloadCfg.command)
		reloadCfg.command = abs
	}

	// compile templates
	t, _ := template.New("reloadjs").Parse(RELOAD_JS)
	reloadCfg.reloadJs = t
	t, _ = template.New("dirlist").Parse(DIR_HTML)
	reloadCfg.dirListTmpl = t
	t, _ = template.New("doc").Parse(HELP_HTML)
	reloadCfg.docTmpl = t

	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	compilePattens()
	if e := os.Chdir(reloadCfg.root); e != nil {
		log.Panic(e)
	}
	go startMonitorFs()
	go processFsEvents()

	http.HandleFunc("/", handler)

	int := ":" + strconv.Itoa(reloadCfg.port)
	p := strconv.Itoa(reloadCfg.port)
	mesg := ""
	if reloadCfg.proxy != 0 {
		mesg += "; proxy site http://127.0.0.1:" + strconv.Itoa(reloadCfg.proxy)
	}
	mesg += "; please visit http://127.0.0.1:" + p
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
