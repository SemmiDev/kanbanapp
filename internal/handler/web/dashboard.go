package web

import (
	"embed"
	"net/http"
	"path"
	"text/template"

	"github.com/SemmiDev/kanbanapp/internal/client"
)

type DashboardWeb interface {
	Dashboard(w http.ResponseWriter, r *http.Request)
}

type dashboardWeb struct {
	categoryClient client.CategoryClient
	embed          embed.FS
}

func NewDashboardWeb(catClient client.CategoryClient, embed embed.FS) *dashboardWeb {
	return &dashboardWeb{catClient, embed}
}

func (d *dashboardWeb) Dashboard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	categories, err := d.categoryClient.GetCategories(userId.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dataTemplate = map[string]interface{}{
		"categories": categories,
	}

	var funcMap = template.FuncMap{
		"categoryInc": func(catId int) int {
			// get list categories first, for make sure we get the updated list of categories from database
			categories, err := d.categoryClient.GetCategories(userId.(string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return catId // if error have occured, return the current category
			}

			for i, category := range categories {
				if category.ID == catId {
					// if the current category is the last category, return the current category
					if i == len(categories)-1 {
						return catId
					}
					// if the current category is not the last category, get the next category
					catId = categories[i+1].ID
					break
				}
			}

			return catId
		},

		"categoryDec": func(catId int) int {
			// get list categories first, for make sure we get the updated list of categories from database
			categories, err := d.categoryClient.GetCategories(userId.(string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return catId // if error have occured, return the current category
			}

			for i, category := range categories {
				if category.ID == catId {
					// if the current category is the first category, return the current category
					if i == 0 {
						return catId
					}
					// if the current category is not the first category, get the previous category
					catId = categories[i-1].ID
					break
				}
			}

			return catId
		},
	}

	var header = path.Join("views", "general", "header.html")
	var footer = path.Join("views", "general", "footer.html")
	var filepath = path.Join("views", "main", "dashboard.html")

	tmpl, err := template.New("dashboard.html").Funcs(funcMap).ParseFS(d.embed, header, footer, filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, dataTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
