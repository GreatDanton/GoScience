package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

var captchaTemplate = template.Must(template.ParseFiles("templates/captchaForm.html"))

// Captcha handles captcha part of the article download process
func Captcha(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		answer := r.Form.Get("answer")
		captchaID := r.Form.Get("id")
		ArticleDoi := r.Form.Get("articleDoi")
		ArticleURL := r.Form.Get("articleURL")

		form := url.Values{
			"answer": {answer},
			"id":     {captchaID},
		}

		// post captcha message to scihub servers
		body := bytes.NewBufferString(form.Encode())
		response, err := http.Post(ArticleURL, "application/x-www-form-urlencoded", body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// TODO: if status code != 200, display server message to the client
		if response.StatusCode != http.StatusOK {
			fmt.Println(response)
		}

		// download article
		downloadArticle(w, r, ArticleDoi)
	}
}
