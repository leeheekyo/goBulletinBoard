package main

import (
    "fmt"
    "log"
    "strconv"
    "net/http"
    "html/template"
    "io/ioutil"
//    "regexp"
    "github.com/ziutek/mymysql/mysql"    //db
    _ "github.com/ziutek/mymysql/native" //db
    "github.com/gorilla/securecookie"    //session
)

const (
	user = "root"
	pass = "kyo"
	dbname = "test"
)

const sepa = "#@*"

type BoardDataDetail struct{
    Seq int
    Title string
    Author string
    DateInfo string
    Body string
    Name string
}

type BoardData struct{
    Seq int 
    Title string
    Author string
    DateInfo string
}

type Board struct {
    BoardTitle string
    BoardCnt int
    BoardPage int
    Keyword string
    Name string
    BoardDatas []BoardData
}

//session check
var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))
  
func getSession(request *http.Request) (name string) {
    if cookie, err := request.Cookie("myCookie"); err == nil {
        cookieValue := make(map[string]string)
        if err = cookieHandler.Decode("myCookie", cookie.Value, &cookieValue); err == nil {
            name = cookieValue["myCookie"]
        }
    }
    return name
}
  
func setSession(name string, response http.ResponseWriter) {
    value := map[string]string{
        "myCookie": name,
    }
    if encoded, err := cookieHandler.Encode("myCookie", value); err == nil {
        cookie := &http.Cookie{
            Name:  "myCookie",
            Value: encoded,
            Path:  "/",
            MaxAge: 3600,
        }
        http.SetCookie(response, cookie)
    }
}
  
func clearSession(response http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "myCookie",
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
        redirectTarget = "/"
    }
    http.Redirect(response, request, redirectTarget, 302)
}
  
func clearSessionHandler(response http.ResponseWriter, request *http.Request) {
    clearSession(response)
    http.Redirect(response, request, "/", 302)
}

func login_check(w http.ResponseWriter, r *http.Request) {
    //db connection
    db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
    db.Connect()

    //parameter setting
    email := r.FormValue("email")
    passwd := r.FormValue("passwd")

    //query
    rows, _ , _ := db.Query("SELECT name FROM login WHERE email='%s' AND passwd='%s'", email, passwd)

    //lgoin check
    if len(rows) == 1 {
        name := rows[0].Str(0)

        //set secure cookie
        setSession(name, w)

        fmt.Println("[" + name + ", " + email + "] Login sucess...")
        fmt.Fprintf(w, "0")
    } else {
        fmt.Println("[" + email + "] Login failed...")
        fmt.Fprintf(w, "Incorrect email or password.")
    }

}

func logout(w http.ResponseWriter, r *http.Request) {
    //get previous session
    name := getSession(r)

    clearSession(w)

    fmt.Println("[" + name + "] Logout!")
    fmt.Fprintf(w,"0")
}

func registration_check(w http.ResponseWriter, r *http.Request) {
    //db connection
    db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
    db.Connect()

    //parameter setting
    email := r.FormValue("email")
    passwd := r.FormValue("passwd")
    //passwd_repeat := r.FormValue("passwd_repeat") //password check work on client side... if you need to check in server then use it.
    name := r.FormValue("name")
    tel := r.FormValue("tel")

    //query
    row_email, _ , _ := db.Query("SELECT 1 FROM login WHERE email='%s'", email)
    row_name, _ , _ := db.Query("SELECT 1 FROM login WHERE name='%s'", name)

    //lgoin check
    if len(row_email) == 1 {
        fmt.Fprintf(w, "Email is duplicated.")
    } else if len(row_name) == 1 {
        fmt.Fprintf(w, "Name is duplicated.")
    } else {
        //regstrate user
        ins, _ := db.Prepare("INSERT INTO login(email, passwd, name, telephone) VALUES(?,?,?,?)")
        ins.Bind(email, passwd, name, tel)
        _, err := ins.Run()

        if err != nil {
            fmt.Println(err)
            fmt.Fprintf(w,"A error is occured when execute the insert statement.")
        } else {
            fmt.Fprintf(w, "0")
        }
    }

}

func call_main(w http.ResponseWriter, r *http.Request) {
    name := getSession(r)
    login_param := template.FuncMap{ "Name" : name, }

    //make a html file
    tmpl := template.Must(template.ParseFiles("WEB-INF/html/head.html"))
    tmpl.Execute(w, login_param)
    main_bytes, _ := ioutil.ReadFile("WEB-INF/html/main.html")
    fmt.Fprintf(w, string(main_bytes))
    tail_bytes, _ := ioutil.ReadFile("WEB-INF/html/tail.html")
    fmt.Fprintf(w, string(tail_bytes))

}

