package routers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Unknwon/macaron"
)

func formatSize(file os.FileInfo) string {
	if file.IsDir() {
		return "-"
	}
	size := file.Size()
	switch {
	case size > 1024*1024:
		return fmt.Sprintf("%.1fM", float64(size)/1024/1024)
	case size > 1024:
		return fmt.Sprintf("%.1fK", float64(size)/1024)
	default:
		return strconv.Itoa(int(size))
	}
	return ""
}

func dirHandler(host, path string, f *os.File, ctx *macaron.Context) {
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
				"href":  href,
				"name":  d.Name(),
				"size":  formatSize(d),
				"mtime": d.ModTime().Format("2006-01-02 15:04:05"),
				"host":  host,
				"path":  filepath.Join(path, d.Name()),
			}
		}
		ctx.HTML(200, "dirlist", map[string]interface{}{
			"dir":   f.Name(),
			"files": files,
		})
	}
}

func NewDirHandler(rootDir string) func(ctx *macaron.Context, req *http.Request, w http.ResponseWriter) {
	return func(ctx *macaron.Context, req *http.Request, w http.ResponseWriter) {
		path := ctx.Params("*")
		if path == "" {
			path = "."
		}
		fullpath := filepath.Join(rootDir, path)
		// log.Println(path)
		file, err := os.Open(fullpath)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}
		defer file.Close()

		finfo, er := file.Stat()
		if er != nil {
			ctx.Error(500, err.Error())
			return
		}
		if finfo.IsDir() {
			dirHandler(req.Host, path, file, ctx)
		} else {
			w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(path))
			http.ServeFile(w, req, filepath.Join(rootDir, path))
		}
	}
}
