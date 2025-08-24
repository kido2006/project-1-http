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

	//func (db *DB) Exec(query string, args ...any) (Result, error)

	_, err := db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, password)
	if err != nil {
		http.Error(w, "Lỗi DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)

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

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="vi">
		<head>
			<meta charset="UTF-8">
			<title>Trang chính</title>
		</head>
		<body>
			<h2>Đăng nhập thành công! Xin chào %s</h2>
			<p>Bạn muốn làm gì tiếp theo?</p>
			<ul>
				<li><a href="/change-password">Đổi mật khẩu</a></li>
				<li><a href="/delete-user" >Xóa tài khoản</a></li>
			</ul>
		</body>
		</html>
	`, username)

	fmt.Fprint(w, html)

}


// Hiển thị form đổi mật khẩu + xử lý đổi mật khẩu
func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "change_password.html", nil)
		return
	}

	username := r.FormValue("username")
	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Tài khoản không tồn tại!", http.StatusUnauthorized)
		return
	}

	if storedPassword != oldPassword {
		http.Error(w, "Mật khẩu cũ không đúng!", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", newPassword, username)
	if err != nil {
		http.Error(w, "Lỗi khi đổi mật khẩu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Đổi mật khẩu thành công!")
}

// Hiển thị form xóa user + xử lý xóa user
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "delete.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Kiểm tra user trước khi xóa
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Tài khoản không tồn tại!", http.StatusUnauthorized)
		return
	}

	if storedPassword != password {
		http.Error(w, "Sai mật khẩu!", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		http.Error(w, "Lỗi khi xóa tài khoản: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Xóa tài khoản thành công!")
}