func call_board(w http.ResponseWriter, r *http.Request) {
    name := getSession(r)
    login_param := template.FuncMap{ "Name" : name, }

    //db connectio
    db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
    db.Connect()

    //parameter setting
    keyword := r.FormValue("keyword")
    page,_ := strconv.Atoi(r.FormValue("page"))
    if page == 0 {
      page = 1
    }

    var cnts []mysql.Row
    var rows []mysql.Row
    //query
    if keyword != ""  {
        cnts, _, _ = db.Query("SELECT CEIL(COUNT(1)/10) FROM board WHERE INSTR( CONCAT(Title,'"+sepa+"',Author,'"+sepa+"',Body), '%s') > 0", keyword)
        rows, _, _ = db.Query("SELECT A.* FROM (SELECT A.*, @rownum:=@rownum+1 rnum FROM (SELECT seq, Title, Author, concat(mod_dt, mod_tm) DateInfo FROM board WHERE INSTR( CONCAT(Title,'"+sepa+"',Author,'"+sepa+"',Body), '%s') > 0 ) A, (select @rownum:=0) R ORDER BY DateInfo DESC) A where rnum>%d and rnum<=%d", keyword, (page-1)*10, page*10)
    } else {
        cnts, _, _ = db.Query("SELECT CEIL(COUNT(1)/10) FROM board")
        rows, _, _ = db.Query("SELECT A.* FROM (SELECT seq, Title, Author, concat(mod_dt, mod_tm) DateInfo, @rownum:=@rownum+1 rnum FROM board, (select @rownum:=0) r ORDER BY seq DESC) A where rnum>%d and rnum<=%d",(page-1)*10, page*10)
    }

    cnt := cnts[0].Int(0)
    var item []BoardData
    var tmp BoardData
    for _, row := range rows {
        tmp.Seq = row.Int(0)
        tmp.Title = row.Str(1)
        tmp.Author = row.Str(2)
        tmp.DateInfo = row.Str(3)
        item=append(item,tmp)
    }

    items := Board{
        BoardTitle:"kyo board",
        BoardCnt:cnt,
        BoardPage:page,
        Keyword:keyword,
        Name:name,
        BoardDatas:item,
    }

    //make a html file
    tmpl := template.Must(template.ParseFiles("WEB-INF/html/head.html"))
    tmpl.Execute(w, login_param)
    tmpl = template.Must(template.ParseFiles("WEB-INF/html/board.html"))
    tmpl.Execute(w, items)
    tail_bytes, _ := ioutil.ReadFile("WEB-INF/html/tail.html")
    fmt.Fprintf(w, string(tail_bytes))
}

func call_board_detail(w http.ResponseWriter, r *http.Request) {
    name := getSession(r)
    login_param := template.FuncMap{ "Name" : name, }

    //db connection
    db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
    db.Connect()

    //parameter setting
    seq,_ := strconv.Atoi(r.FormValue("seq"))

    //query
    rows, _, _ := db.Query("SELECT Seq, Title, Author, concat(mod_dt, mod_tm) DateInfo, body FROM board where seq=%d", seq)

    if len(rows) == 1 {
        row := rows[0]

        item := BoardDataDetail{
            Seq:row.Int(0),
            Title:row.Str(1),
            Author:row.Str(2),
            DateInfo:row.Str(3),
            Body:row.Str(4),
            Name:name,
        }

        //make a html file
        tmpl := template.Must(template.ParseFiles("WEB-INF/html/head.html"))
        tmpl.Execute(w, login_param)
        tmpl = template.Must(template.ParseFiles("WEB-INF/html/board_detail.html"))
        tmpl.Execute(w, item)
        tail_bytes, _ := ioutil.ReadFile("WEB-INF/html/tail.html")
        fmt.Fprintf(w, string(tail_bytes))
    } else {
        fmt.Println("[" + name + "] Abnormal approach on board_detail.do.")
        fmt.Fprintf(w, "Abnormal approach."); 
        //http.Redirect(w, r,"/board.do", 302)
    }

}

func call_board_add(w http.ResponseWriter, r *http.Request) {
    name := getSession(r)
    login_param := template.FuncMap{ "Name" : name, }

    //make a html file
    tmpl := template.Must(template.ParseFiles("WEB-INF/html/head.html"))
    tmpl.Execute(w, login_param)
    tmpl = template.Must(template.ParseFiles("WEB-INF/html/board_add.html"))
    tmpl.Execute(w, login_param)
    tail_bytes, _ := ioutil.ReadFile("WEB-INF/html/tail.html")
    fmt.Fprintf(w, string(tail_bytes))

}

