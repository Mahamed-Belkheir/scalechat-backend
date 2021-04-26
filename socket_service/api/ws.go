package api

import (
	"net/http"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/broker"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/pool"
	"nhooyr.io/websocket"
)

type wsHandler struct {
	br     *broker.MessageBroker
	pool   *pool.Pool
	config service.Config
	jwt    JWT
}

func newWSHandler(br *broker.MessageBroker, pool *pool.Pool, config service.Config, jwt JWT) *wsHandler {
	return &wsHandler{
		br,
		pool,
		config,
		jwt,
	}
}

func (h *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("name")
	if roomName == "" {
		http.Error(w, "name is a required query parameter", http.StatusBadRequest)
		return
	}
	userId, err := h.jwt.verify(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.pool.IsFull() {
		http.Error(w, "service unavilable", http.StatusServiceUnavailable)
		return
	}
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	work := newConnection(userId, roomName, conn, h.br)
	h.pool.AddJob(&work)
}
