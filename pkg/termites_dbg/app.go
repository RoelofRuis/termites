package termites_dbg

import (
	"embed"
	"io/fs"
)

//go:embed all:app
var embeddedApp embed.FS

func App() (fs.FS, error) {
	return fs.Sub(embeddedApp, "app")
}
