package api

import (
	"encoding/json"
	"log"
	"net/http"

	app "github.com/Mahamed-Belkheir/scalechat-backend/chat_service/app"
)

type ChatHandler struct {
	jwt            JWT
	chatController app.ChatControl
}

type h map[string]interface{}

func (c ChatHandler) json(obj interface{}, w http.ResponseWriter) {
	str, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(str)
}

func (c ChatHandler) getChat(w http.ResponseWriter, r *http.Request) {
	chats, err := c.chatController.GetChats()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	c.json(chats, w)
}
func (c ChatHandler) postChat(userId string, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := r.Form.Get("name")
	if name == "" {
		http.Error(w, "name is a required property", http.StatusBadRequest)
		return
	}
	err = c.chatController.AddChat(name, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.json(h{
		"status": "success",
	}, w)
}
func (c ChatHandler) delChat(userId string, w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is a required property", http.StatusBadRequest)
		return
	}
	err := c.chatController.DeleteChat(userId, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.json(h{
		"status": "success",
	}, w)
}

func (c ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := c.jwt.verify(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case "GET":
		c.getChat(w, r)
		break
	case "POST":
		c.postChat(userId, w, r)
		break
	case "DELETE":
		c.delChat(userId, w, r)
		break
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		break
	}
}

var _ http.Handler = ChatHandler{}
