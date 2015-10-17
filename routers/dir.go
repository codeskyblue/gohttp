package routers

import (
	"fmt"
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

func inspectFileInfo(info os.FileInfo) map[string]interface{} {
	ftype := "file"
	if info.IsDir() {
		ftype = "directory"
	}
	return map[string]interface{}{
		"name":  info.Name(),
		"type":  ftype,
		"size":  info.Size(),
		"mtime": info.ModTime().Unix(),
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
		data = append(data, inspectFileInfo(finfo))
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
			if req.FormValue("download") == "true" {
				w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(abspath))
			}
			http.ServeFile(w, req, abspath)
		}
	}
}
