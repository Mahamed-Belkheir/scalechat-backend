package api

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/app"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type connection struct {
	userId      string
	roomName    string
	con         *websocket.Conn
	sockApp     app.SocketApplication
	sendQueue   chan service.Message
	closeSignal chan bool
	wg          *sync.WaitGroup
	once        *sync.Once
}

func newConnection(userId, roomName string, con *websocket.Conn, sockApp app.SocketApplication) connection {
	return connection{
		userId,
		roomName,
		con,
		sockApp,
		make(chan service.Message),
		make(chan bool),
		&sync.WaitGroup{},
		&sync.Once{},
	}
}

func (c *connection) exit() {
	log.Printf("debug: user %s disconnected from room %s", c.userId, c.roomName)
	c.sockApp.Unregister(c.userId, c.roomName, c.sendQueue)
	close(c.closeSignal)
	close(c.sendQueue)
}

func (c *connection) Run() {
	log.Printf("debug: user %s connected to room %s", c.userId, c.roomName)
	c.sockApp.Register(c.userId, c.roomName, c.sendQueue)
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
			log.Printf("debug: user %s sent a message to room %s \n %s", c.userId, c.roomName, msg.Body)
			msg.UserID = c.userId
			msg.Room = c.roomName
			msg.CreatedAt = time.Now().Unix()
			c.sockApp.SendMessage(msg)
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
			log.Printf("debug: sending message to user %s on room %s \n %s", c.userId, c.roomName, msg.Body)
			c.con.Write(context.Background(), websocket.MessageText, res)
		}
	}
}
