package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	//https://github.com/golang/oauth2
	"golang.org/x/oauth2"
)

//https://godoc.org/golang.org/x/oauth2/google
var googleEndpotin = oauth2.Endpoint{
	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://accounts.google.com/o/oauth2/token",
}

//https://godoc.org/golang.org/x/oauth2#example-Config
var googleOauthConfig = &oauth2.Config{
	ClientID:     "1013783811015-g2ajlk3sf78atgjl3e8p2c1ov2beg85p.apps.googleusercontent.com",
	ClientSecret: "JAFkDYpCMKL6BVzZM_BUVOiq",
	RedirectURL:  "http://localhost:8080/googlecallback",
	//https://developers.google.com/identity/protocols/oauth2/scopes
	Scopes:   []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: googleEndpotin,
}

//https://godoc.org/golang.org/x/oauth2/facebook
var facebookEndpoint = oauth2.Endpoint{
	AuthURL:  "https://www.facebook.com/v3.2/dialog/oauth",
	TokenURL: "https://graph.facebook.com/v3.2/oauth/access_token",
}

var facebookOauthConfig = &oauth2.Config{
	ClientID:     "1707037259472125",
	ClientSecret: "81361d8fadf9fb31c5b74670ad92cf71",
	RedirectURL:  "http://localhost:8080/facebookcallback",
	//https://developers.facebook.com/docs/facebook-login/permissions/#reference-email
	Scopes:   []string{"public_profile", "email"},
	Endpoint: facebookEndpoint,
}

func login(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Println(err)
	}

	t.Execute(w, nil)
}

func googlelogin(w http.ResponseWriter, r *http.Request) {

	url := googleOauthConfig.AuthCodeURL("radom")
	fmt.Println(url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func googleCallback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	fmt.Println(state)

	//Authorization Grant code
	code := r.FormValue("code")
	fmt.Println(code)

	//Access token
	token, _ := googleOauthConfig.Exchange(oauth2.NoContext, code)
	fmt.Println(token)

	response, _ := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	// defer somelike stack LIFO
	defer response.Body.Close()

	contents, _ := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)

}

func facebooklogin(w http.ResponseWriter, r *http.Request) {

	url := facebookOauthConfig.AuthCodeURL("radom")
	fmt.Println(url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func facebookcallback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	fmt.Println(state)

	//Authorization Grant code
	code := r.FormValue("code")
	fmt.Println(code)

	//Access token
	token, _ := facebookOauthConfig.Exchange(oauth2.NoContext, code)
	fmt.Println(token)

	//https://developers.facebook.com/docs/graph-api/using-graph-api/
	response, _ := http.Get("https://graph.facebook.com/me?access_token=" + token.AccessToken)
	// defer somelike stack LIFO
	defer response.Body.Close()

	contents, _ := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)
}

func main() {

	http.HandleFunc("/", login)
	http.HandleFunc("/googlelogin", googlelogin)
	http.HandleFunc("/googlecallback", googleCallback)
	http.HandleFunc("/facebooklogin", facebooklogin)
	http.HandleFunc("/facebookcallback", facebookcallback)

	fmt.Println("服務器即將開啓，訪問地址 http://localhost:8080")

	//err := http.ListenAndServe(":0", nil)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println("服務器開啓錯誤: ", err)
	}
}
