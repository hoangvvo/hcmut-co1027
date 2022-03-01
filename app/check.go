package app

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/hoangvvo/hcmut-co1027/testcase"
	"github.com/julienschmidt/httprouter"
)

var tempCheck = template.Must(template.ParseFiles("template/check.html", "template/base.html"))

type DataCheck struct {
	Suites     []testcase.TestSuite
	Error      error
	Results    []testcase.Result
	ResultStat *ResultStat
	PrevAnswer *string
}

type ResultStat struct {
	Total         int
	Correct       int
	Percentage    string
	ExecutionTime int64
	TestSuite     string
}

func DoCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	testSuites, err := testcase.GetSuites()
	if err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
	}

	answer := r.FormValue("answer")
	suiteName := r.FormValue("suite")

	executionTimeStart := time.Now()

	results, err := testcase.RunSuite(suiteName, answer)

	executionTimeEnd := time.Now()

	var incorrectResults []testcase.Result
	for _, result := range results {
		if len(result.Error) > 0 {
			if len(result.MyOutput) == 0 {
				result.MyOutput = result.Error
			}
			incorrectResults = append(incorrectResults, result)
		}
	}
	total := len(results)
	correctCount := len(results) - len(incorrectResults)
	percentage := fmt.Sprintf("%.2f", float32(correctCount)/float32(total)*100)

	if err := tempCheck.Execute(w, DataCheck{
		Suites:  testSuites,
		Error:   err,
		Results: incorrectResults,
		ResultStat: &ResultStat{
			Total:         total,
			Correct:       correctCount,
			Percentage:    percentage,
			ExecutionTime: executionTimeEnd.UnixMilli() - executionTimeStart.UnixMilli(),
			TestSuite:     suiteName,
		},
		PrevAnswer: &answer,
	}); err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
	}
}

func Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	testSuites, err := testcase.GetSuites()
	if err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
	}
	if err := tempCheck.Execute(w, DataCheck{
		Suites: testSuites,
	}); err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
	}
}
