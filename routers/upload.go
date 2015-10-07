package routers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/macaron"
)

// Handle Upload file
func NewUploadHandler(rootDir string) func(req *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
	return func(req *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
		err := req.ParseMultipartForm(100 << 20) // max memory 100M
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Hi", err)
			return
		}
		log.Println(ctx.Params("*"))
		path := strings.Replace(ctx.Params("*"), "..", "", -1)
		dirpath := filepath.Join(rootDir, path)
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
	}
}
