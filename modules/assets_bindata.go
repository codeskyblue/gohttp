// +build bindata

package modules

import (
	"github.com/codeskyblue/gohttp/public"
	"github.com/codeskyblue/gohttp/templates"
	"github.com/go-macaron/bindata"
	"gopkg.in/macaron.v1"
)

var Public = macaron.Static("public",
	macaron.StaticOptions{
		Prefix: "-",
		FileSystem: bindata.Static(bindata.Options{
			Asset:      public.Asset,
			AssetDir:   public.AssetDir,
			AssetNames: public.AssetNames,
			Prefix:     "",
		}),
	})

var Renderer = macaron.Renderer(macaron.RenderOptions{
	TemplateFileSystem: bindata.Templates(bindata.Options{
		Asset:      templates.Asset,
		AssetDir:   templates.AssetDir,
		AssetNames: templates.AssetNames,
		Prefix:     "",
	}),
})
