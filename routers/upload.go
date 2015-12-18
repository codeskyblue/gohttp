package routers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/macaron.v1"
)

// Handle Upload file
func NewUploadHandler(rootDir string) func(req *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
	return func(req *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
		err := req.ParseMultipartForm(1024 << 20) // max memory 100M
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(req.MultipartForm.File["file"]) == 0 {
			http.Error(w, "Need multipart file", http.StatusInternalServerError)
			return
		}

		path := strings.Replace(ctx.Params("*"), "..", "", -1)
		dirpath := filepath.Join(rootDir, path)

		name := req.FormValue("name")
		version := req.FormValue("version")
		if name != "" && version != "" {
			base := filepath.Join(dirpath, name)
			dirpath = filepath.Join(base, name+"-"+version)
			os.MkdirAll(dirpath, 0755)

			//symlinkPath := filepath.Join(base, name+"-latest")
			//os.Symlink(name+"-"+version, symlinkPath)
		}
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
				log.Println("Handle upload file:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"message": "Upload success",
		})
	}
}
