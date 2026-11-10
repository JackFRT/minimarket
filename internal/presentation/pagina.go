package presentation

import "embed"

//go:embed pagina.html
var home string

//go:embed all:web
var webFS embed.FS
