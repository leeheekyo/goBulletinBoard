# boardUsingGo
It is server program using GO language.

This directory is the main proeject of this repository.

Before executing this directory, you should install the package using following command.

    $ go get github.com/ziutek/mymysql/thrsafe
    $ go get github.com/ziutek/mymysql/autorc
    $ go get github.com/ziutek/mymysql/godrv

    $ go get -v github.com/ziutek/mymysql/...
    $ go get github.com/gorilla/securecookie

And before executing this directory, you should create tags table which follow description.

    mysql> create table login(email varchar(50) PRIMARY KEY, passwd varchar(512), name varchar(30), telephone varchar(12));
    mysql> create table board( seq INT(10) AUTO_INCREMENT PRIMARY KEY, title VARCHAR(250), author VARCHAR(30), body LONGTEXT, mod_dt CHAR(8), mod_tm CHAR(6), reg_dt CHAR(8), reg_tm CHAR(6) );

And run it using command `go run go/main.go`.

### tutorial
 - start the go commant
![](images/start.png)

 - go to the webpage(127.0.0.1:8080)
![](images/start_page.png)

 - click the login button
![](images/_login_click_page.png)

 - click the "Sign up" (it located bottom of the login form) 
![](images/sign_up_click_page.png)

 - fill out the registration form and click the "Restrater" button
![](images/registration_page.png)

 - fill out the login form and click the "Login" button
![](images/login_fill_out_page.png)

 - click the "Board" tab.
![](images/board_page.png)

 - you can add the board and show the detail infomation of board and modify and delete it.
![](images/board_detail_page.png)

