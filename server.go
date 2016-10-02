package main

import (
	"fmt"
  "log"
	"net/http"
  "os"

  // Third Party Libraries
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
  "github.com/julienschmidt/httprouter"
)

var (
    code  = ""
    token = ""
)

// Your credentials should be obtained from the Google
// Developer Console (https://console.developers.google.com).
var oauthCfg = &oauth2.Config{
  ClientID:     os.Getenv("GO_GUESTBOOK_GOOGLE_CLIENT_ID"),
  ClientSecret: os.Getenv("GO_GUESTBOOK_GOOGLE_CLIENT_SECRET"),
  RedirectURL:  "http://localhost:8000/auth/google/callback",
  Scopes: []string{"profile"},
  Endpoint: google.Endpoint,
}

func GoogleHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
  http.Redirect(w,r, oauthCfg.AuthCodeURL("state"), http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  // Uses temporary token provided by Google to get Token. Temporary token is
  // found in the "code" attribute in the request.
  fmt.Println("****** New User ******")
  code := r.FormValue("code")
  token, err := oauthCfg.Exchange(oauth2.NoContext, code)
  if err != nil {
    fmt.Println("Failed Login")
    http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
  } else {
    fmt.Printf("Successful Login! Token Value: %v\n", token)
    http.Redirect(w, r, "/index", http.StatusTemporaryRedirect)
  }
  fmt.Println("******************\n")
}

func loginPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Println("LOAD LOGIN PAGE")
  http.ServeFile(w, r, "src/login/index.html")
}
func indexPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  http.ServeFile(w, r, "src/index/index.html")
}
func errorPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  http.ServeFile(w, r, "src/error/index.html")
}

func main() {
  router := httprouter.New()

  // OAuth
  router.GET("/auth/google/", GoogleHandler)
  router.GET("/auth/google/callback", GoogleCallbackHandler)

  // Routes
  router.GET("/", loginPageHandler)
  router.GET("/index/", indexPageHandler)
  router.GET("/error/", errorPageHandler)

  log.Fatal(http.ListenAndServe(":8000", router))
}
