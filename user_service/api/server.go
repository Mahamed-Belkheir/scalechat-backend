package api

import (
	"log"
	"net/http"

	service "github.com/Mahamed-Belkheir/scalechat-backend/user_service"
	app "github.com/Mahamed-Belkheir/scalechat-backend/user_service/app"
)

func StartWebServer(config service.Config, auth app.Authentication) {
	jwt := NewJWT(config)

	http.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			if username == "" {
				http.Error(rw, "username is required", http.StatusBadRequest)
			}
			password := r.FormValue("password")
			if password == "" {
				http.Error(rw, "password is required", http.StatusBadRequest)
			}
			user, err := auth.Login(username, password)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
			}
			token, err := jwt.sign(user)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			rw.Write([]byte(token))
		} else {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/register", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			if username == "" {
				http.Error(rw, "username is required", http.StatusBadRequest)
			}
			password := r.FormValue("password")
			if password == "" {
				http.Error(rw, "password is required", http.StatusBadRequest)
			}
			user, err := auth.Register(username, password)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
			}
			token, err := jwt.sign(user)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			rw.Write([]byte(token))

		} else {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
