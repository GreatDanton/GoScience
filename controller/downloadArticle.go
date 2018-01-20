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
		r.ParseForm()
		doi := template.HTMLEscapeString(r.Form.Get("doi"))
		downloadArticle(w, r, doi)
	}
}

func downloadArticle(w http.ResponseWriter, r *http.Request, doi string) {
	article := parse.Article{}
	err := article.GetPdf(doi)
	if err != nil {
		fmt.Println("################## GetPdf error")
		fmt.Println(err)
		// server returned captcha, display captcha image & relevant template
		if err == parse.ErrCaptchaPresent {
			captcha := article.Captcha
			err = captchaTemplate.Execute(w, captcha)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

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
	pdfName := article.Name
	pdf := article.PdfStream
	// opens up a browser popup for pdf download
	w.Header().Set("Content-Disposition", "attachment; filename="+pdfName)
	http.ServeContent(w, r, pdfName, time.Now(), bytes.NewReader(pdf))
}
