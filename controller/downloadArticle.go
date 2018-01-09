package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/greatdanton/goScience/parse"
)

// load template into memory at compile time (better performance than loading it
// each time on function call)
var templateDownload = template.Must(template.ParseFiles("templates/download.html"))

// downloadForm is used for populating fields & displaying error
// messages in download.html template
type downloadForm struct {
	Doi      string
	LabelDoi string
}

// DownloadArticle handles client article download requests
func DownloadArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := templateDownload.Execute(w, nil)
		if err != nil {
			log.Println(err)
		}
	case "POST":
		downloadRequest(w, r)
	}
}

func downloadRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	doi := template.HTMLEscapeString(r.Form.Get("doi"))
	pdf, pdfName, err := parse.GetPdf(doi)
	// inform user about wrong doi
	if err != nil {
		/*
			// TODO: error contains captcha
			if err == parse.ErrCaptchaPresent {
			parsedHTML := string(pdf)
			tmpl, err := template.New(pdfName).Parse(parsedHTML)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				fmt.Println(err)
			}
			return
		} */

		// display error message to the end user
		msg := fmt.Sprintf("%v", err)
		doi := r.Form.Get("doi")
		data := downloadForm{Doi: doi, LabelDoi: msg}
		err := templateDownload.Execute(w, data)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	// opens up a browser popup for pdf download
	w.Header().Set("Content-Disposition", "attachment; filename="+pdfName)
	http.ServeContent(w, r, pdfName, time.Now(), bytes.NewReader(pdf))
}