func board_add_check(w http.ResponseWriter, r *http.Request) {
    name := getSession(r)
    
    title := r.FormValue("title")
    author := name //get value from session value.
    body := r.FormValue("body")

    //connect db
    db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
    db.Connect()

    stmt, _ := db.Prepare("INSERT INTO board(title, author, body, mod_dt, mod_tm, reg_dt, reg_tm) VALUES(?, ?, ?, DATE_FORMAT(SYSDATE(0),'%Y%m%d'), DATE_FORMAT(SYSDATE(0),'%H%i%s'), DATE_FORMAT(SYSDATE(0),'%Y%m%d'), DATE_FORMAT(SYSDATE(0),'%H%i%s'))")
    _, err := stmt.Run(title, author, body)

    if err == nil {
        fmt.Fprintf(w, "0")
    } else {
        fmt.Println(err)
        fmt.Fprintf(w, "there are some error when execute insert statement");
    }

}

func board_modify(w http.ResponseWriter, r *http.Request) {
    name := getSession(r)
    author := r.FormValue("author")

    if( name != "" &&  name == author ){
        login_param := template.FuncMap{ "Name" : name, }
    
        //db connection
        db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
        db.Connect()

        //parameter setting
        seq,_ := strconv.Atoi(r.FormValue("seq"))

        //query
        rows, _, _ := db.Query("SELECT Seq, Title, Author, concat(mod_dt, mod_tm) DateInfo, body FROM board where seq=%d", seq)

        if len(rows) == 1 {
            row := rows[0]

            item := BoardDataDetail{
                Seq:row.Int(0),
                Title:row.Str(1),
                Author:row.Str(2),
                DateInfo:row.Str(3),
                Body:row.Str(4),
                Name:name,
            }

            //make a html file
            tmpl := template.Must(template.ParseFiles("WEB-INF/html/head.html"))
            tmpl.Execute(w, login_param)
            tmpl = template.Must(template.ParseFiles("WEB-INF/html/board_modify.html"))
            tmpl.Execute(w, item)
            tail_bytes, _ := ioutil.ReadFile("WEB-INF/html/tail.html")
            fmt.Fprintf(w, string(tail_bytes))
        } else {
            fmt.Println("[" + name + "] Abnormal approach(no data) on board_modify.do.")
            fmt.Fprintf(w, "Nodata error.");
        }
    } else {
        fmt.Println("[" + name + "] Abnormal approach(permission) on board_modify.do.")
        fmt.Fprintf(w, "Permission error.");
    }
}

func board_modify_check(w http.ResponseWriter, r *http.Request) {
    name := getSession(r)
    
    seq,_ := strconv.Atoi(r.FormValue("seq"))
    title := r.FormValue("title")
    author := r.FormValue("author")
    body := r.FormValue("body")

    if( name != "" &&  name == author ){
        //connect db
        db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
        db.Connect()

        stmt, _ := db.Prepare("UPDATE board SET title=?, body=?, mod_dt=DATE_FORMAT(SYSDATE(0),'%Y%m%d'), mod_tm=DATE_FORMAT(SYSDATE(0),'%H%i%s') where seq=?")
        _, err := stmt.Run(title, body, seq)

        if err == nil {
            fmt.Fprintf(w, "0")
        } else {
            fmt.Println(err)
            fmt.Fprintf(w, "there are some error when execute modify statement");
        }
    } else {
        fmt.Println("[" + name + "] Abnormal approach on board_modify_check.do.")
        fmt.Fprintf(w, "Permission error.");
    }
}

func board_delete(w http.ResponseWriter, r *http.Request) {
    name := getSession(r)
    
    seq := r.FormValue("seq")
    author := r.FormValue("author")

    if( name!="" &&  name == author ){
        //connect db
        db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
        db.Connect()

        stmt, _ := db.Prepare("DELETE FROM board where seq=?")
        _, err := stmt.Run(seq)

        if err == nil {
            fmt.Fprintf(w, "0")
        } else {
            fmt.Println(err)
            fmt.Fprintf(w, "there are some error when execute delete statement");
        }
    } else {
        fmt.Println("[" + name + "] Abnormal approach on board_delete.do.")
        fmt.Fprintf(w, "Permission error.");
    }
}

func main() {
    http.HandleFunc("/", call_main)
    http.HandleFunc("/board.do", call_board)
    http.HandleFunc("/board_detail.do", call_board_detail)
    http.HandleFunc("/board_add.do", call_board_add)
    http.HandleFunc("/board_add_check.do", board_add_check)
    http.HandleFunc("/board_modify.do", board_modify)
    http.HandleFunc("/board_modify_check.do", board_modify_check)
    http.HandleFunc("/board_delete.do", board_delete)
    http.HandleFunc("/login_check.do", login_check)
    http.HandleFunc("/logout.do", logout)
    http.HandleFunc("/registration_check.do", registration_check)
    http.Handle("/WEB-INF/", http.StripPrefix("/WEB-INF/", http.FileServer(http.Dir("WEB-INF"))))
 
    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
