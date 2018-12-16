# boardUsingGo
it is server program using GO language.

this directory is the test page for db connection and refer to "https://github.com/ziutek/mymysql".

before executing this directory, you should install the package using following command.

	$ go get github.com/ziutek/mymysql/thrsafe
	$ go get github.com/ziutek/mymysql/autorc
	$ go get github.com/ziutek/mymysql/godrv

	$ go get -v github.com/ziutek/mymysql/...

before executing this directory, you should create tags table which follow description.

create table tags(ID varchar(30), Name varchar(30));
insert into tags values('test1','test1');
insert into tags values('test2','test2');

and run it using command `go run main.go`.

