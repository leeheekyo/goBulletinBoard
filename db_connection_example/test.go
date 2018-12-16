package main

import (
	"os"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func main() {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "kyo", "test")

	err := db.Connect()
	if err != nil {
		panic(err)
	}

	rows, res, err := db.Query("select * from tags where id = '%s'","test1")
	if err != nil {
		panic(err)
	}
	if res == nil {
		panic(res)
	}

	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}
		// You can get specific value from a row
		val1 := row[1].([]byte)

		// You can use it directly if conversion isn't needed
		os.Stdout.Write(val1)
		fmt.Println()
	}
}
