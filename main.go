package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"strconv"
	"strings"

	"github.com/codeskyblue/gohttp/modules"
	"github.com/codeskyblue/gohttp/routers"
	"github.com/go-macaron/auth"
	"github.com/go-macaron/gzip"
	"gopkg.in/macaron.v1"
)

type Configure struct {
	port     int
	root     string
	private  bool
	httpauth string
	cert     string
	key      string
}

var gcfg = Configure{}

var m *macaron.Macaron

func init() {
	m = macaron.Classic()
	m.Use(modules.Public)
	m.Use(modules.Renderer)
	m.Use(gzip.Gziper())

	flag.IntVar(&gcfg.port, "port", 8000, "Which port to listen")
	flag.StringVar(&gcfg.root, "root", ".", "Watched root directory for filesystem events, also the HTTP File Server's root directory")
	flag.BoolVar(&gcfg.private, "private", false, "Only listen on lookback interface, otherwise listen on all interface")
	flag.StringVar(&gcfg.httpauth, "auth", "", "Basic Authentication (ex: username:password)")
	flag.StringVar(&gcfg.cert, "cert", "", "TLS cert.pem")
	flag.StringVar(&gcfg.key, "key", "", "TLS key.pem")
}

func initRouters() {
	m.Get("/$qrcode", routers.Qrcode)
	m.Get("/*", routers.NewStaticHandler(gcfg.root))
	m.Post("/*", routers.NewUploadHandler(gcfg.root))
	m.Get("/$zip/*", routers.NewZipDownloadHandler(gcfg.root))
	m.Get("/$plist/*", routers.NewPlistHandler(gcfg.root))
	m.Get("/$ipaicon/*", routers.NewIpaIconHandler(gcfg.root))
	m.Get("/$ipa/*", routers.IPAHandler)
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
	// HTTP Basic Authentication
	userpass := strings.SplitN(gcfg.httpauth, ":", 2)
	if len(userpass) == 2 {
		user, pass := userpass[0], userpass[1]
		m.Use(auth.Basic(user, pass))
	}

	m.Get("/-/:rand(.*).hot-update.:ext(.*)", ReloadProxy)
	m.Get("/-/:name(.*).bundle.js", ReloadProxy)
}

func main() {
	flag.Parse()
	initRouters()

	http.Handle("/", m)

	addr := ":" + strconv.Itoa(gcfg.port)
	p := strconv.Itoa(gcfg.port)
	mesg := "; please visit http://127.0.0.1:" + p
	if gcfg.private {
		addr = "localhost" + addr
		log.Printf("listens on 127.0.0.1@" + p + mesg)
	} else {
		log.Printf("listens on 0.0.0.0@" + p + mesg)
	}

	var err error
	if gcfg.key != "" && gcfg.cert != "" {
		err = http.ListenAndServeTLS(addr, gcfg.cert, gcfg.key, nil)
	} else {
		err = http.ListenAndServe(addr, nil)
	}
	log.Fatal(err)
}
