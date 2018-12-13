package main
  
import (
    "fmt"
    "github.com/gorilla/securecookie"
    "net/http"
)
  
  
var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))
  
func getSession(request *http.Request) (yourName string) {
    if cookie, err := request.Cookie("your-name"); err == nil {
        cookieValue := make(map[string]string)
        if err = cookieHandler.Decode("your-name", cookie.Value, &cookieValue); err == nil {
            yourName = cookieValue["your-name"]
        }
    }
    return yourName
}
  
func setSession(yourName string, response http.ResponseWriter) {
    value := map[string]string{
        "your-name": yourName,
    }
    if encoded, err := cookieHandler.Encode("your-name", value); err == nil {
        cookie := &http.Cookie{
            Name:  "your-name",
            Value: encoded,
            Path:  "/",
            MaxAge: 3600,
        }
        http.SetCookie(response, cookie)
    }
}
  
func clearSession(response http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "your-name",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(response, cookie)
}
  
  
func setSessionHandler(response http.ResponseWriter, request *http.Request) {
    name := request.FormValue("name")    
    redirectTarget := "/"
    if name != "" {
        setSession(name, response)
        redirectTarget = "/page1"
    }
    http.Redirect(response, request, redirectTarget, 302)
}
  
 
  
func clearSessionHandler(response http.ResponseWriter, request *http.Request) {
    clearSession(response)
    http.Redirect(response, request, "/", 302)
}
  
 
  
const indexPage = `
<h1>Sesssion Test</h1>
<form action="/start" method="post">
    <label for="name">Your Name</label>
    <input name="name" id="name" type="text">
    <button type="submit">Set Sesssion</button>
</form>
`
  
func indexPageHandler(response http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(response, indexPage)
}
  
  
const frontPage = `
<h1>Check Sesssion</h1>
<hr>
<small>Your Name: %s</small>
<form action="/clear" method="post">
    <button type="submit">Clear Sesssion</button>
</form>
`
  
func sessionPageHandler(response http.ResponseWriter, request *http.Request) {
    yourName := getSession(request)
    if yourName != "" {
        fmt.Fprintf(response, frontPage, yourName)
    } else {
        http.Redirect(response, request, "/", 302)
    }
}
  
  
func main() {
    http.HandleFunc("/", indexPageHandler)
    http.HandleFunc("/page1", sessionPageHandler)
  
    http.HandleFunc("/start", setSessionHandler)//.Methods("POST")
    http.HandleFunc("/clear", clearSessionHandler)//.Methods("POST")
  
    //http.Handle("/", router)
    http.ListenAndServe(":8080", nil)
}

