package app

import (
	"html/template"
	"net/http"

	"github.com/hoangvvo/hcmut-co1027/testcase"
	"github.com/julienschmidt/httprouter"
)

func DoUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(20 << 20)
	file, handler, err := r.FormFile("suite")
	if err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
		return
	}
	defer file.Close()
	err = testcase.AddSuite(handler.Filename, file)
	if err != nil {
		sendResponseErr(w, http.StatusBadRequest, err)
		return
	}
	sendResponseSuccess(w)
}

func DeleteUpload(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	suiteName := params.ByName("name")
	if err := testcase.DeleteSuite(suiteName); err != nil {
		sendResponseErr(w, http.StatusBadRequest, err)
		return
	}
	sendResponseSuccess(w)
}

var tempUpload = template.Must(template.ParseFiles("template/upload.html", "template/base.html"))

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	testSuites, err := testcase.GetSuites()
	if err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
		return
	}
	if err := tempUpload.Execute(w, DataCheck{
		Suites: testSuites,
	}); err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
	}
}
