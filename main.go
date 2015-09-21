package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Unknwon/macaron"
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
	m.Use(macaron.Renderer())
}

func main() {
	flag.IntVar(&gcfg.port, "port", 8000, "Which port to listen")
	flag.StringVar(&gcfg.root, "root", ".", "Watched root directory for filesystem events, also the HTTP File Server's root directory")
	flag.BoolVar(&gcfg.private, "private", false, "Only listen on lookback interface, otherwise listen on all interface")
	flag.Parse()

	m.Get("/_qr", routers.Qrcode)

	m.Get("/_/*", func(r *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
		http.ServeFile(w, r, filepath.Join("public", ctx.Params("*")))
	})

	// Handle Upload file
	m.Post("/*", func(req *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
		err := req.ParseMultipartForm(100 << 20) // max memory 100M
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Hi", err)
			return
		}
		log.Println(ctx.Params("*"))
		path := strings.Replace(ctx.Params("*"), "..", "", -1)
		dirpath := filepath.Join(gcfg.root, path)
		for _, mfile := range req.MultipartForm.File["file"] {
			file, err := mfile.Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dst, err := os.Create(filepath.Join(dirpath, mfile.Filename)) // BUG(ssx): There is a leak here
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		ctx.JSON(200, map[string]string{
			"message": "Upload success",
		})
	})

	m.Get("/*", routers.NewDirHandler(gcfg.root))

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
