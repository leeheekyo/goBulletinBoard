# boardUsingGo
It is server program using GO language.

This directory is the test page for db connection and refer to "https://github.com/ziutek/mymysql".

Before executing this directory, you should install the package using following command.

	$ go get github.com/ziutek/mymysql/thrsafe
	$ go get github.com/ziutek/mymysql/autorc
	$ go get github.com/ziutek/mymysql/godrv

	$ go get -v github.com/ziutek/mymysql/...

Before executing this directory, you should create tags table which follow description.

    mysql> create table tags(ID varchar(30), Name varchar(30));
    mysql> insert into tags values('test1','test1');
    mysql> insert into tags values('test2','test2');

And run it using command `go run main.go`.

REF) There are two example to execute sql statement.
    //case 1
    db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
    db.Connect()
    rows, _, _ := db.Query("SELECT col1, col2 FROM your_table")
    col1 := rows[0].Str(0)
    col2 := rows[0].Str(2)

    //case 2
    db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
    db.Connect()
    stmt, _ := db.Prepare("INSERT INTO your_table VALUES (?, ?)")
    _, err := stmt.Run(title, author, body)
    if( err != nil ){
        fmt.Println("Success")
    }

