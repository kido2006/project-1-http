package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Đăng ký
func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl.ExecuteTemplate(w, "register.html", nil)

	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			//http.Error(w, "Username and password cannot be blank!", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Username and password cannot be blank!"))
			return
		}

		var exists string
		err := db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if exists != "" {
			http.Error(w, "Username already exists, choose another one!", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, string(hashedPassword))
		if err != nil {
			http.Error(w, "Error when adding user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Đăng nhập
func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl.ExecuteTemplate(w, "login.html", nil)

	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
		if err != nil {
			http.Error(w, "Wrong account or password!", http.StatusUnauthorized)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) != nil {
			http.Error(w, "Wrong account or password!", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "Login successful! Hello %s", username)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Đổi mật khẩu
func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl.ExecuteTemplate(w, "change_password.html", nil)

	case http.MethodPost:
		username := r.FormValue("username")
		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
		if err != nil {
			http.Error(w, "Tài khoản không tồn tại!", http.StatusUnauthorized)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(oldPassword)) != nil {
			http.Error(w, "Mật khẩu cũ không đúng!", http.StatusUnauthorized)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error ", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", string(hashedPassword), username)
		if err != nil {
			http.Error(w, "Lỗi khi đổi mật khẩu: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Đổi mật khẩu thành công!")

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Xóa user
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl.ExecuteTemplate(w, "delete.html", nil)

	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
		if err != nil {
			http.Error(w, "Tài khoản không tồn tại!", http.StatusUnauthorized)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) != nil {
			http.Error(w, "Sai mật khẩu!", http.StatusUnauthorized)
			return
		}

		_, err = db.Exec("DELETE FROM users WHERE username = ?", username)
		if err != nil {
			http.Error(w, "Lỗi khi xóa tài khoản: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Xóa tài khoản thành công!")

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
