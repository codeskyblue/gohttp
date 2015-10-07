package routers

import (
	"net/http"
	"path/filepath"

	"github.com/Unknwon/macaron"
)

func AssetsHandler(r *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
	http.ServeFile(w, r, filepath.Join("public", ctx.Params("*")))
}
