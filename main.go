package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/goScience/controller"
	"github.com/greatdanton/goScience/global"
)

// Configuration struct created for reading config from file
type Configuration struct {
	Port      string
	Password  string
	ScihubURL string
}

// main function
func main() {
	config, err := ReadConfiguration()
	if err != nil {
		fmt.Println(err)
		return
	}
	PORT := config.Port
	global.PASSWORD = config.Password
	global.ScihubURL = config.ScihubURL

	// handling download section
	http.HandleFunc("/", authMiddleware(controller.DownloadArticle))
	http.HandleFunc("/login", loginMiddleware(controller.Login))

	// serving css & public stuff
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// start webserver
	log.Print("Started server on http://127.0.0.1:" + PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// authMiddleware checks if user is already authenticated. If the user is
// not authenticated it sends him to /login otherwise he is able to
// access downloading part of the application
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if cookie with hashed password exist
		cookie, err := r.Cookie("GoScience")
		if err != nil { // cookie does not exist
			fmt.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// check if password in cookie is the same as server set password
		passHash := cookie.Value
		err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(global.PASSWORD))
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// password is correct, serve the request
		next.ServeHTTP(w, r)
	})
}

// loginMiddleware checks if user is already authenticated (and redirects him/her to main download page).
func loginMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("GoScience")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		passHash := cookie.Value
		err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(global.PASSWORD))
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		// password is correct, just redirect user to download page -> "/"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
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

	// check if scihub url is present in configuration
	if len(config.ScihubURL) < 1 {
		return Configuration{}, fmt.Errorf("ScihubURL is not present in configuration")
	}

	return config, nil
}
