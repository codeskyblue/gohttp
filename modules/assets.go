// +build !bindata

package modules

import "gopkg.in/macaron.v1"

var (
	Public = macaron.Static("public", macaron.StaticOptions{
		Prefix: "-",
	})
	Renderer = macaron.Renderer()
)
