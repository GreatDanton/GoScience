package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/greatdanton/goScience/global"
	"golang.org/x/crypto/bcrypt"
)

var templateLogin = template.Must(template.ParseFiles("templates/login.html"))

type loginForm struct {
	Password   string
	ErrorLabel string
}

// Login takes care of handling user login and displaying error
// message on wrong password input
func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderLogin(w, r, loginForm{})
	case "POST":
		userLogin(w, r)
	}
}

// render login template and display possible errors via loginForm struct
func renderLogin(w http.ResponseWriter, r *http.Request, form loginForm) {
	err := templateLogin.Execute(w, form)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	password := r.Form.Get("password")
	if password != global.PASSWORD {
		data := loginForm{}
		data.Password = password
		data.ErrorLabel = "Wrong password"
		renderLogin(w, r, data)
		return
	}

	// password is okay, create cookie with hashed password
	// explicitly creating cookie from user inputted password
	cookie, err := createCookie(password)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createCookie(password string) (http.Cookie, error) {
	// expires in one week
	expiration := time.Now().Add(7 * 24 * time.Hour)

	passHash, err := hashPassword(password)
	cookie := http.Cookie{Name: "GoScience", Value: passHash, Expires: expiration, Path: "/", HttpOnly: true}
	return cookie, err
}

func hashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	pass := string(hashedPassword)
	return pass, nil
}
