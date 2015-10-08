// +build !bindata

package modules

import "github.com/Unknwon/macaron"

var (
	Public = macaron.Static("public", macaron.StaticOptions{
		Prefix: "-",
	})
	Renderer = macaron.Renderer()
)
