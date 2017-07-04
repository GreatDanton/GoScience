package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/greatdanton/goScience/parse"
)

// PASSWORD is global variable, set by ReadConfiguration function.
// It is used to prevent bots wasting our bandwith.
var PASSWORD string

// load template into memory at compile time (better performance than loading it
// each time on function call)
var templateDownload = template.Must(template.ParseFiles("templates/download.html"))

// downloadForm is used for populating fields & displaying error
// messages in download.html template
type downloadForm struct {
	Password  string
	LabelPass string
	Doi       string
	LabelDoi  string
	Token     string
}

// Configuration struct created for reading config from file
type Configuration struct {
	Port     string
	Password string
}

// ReadConfiguration reads from "conf.json" file and returns Configuration struct
// which is used in main func
func ReadConfiguration() (Configuration, error) {
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		fmt.Println("Please add conf.json file")
		return Configuration{}, err
	}

	config := Configuration{}
	if err := json.Unmarshal(data, &config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}

//
// TODO: move function to separate package? (views?)
// downloads article with doi set with doi field parameters in download.html template
// and distribute it to the client
func downloadArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// TODO: fix token part
		// I am not entirely sure what to do with token? According to the book
		//(https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.4.html)
		// implementing it should take care of duplicate submissions
		// Do I even need it?
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t := templateDownload
		data := downloadForm{Token: token}

		err := t.Execute(w, data)
		if err != nil {
			log.Print(err)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		token := r.Form.Get("token")
		// TODO: What do I do with token? Return error if no token is present?
		if token == "" {
			return
		}
		// safe password field escape
		password := template.HTMLEscapeString(r.Form.Get("password"))
		// inform user if password is not correct
		if password != PASSWORD {
			data := parseFormDownload(r, "PASS", "Wrong Password")
			t := templateDownload
			err := t.Execute(w, data)
			if err != nil {
				log.Print(err)
			}
			return
		}

		doi := template.HTMLEscapeString(r.Form.Get("doi"))
		pdf, pdfName, err := parse.GetPdf(doi)
		// inform user about wrong doi
		if err != nil {
			label := "DOI"
			msg := "Article with this doi does not exist"
			error := fmt.Sprintf("%v", err)

			if strings.Contains(error, "502") { // check for Bad Gateway 502
				label = "DOI"
				msg = "Scihub servers are over capacity, try again later"
			}

			data := parseFormDownload(r, label, msg)
			t := templateDownload
			err := t.Execute(w, data)
			if err != nil {
				fmt.Println(err)
				log.Print(err)
			}
			return
		}
		// opens up a browser popup for pdf download
		w.Header().Set("Content-Disposition", "attachment; filename="+pdfName)

		// display pdf with built in pdf viewer
		//(TODO: delete this line, leaving it for now in case I need it later)
		//w.Header().Set("Content-disposition", "inline; filename=article.pdf")
		http.ServeContent(w, r, pdfName, time.Now(), bytes.NewReader(pdf))

	}
}

//
// parses download form, and reurns template with data filled in
// label -> which info label text should be changed
// info -> text for info label
func parseFormDownload(r *http.Request, label string, info string) downloadForm {
	r.ParseForm()
	data := downloadForm{}
	data.Doi = r.Form.Get("doi")
	data.Password = r.Form.Get("password")
	data.Token = r.Form.Get("token")

	if label == "PASS" {
		data.LabelPass = info
	} else if label == "DOI" {
		data.LabelDoi = info
	}

	return data
}

// main function
func main() {
	config, err := ReadConfiguration()
	if err != nil {
		fmt.Println(err)
		return
	}
	PORT := config.Port
	PASSWORD = config.Password

	log.Print("Started server on http://127.0.0.1:" + PORT)

	// handling download section
	http.HandleFunc("/", downloadArticle)
	// serving css & public stuff
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// start webserver
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
