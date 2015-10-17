package routers

import (
	"net/http"
	"path/filepath"

	"gopkg.in/macaron.v1"
)

func AssetsHandler(r *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
	http.ServeFile(w, r, filepath.Join("public", ctx.Params("*")))
}
