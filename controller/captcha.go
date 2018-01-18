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
		captchaStr := r.Form.Get("answer")
		captchaID := r.Form.Get("id")
		ArticleDoi := r.Form.Get("articleDoi")
		ArticleLink := r.Form.Get("articleLink")

		form := url.Values{
			"answer": {captchaStr},
			"id":     {captchaID},
		}

		body := bytes.NewBufferString(form.Encode())
		fmt.Println("######## ARTICLE LINK", ArticleLink)
		response, err := http.Post(ArticleLink, "application/x-www-form-urlencoded", body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(response)
		// download article
		downloadArticle(w, r, ArticleDoi)
	}
}
