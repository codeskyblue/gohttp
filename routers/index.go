package routers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/macaron.v1"
)

// record download count
var downloadMap map[string]int64

func init() {
	downloadMap = make(map[string]int64)
}

func dumpDataBackground(dumpFile string, interval time.Duration) {
	for {
		time.Sleep(interval)
		data, _ := json.Marshal(downloadMap)
		ioutil.WriteFile(dumpFile, data, 0644)
	}
}

func NewStaticHandler(root string) interface{} {
	dumpFile := filepath.Join(root, ".gohttp.stat.json")
	data, err := ioutil.ReadFile(dumpFile)
	if err == nil {
		json.Unmarshal(data, &downloadMap)
	}
	go dumpDataBackground(dumpFile, time.Minute)

	return func(ctx *macaron.Context, w http.ResponseWriter, req *http.Request) {
		format := req.FormValue("format")
		if format == "" {
			format = "html"
		}
		relpath := filepath.Clean(req.URL.Path)
		abspath := filepath.Join(root, relpath) //req.URL.Path)
		finfo, err := os.Stat(abspath)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}
		if finfo.IsDir() {
			switch format {
			case "html":
				ctx.HTML(200, "index", nil)
				return
			case "json":
				data, err := listDirectory(root, relpath)
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
			downloadMap[relpath] = downloadMap[relpath] + 1
			http.ServeFile(w, req, abspath)
		}
	}
}

func listDirectory(root, relpath string) (data []interface{}, err error) {
	dir := filepath.Join(root, relpath)
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
		data = append(data, inspectFileInfo(root, relpath, finfo))
	}
	return
}

func deepPath(basedir, name string) string {
	isDir := true
	// loop max 5, incase of for loop not finished
	maxDepth := 5
	for depth := 0; depth <= maxDepth && isDir; depth += 1 {
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

func inspectFileInfo(root, relpath string, info os.FileInfo) map[string]interface{} {
	basedir := filepath.Join(root, relpath)
	name := info.Name()
	if info.IsDir() {
		return map[string]interface{}{
			"name":  deepPath(basedir, name),
			"type":  "directory",
			"size":  info.Size(),
			"mtime": info.ModTime().Unix(),
		}
	} else {
		reqpath := filepath.Join(relpath, name)
		log.Println("inspect", reqpath)
		return map[string]interface{}{
			"name":     info.Name(),
			"type":     "file",
			"size":     info.Size(),
			"mtime":    info.ModTime().Unix(),
			"download": downloadMap[reqpath],
		}
	}

}
