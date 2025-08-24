package main

import (
	"fmt"
	"net/http"
)

// Hiển thị form đăng ký + xử lý đăng ký
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "register.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	_, err := db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, password)
	if err != nil {
		http.Error(w, "Lỗi DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)

	fmt.Fprintln(w, "Đăng ký thành công! <a href='/login'>Đăng nhập</a>")

}

// Hiển thị form đăng nhập + xử lý đăng nhập
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Sai tài khoản hoặc mật khẩu!", http.StatusUnauthorized)
		return
	}
	if storedPassword != password {
		http.Error(w, "Sai tài khoản hoặc mật khẩu!", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Đăng nhập thành công! Xin chào %s", username)
}
