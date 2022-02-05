package handler

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

func (h *handler) err(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusInternalServerError)

	queryValues := r.URL.Query()

	generateHTML(w, r, queryValues.Get("msg"), "layout", "error")
}

func errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println(err)
	url := "/err?msg=" + err.Error()
	http.Redirect(w, r, url, http.StatusFound)
}

func generateHTML(w http.ResponseWriter, r *http.Request, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("/opt/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		errorResponse(w, r, err)
	}
}
