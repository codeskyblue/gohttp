package routers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/macaron.v1"
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

func deepPath(basedir, name string) string {
	isDir := true
	for isDir {
		finfos, err := ioutil.ReadDir(filepath.Join(basedir, name))
		if err != nil || len(finfos) != 1 {
			return name
		}
		if finfos[0].IsDir() {
			name = filepath.Join(name, finfos[0].Name())
		} else {
			break
		}
	}
	return name
}

func inspectFileInfo(basedir string, info os.FileInfo) map[string]interface{} {
	name := info.Name()
	if info.IsDir() {
		return map[string]interface{}{
			"name":  deepPath(basedir, name),
			"type":  "directory",
			"size":  info.Size(),
			"mtime": info.ModTime().Unix(),
		}
	} else {
		return map[string]interface{}{
			"name":  info.Name(),
			"type":  "file",
			"size":  info.Size(),
			"mtime": info.ModTime().Unix(),
		}
	}

}

func listDirectory(dir string) (data []interface{}, err error) {
	file, err := os.Open(dir)
	if err != nil {
		return
	}
	defer file.Close()
	files, err := file.Readdir(-1)
	if err != nil {
		return
	}
	data = make([]interface{}, 0, len(files))
	for _, finfo := range files {
		data = append(data, inspectFileInfo(dir, finfo))
	}
	return
}

func NewStaticHandler(root string) interface{} {
	return func(ctx *macaron.Context, w http.ResponseWriter, req *http.Request) {
		format := req.FormValue("format")
		if format == "" {
			format = "html"
		}
		abspath := filepath.Join(root, req.URL.Path)
		finfo, err := os.Stat(abspath)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}
		if finfo.IsDir() {
			switch format {
			case "html":
				ctx.HTML(200, "dirlist", nil)
				return
			case "json":
				data, err := listDirectory(abspath)
				if err != nil {
					ctx.Error(500, err.Error())
					return
				}
				ctx.JSON(200, data)
			}
		} else {
			if req.FormValue("preview") == "true" {
				ctx.HTML(200, "preview", nil)
				return
			}
			if req.FormValue("download") == "true" {
				w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(abspath))
			}
			http.ServeFile(w, req, abspath)
		}
	}
}
