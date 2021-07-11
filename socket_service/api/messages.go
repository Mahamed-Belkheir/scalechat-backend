package api

import (
	"encoding/json"
	"net/http"

	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/app"
)

type messagesHandler struct {
	messagesApp app.MessagesApplication
	jwt         JWT
}

func newMessagesHandler(messagesApp app.MessagesApplication, jwt JWT) messagesHandler {
	return messagesHandler{messagesApp, jwt}
}

type h map[string]interface{}

func (c messagesHandler) json(obj interface{}, w http.ResponseWriter) {
	str, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(str)
}

func (m messagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := m.jwt.verify(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case "":
	case "GET":
		m.getHandle(w, r)
		return
	case "DELETE":
		m.deleteHandle(w, r, userId)
		return
	default:
		http.Error(w, "method not supported", http.StatusMethodNotAllowed)
	}
}

func (m messagesHandler) getHandle(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("name")
	if room == "" {
		http.Error(w, "name is a required query parameter", http.StatusBadRequest)
		return
	}
	messages, err := m.messagesApp.GetRoomMessages(room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m.json(messages, w)
}

func (m messagesHandler) deleteHandle(w http.ResponseWriter, r *http.Request, userId string) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is a required query parameter", http.StatusBadRequest)
		return
	}
	err := m.messagesApp.DeleteUserMessage(id, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.json(h{
		"status": "success",
	}, w)
}
