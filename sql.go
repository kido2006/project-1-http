package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var tpl *template.Template

func main() {
	var err error
	// username: niga , password: 123456789 , database: myapp
	db, err = sql.Open("mysql", "niga:123456789@tcp(localhost:3306)/myapp")
	if err != nil {
		log.Fatal("MySQL connection error:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Can not connect to MySQL:", err)
	} 
	
	fmt.Println("MySQL connection successful!")
	

	// load templates
	tpl = template.Must(template.ParseGlob("*.html"))

	// route
	//http.HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/change-password", changePasswordHandler)
	http.HandleFunc("/delete", deleteUserHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
