package app

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/hoangvvo/hcmut-co1027/conf"
	"github.com/hoangvvo/hcmut-co1027/runner"
	"github.com/julienschmidt/httprouter"
)

var tempCheck = template.Must(template.ParseFiles("template/check.html", "template/base.html"))
var tempCheckResult = template.Must(template.ParseFiles("template/check-result.html", "template/base.html"))

type DataCheck struct {
	Suites     []runner.TestSuite
	Error      error
	PrevAnswer string
}

type CheckResult struct {
	Total         int
	Correct       int
	ExecutionTime int64
	Results       []runner.Result
}

type DataCheckResult struct {
	TestSuite runner.TestSuite
	RunDir    string
}

var hashKey = []byte("p_ps}}#C*uZUh,v`Ntk*8LE(LNvD4:Gx") // FIXME: to env
var blockKey = []byte("blockKeyblockKey")
var s = securecookie.New(hashKey, blockKey)

const cookieName = "sid-v0"
const SESSION_PROP = "compile"

func readSess(w http.ResponseWriter, r *http.Request) *runner.CompileResult {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie == nil {
		return nil
	}

	value := make(map[string]string)

	if err = s.Decode(cookieName, cookie.Value, &value); err != nil {
		setSess(w, nil) // invalid session, should just remove
		return nil
	}

	if value["runDir"] == "" {
		return nil
	}

	return &runner.CompileResult{
		RunDir:    value["runDir"],
		SuiteName: value["suiteName"],
	}
}

func setSess(w http.ResponseWriter, compileRes *runner.CompileResult) error {
	if compileRes == nil {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		}
		http.SetCookie(w, cookie)
		return nil
	}
	value := map[string]string{
		"runDir":    compileRes.RunDir,
		"suiteName": compileRes.SuiteName,
	}
	encoded, err := s.Encode(cookieName, value)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}

var testId = 0 // FIXME: concurrency issue

func CheckTestPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	testId += 1
	l := log.New(os.Stdout, "test "+strconv.Itoa(testId)+": ", log.LstdFlags)

	compileRes := readSess(w, r)
	if compileRes == nil {
		sendResponseErr(w, http.StatusBadRequest, errors.New("no session"))
		return
	}

	urlQuery := r.URL.Query()
	if urlQuery.Get("runDir") != compileRes.RunDir {
		sendResponseErr(w, http.StatusBadRequest, errors.New("mismatched runDir"))
		return
	}

	suiteDir := filepath.Join(conf.SuitesDir, compileRes.SuiteName)

	var CaseDirs []string
	if urlQuery.Has("all") {
		l.Println(compileRes.SuiteName + ": test all")
		dirEntries, err := os.ReadDir(suiteDir)
		if err != nil {
			sendResponseErr(w, http.StatusInternalServerError, err)
			return
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				CaseDirs = append(CaseDirs, filepath.Join(suiteDir, entry.Name()))
			}
		}
	} else {
		caseNames := strings.Split(urlQuery.Get("cases"), ",")
		l.Println(compileRes.SuiteName + ": test " + strconv.Itoa(len(caseNames)) + " cases")
		for _, caseName := range caseNames {
			CaseDirs = append(CaseDirs, filepath.Join(suiteDir, caseName))
		}
	}

	executionTimeStart := time.Now()

	results, err := runner.Run(compileRes.RunDir, CaseDirs)

	executionTimeEnd := time.Now()

	if err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	execTime := executionTimeEnd.UnixMilli() - executionTimeStart.UnixMilli()

	l.Println("done in " + strconv.Itoa(int(execTime)) + "ms")

	correctCount := 0
	for _, result := range results {
		if len(result.Error) > 0 {
			if len(result.ResultGot) == 0 {
				result.ResultGot = result.Error
			}
			correctCount += 1
		}
	}
	total := len(results)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age=31536000")

	err = json.NewEncoder(w).Encode(CheckResult{
		Total:         total,
		Correct:       correctCount,
		ExecutionTime: execTime,
		Results:       results,
	})

	if err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
	}
}

func CheckDeleteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	compileRes := readSess(w, r)
	if compileRes != nil {
		err := runner.DeleteCompiled(compileRes.RunDir)
		if err != nil {
			sendResponseErr(w, http.StatusInternalServerError, err)
			return
		}
		err = setSess(w, nil)
		if err != nil {
			sendResponseErr(w, http.StatusInternalServerError, err)
			return
		}
	}
	sendResponseSuccess(w)
}

var compileId = 0 // FIXME: concurrency issue

func CheckCompileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess := readSess(w, r)
	if sess != nil {
		// remove previous sess
		runner.DeleteCompiled(sess.RunDir)
		setSess(w, nil)
	}

	answer := strings.TrimSpace(r.FormValue("answer"))
	answerFileName := strings.TrimSpace(r.FormValue("answer-filename"))
	suiteName := strings.TrimSpace(r.FormValue("suite"))

	if answer == "" || answerFileName == "" || suiteName == "" {
		sendResponseErr(w, http.StatusBadRequest, errors.New("invalid input"))
		return
	}

	compileId += 1
	l := log.New(os.Stdout, "compile "+strconv.Itoa(testId)+": ", log.LstdFlags)

	executionTimeStart := time.Now()

	compiledRes, err := runner.Compile(answer, answerFileName, suiteName)

	executionTimeEnd := time.Now()

	// error happen, just render normally with error message
	if err != nil {
		testSuites, errGet := runner.GetSuites()
		if errGet != nil {
			sendResponseErr(w, http.StatusInternalServerError, err)
			return
		}
		if err := tempCheck.Execute(w, DataCheck{
			Suites: testSuites,
			Error:  err,
		}); err != nil {
			sendResponseErr(w, http.StatusInternalServerError, err)
			return
		}
		return
	}

	execTime := executionTimeEnd.UnixMilli() - executionTimeStart.UnixMilli()

	l.Println("done in " + strconv.Itoa(int(execTime)) + "ms")

	// set session
	if err = setSess(w, compiledRes); err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	// redirect to render with result page
	http.Redirect(w, r, conf.AppURI+"/check/result", http.StatusSeeOther)
}

func CheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	testSuites, err := runner.GetSuites()
	if err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
		return
	}
	if err := tempCheck.Execute(w, DataCheck{
		Suites: testSuites,
	}); err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
		return
	}
}

func CheckResultHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess := readSess(w, r)

	if sess == nil {
		http.Redirect(w, r, conf.AppURI+"/check", http.StatusSeeOther)
		return
	}

	ts, err := runner.GetSuite(sess.SuiteName)
	if err != nil || ts == nil {
		// error getting suite redirect while deleting this cookie
		setSess(w, nil)
		http.Redirect(w, r, conf.AppURI+"/check", http.StatusSeeOther)
		return
	}

	if err := tempCheckResult.Execute(w, DataCheckResult{
		TestSuite: *ts,
		RunDir:    sess.RunDir,
	}); err != nil {
		sendResponseErr(w, http.StatusInternalServerError, err)
	}
}
