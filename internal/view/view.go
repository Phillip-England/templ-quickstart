package view

import (
	"net/http"
	"path/filepath"
	"snake-scape/internal/component"
	"snake-scape/internal/middleware"
	"snake-scape/internal/template"

	"github.com/a-h/templ"
)

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	filePath := "favicon.ico"
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func Home(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	
	}
	template.Base(
		"SnakeScape - Build RuneScape bots with Python",
		[]templ.Component{
			component.TextAndTitle("I'm a Component!", "I am included as a content item in the Base Template!"),
			component.TextAndTitle("I'm another Component!", "I am also included in the Base Template!"),
		},
	).Render(ctx, w)
}