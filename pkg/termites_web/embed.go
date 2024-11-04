package termites_web

import (
	"embed"
	"net/http"
)

//go:embed connect.js
var embeddedJS embed.FS

func EmbeddedJS() http.Handler {
	return http.FileServer(http.FS(embeddedJS))
}
