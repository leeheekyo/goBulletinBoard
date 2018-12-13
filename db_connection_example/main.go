package main

import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

/*
 * Tag... - a very simple struct
 */
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
    fmt.Println("Go MySQL Tutorial")
    
    // Open up our database connection.
    // I've set up a database on my local machine using phpmyadmin.
    // The database is called testDb
    // userID:password@tcp(127.0.0.1:3306)/dbName
    db, err := sql.Open("mysql", "root:kyo@tcp(127.0.0.1:3306)/test")
    
    // if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
    }
    
    // defer the close till after the main function has finished
    // executing 
    defer db.Close()

    // Execute the query
    results, err := db.Query("SELECT id, name FROM tags")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    for results.Next() {
        var tag Tag
        // for each row, scan the result into our tag composite object
        err = results.Scan(&tag.ID, &tag.Name)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
        // and then print out the tag's Name attribute
        //log.Printf(tag.Name)
    	fmt.Println(tag.Name)
    }
}
