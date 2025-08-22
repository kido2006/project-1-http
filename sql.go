package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
    // MySQL
    var err error
    db, err = sql.Open("mysql", "niga:123456789@tcp(localhost:3306)/myapp")
    if err != nil {
        log.Fatal("Lỗi kết nối MySQL:", err)
    }
    defer db.Close()
    if err = db.Ping(); err != nil {
        log.Fatal("Không thể kết nối MySQL:", err)
    }

    fmt.Println("Kết nối MySQL thành công!")

    //  HTTP
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Helo ửold")
    })

    fmt.Println("Server chạy trên LAN")
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
