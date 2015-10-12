package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"strconv"

	"github.com/Unknwon/macaron"
	"github.com/codeskyblue/file-server/modules"
	"github.com/codeskyblue/file-server/routers"
)

type Configure struct {
	port    int
	root    string
	private bool
}

var gcfg = Configure{}

var m *macaron.Macaron

func init() {
	m = macaron.Classic()
	m.Use(modules.Public)
	m.Use(modules.Renderer)

	flag.IntVar(&gcfg.port, "port", 8000, "Which port to listen")
	flag.StringVar(&gcfg.root, "root", ".", "Watched root directory for filesystem events, also the HTTP File Server's root directory")
	flag.BoolVar(&gcfg.private, "private", false, "Only listen on lookback interface, otherwise listen on all interface")
}

func initRouters() {
	m.Get("/_qr", routers.Qrcode)
	m.Get("/*", routers.NewStaticHandler(gcfg.root))
	m.Post("/*", routers.NewUploadHandler(gcfg.root))
	ReloadProxy := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Debug, Hot reload", r.Host)
		resp, err := http.Get("http://localhost:3000" + r.RequestURI)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	}
	m.Get("/-/:rand(.*).hot-update.:ext(.*)", ReloadProxy)
	m.Get("/-/bundle.js", ReloadProxy)
}

func main() {
	flag.Parse()
	initRouters()

	http.Handle("/", m)

	int := ":" + strconv.Itoa(gcfg.port)
	p := strconv.Itoa(gcfg.port)
	mesg := "; please visit http://127.0.0.1:" + p
	if gcfg.private {
		int = "localhost" + int
		log.Printf("listens on 127.0.0.1@" + p + mesg)
	} else {
		log.Printf("listens on 0.0.0.0@" + p + mesg)
	}
	if err := http.ListenAndServe(int, nil); err != nil {
		log.Fatal(err)
	}
}
