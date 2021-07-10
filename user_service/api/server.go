package api

import (
	"encoding/json"
	"log"
	"net/http"

	service "github.com/Mahamed-Belkheir/scalechat-backend/user_service"
	app "github.com/Mahamed-Belkheir/scalechat-backend/user_service/app"
)

type result map[string]interface{}

func (r result) serialize() ([]byte, error) {
	return json.Marshal(r)
}

func StartWebServer(config service.Config, auth app.Authentication) {
	jwt := NewJWT(config)

	http.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			if username == "" {
				http.Error(rw, "username is required", http.StatusBadRequest)
				return
			}
			password := r.FormValue("password")
			if password == "" {
				http.Error(rw, "password is required", http.StatusBadRequest)
				return
			}
			user, err := auth.Login(username, password)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}
			if user == nil {
				http.Error(rw, "username not found", http.StatusNotFound)
				return
			}
			token, err := jwt.sign(user)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			message, err := result{
				"status": "success",
				"token":  token,
			}.serialize()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			rw.Header().Add("Content-Type", "application/json; charset=utf-8")
			rw.Write(message)
		} else {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/register", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			if username == "" {
				http.Error(rw, "username is required", http.StatusBadRequest)
				return
			}
			password := r.FormValue("password")
			if password == "" {
				http.Error(rw, "password is required", http.StatusBadRequest)
				return
			}
			user, err := auth.Register(username, password)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}
			token, err := jwt.sign(user)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			message, err := result{
				"status": "success",
				"token":  token,
			}.serialize()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			rw.Header().Add("Content-Type", "application/json; charset=utf-8")
			rw.Write(message)

		} else {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Printf("info: server listening at %v", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
