// +build bindata

package modules

import (
	"github.com/Unknwon/macaron"
	"github.com/codeskyblue/file-server/public"
	"github.com/codeskyblue/file-server/templates"
	"github.com/macaron-contrib/bindata"
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
