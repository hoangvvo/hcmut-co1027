package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/hoangvvo/hcmut-co1027/app"
	"github.com/hoangvvo/hcmut-co1027/conf"
	"github.com/julienschmidt/httprouter"
)

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func main() {
	router := httprouter.New()
	router.GET("/", app.Index)
	router.GET("/check", app.CheckHandler)
	router.GET("/check/result", app.CheckResultHandler)
	router.POST("/check", app.CheckCompileHandler)
	router.DELETE("/check/result", app.CheckDeleteHandler)
	router.POST("/check/run", app.CheckTestPostHandler)
	router.GET("/upload", app.Upload)
	router.POST("/upload", app.DoUpload)
	router.DELETE("/upload/:name", app.DeleteUpload)

	if err := os.MkdirAll(conf.CasesDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(conf.ArchiveDir, os.ModePerm); err != nil {
		panic(err)
	}

	router.ServeFiles("/cases/*filepath", http.Dir(conf.CasesDir))
	router.ServeFiles("/case-archives/*filepath", http.Dir(conf.CasesDir))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", gziphandler.GzipHandler(router)); err != nil {
		log.Fatal(err)
	}
}
