package web

import (
	"embed"
	"net/http"
	"path"
	"text/template"
)

type HomeWeb interface {
	Index(w http.ResponseWriter, r *http.Request)
}

type homeWeb struct {
	embed embed.FS
}

func NewHomeWeb(embed embed.FS) *homeWeb {
	return &homeWeb{embed}
}

func (h *homeWeb) Index(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "main", "index.html")
	var navbar = path.Join("views", "general", "navbar.html")
	var header = path.Join("views", "general", "header.html")
	var footer = path.Join("views", "general", "footer.html")

	var tmpl = template.Must(template.ParseFS(h.embed, filepath, header, navbar, footer))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
