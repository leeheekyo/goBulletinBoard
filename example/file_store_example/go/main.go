package main

import (
    "net/http"
    "fmt"
    "os"
    "io"
    "log"
    "strconv"
    "time"
    "text/template"
    "crypto/md5"
)

// main logic 
func index(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("index call\n")
	http.ServeFile(w, r, "html/index.html")
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("upload call\n")
       fmt.Println("method:", r.Method)
       if r.Method == "GET" {
           crutime := time.Now().Unix()
           h := md5.New()
           io.WriteString(h, strconv.FormatInt(crutime, 10))
           token := fmt.Sprintf("%x", h.Sum(nil))

           t, _ := template.ParseFiles("upload.gtpl")
           t.Execute(w, token)
       } else {
           r.ParseMultipartForm(32 << 20)
           file, handler, err := r.FormFile("uploadfile")
           if err != nil {
               fmt.Println(err)
               return
           }
           defer file.Close()
           fmt.Fprintf(w, "%v", handler.Header)
           f, err := os.OpenFile("test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
           if err != nil {
               fmt.Println(err)
               return
           }
           defer f.Close()
           io.Copy(f, file)
       }
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/upload", upload)

    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
