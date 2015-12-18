package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/codeskyblue/gohttp/modules"
	"github.com/codeskyblue/gohttp/routers"
	"github.com/go-macaron/auth"
	"github.com/go-macaron/gzip"
	//"github.com/goftp/posixfs-driver"
	//goftp "github.com/goftp/server"
	"gopkg.in/macaron.v1"
)

const VERSION = "0.1.2"

type Configure struct {
	port     int
	root     string
	private  bool
	httpauth string
	cert     string
	key      string
	gzip     bool
	ftp      bool
	ftpPort  int
	ftpAuth  string
	upload   bool
	zipable  bool
}

var gcfg = Configure{}

var m *macaron.Macaron

func init() {
	m = macaron.Classic()
	m.Use(modules.Public)

	if _, err := os.Stat("templates"); err == nil {
		m.Use(macaron.Renderer())
	} else {
		m.Use(modules.Renderer)
	}

	kingpin.HelpFlag.Short('h')
	kingpin.Version(VERSION)
	kingpin.Flag("port", "Port to listen").Default("8000").IntVar(&gcfg.port)
	kingpin.Flag("root", "File root directory").Default(".").StringVar(&gcfg.root)
	kingpin.Flag("private", "Only listen on loopback address").BoolVar(&gcfg.private)
	kingpin.Flag("httpauth", "HTTP basic auth (ex: user:pass)").Default("").StringVar(&gcfg.httpauth)
	kingpin.Flag("cert", "TLS cert.pem").StringVar(&gcfg.cert)
	kingpin.Flag("key", "TLS key.pem").StringVar(&gcfg.key)
	kingpin.Flag("gzip", "Enable Gzip support").BoolVar(&gcfg.gzip)
	//kingpin.Flag("ftp", "Enable FTP support").BoolVar(&gcfg.ftp)
	//kingpin.Flag("ftp-port", "FTP listen port").Default("2121").IntVar(&gcfg.ftpPort)
	//kingpin.Flag("ftp-auth", "FTP auth (ex: user:pass)").Default("admin:123456").StringVar(&gcfg.ftpAuth)
	kingpin.Flag("upload", "Enable upload support").BoolVar(&gcfg.upload)
	kingpin.Flag("zipable", "Enable archieve folder into zip").BoolVar(&gcfg.zipable)
}

func initRouters() {
	m.Get("/*", routers.NewStaticHandler(routers.IndexOptions{
		Root:    gcfg.root,
		Upload:  gcfg.upload,
		Zipable: gcfg.zipable,
	}))
	m.Get("/$qrcode", routers.Qrcode)
	m.Get("/$plist/*", routers.NewPlistHandler(gcfg.root))
	m.Get("/$ipaicon/*", routers.NewIpaIconHandler(gcfg.root))
	m.Get("/$ipa/*", routers.IPAHandler)
	if gcfg.upload {
		m.Post("/*", routers.NewUploadHandler(gcfg.root))
	}
	if gcfg.zipable {
		m.Get("/$zip/*", routers.NewZipDownloadHandler(gcfg.root))
	}
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
	if gcfg.gzip {
		m.Use(gzip.Gziper())
	}
	m.Get("/-/:rand(.*).hot-update.:ext(.*)", ReloadProxy)
	m.Get("/-/:name(.*).bundle.js", ReloadProxy)
}

type FTPAuth struct {
	Username string
	Password string
}

func (this *FTPAuth) CheckPasswd(user, pass string) (bool, error) {
	ok := (this.Username == user && this.Password == pass)
	return ok, nil
}

func main() {
	kingpin.Parse()
	initRouters()

	http.Handle("/", m)

	addr := ":" + strconv.Itoa(gcfg.port)
	p := strconv.Itoa(gcfg.port)
	//mesg := "; please visit http://127.0.0.1:" + p
	if gcfg.private {
		addr = "localhost" + addr
		log.Printf("listens on 127.0.0.1@" + p) // + mesg)
	} else {
		log.Printf("listens on 0.0.0.0@" + p) // + mesg)
	}

	/*
		if gcfg.ftp {
			//log.Println("Enable FTP")
			auths := strings.SplitN(gcfg.ftpAuth, ":", 2)
			if len(auths) != 2 {
				log.Fatal("ftp auth format error")
			}
			auth := FTPAuth{auths[0], auths[1]}
			ftpserv := goftp.NewServer(&goftp.ServerOpts{
				Port:    gcfg.ftpPort,
				Factory: posixfs.NewPosixFSFactory(gcfg.root),
				Auth:    &auth,
			})
			go ftpserv.ListenAndServe()
		}
	*/

	var err error
	if gcfg.key != "" && gcfg.cert != "" {
		err = http.ListenAndServeTLS(addr, gcfg.cert, gcfg.key, nil)
	} else {
		err = http.ListenAndServe(addr, nil)
	}
	log.Fatal(err)
}
