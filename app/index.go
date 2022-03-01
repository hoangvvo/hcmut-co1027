package app

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tmplHome = template.Must(template.ParseFiles("template/index.html", "template/base.html"))

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := tmplHome.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
