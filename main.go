package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	router.GET("/check", app.Check)
	router.POST("/check", app.DoCheck)
	router.GET("/upload", app.Upload)
	router.POST("/upload", app.DoUpload)
	router.DELETE("/upload/:name", app.DeleteUpload)

	err := os.MkdirAll(conf.CASEDIR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	router.ServeFiles("/tests/*filepath", http.Dir(conf.CASEDIR))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
