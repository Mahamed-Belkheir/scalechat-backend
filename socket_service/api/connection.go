package api

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/broker"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type connection struct {
	userId      string
	roomName    string
	con         *websocket.Conn
	br          *broker.MessageBroker
	sendQueue   chan service.Message
	closeSignal chan bool
	wg          *sync.WaitGroup
	once        *sync.Once
}

func newConnection(userId, roomName string, con *websocket.Conn, br *broker.MessageBroker) connection {
	return connection{
		userId,
		roomName,
		con,
		br,
		make(chan service.Message),
		make(chan bool),
		&sync.WaitGroup{},
		&sync.Once{},
	}
}

func (c *connection) exit() {
	c.br.Unregister(c.userId, c.roomName)
	close(c.closeSignal)
	close(c.sendQueue)
}

func (c *connection) Run() {
	c.br.Register(c.userId, c.roomName, c.sendQueue)
	c.wg.Add(2)
	go c.recieve()
	go c.transmit()
	c.wg.Wait()
}

func (c *connection) recieve() {
	defer c.wg.Done()
	for {
		select {
		case <-c.closeSignal:
			return
		default:
			var msg service.Message
			err := wsjson.Read(context.Background(), c.con, &msg)
			if err != nil {
				log.Printf("error: client %v ran into error %v", c.userId, err)
				c.once.Do(c.exit)
				return
			}
			msg.UserID = c.userId
			msg.Room = c.roomName
			msg.CreatedAt = time.Now().Unix()
			c.br.SendMessage(msg)
		}
	}
}

func (c *connection) transmit() {
	defer c.wg.Done()
	for {
		select {
		case <-c.closeSignal:
			return
		case msg := <-c.sendQueue:
			res, err := json.Marshal(msg)
			if err != nil {
				log.Printf("error: parsing json %v", err)
				c.once.Do(c.exit)
				return
			}
			c.con.Write(context.Background(), websocket.MessageText, res)
		}
	}
}
