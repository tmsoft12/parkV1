package operator

import (
	"fmt"
	modelscar "park/models/modelsCar"

	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan modelscar.Car_Model)
var Refresh = make(chan struct{})

func Ws(c *websocket.Conn) {
	defer func() {
		delete(clients, c)
		c.Close()
	}()

	clients[c] = true

	for {
		var msg interface{}
		fmt.Println(msg)
		if err := c.ReadJSON(&msg); err != nil {
			break
		}

		if _, ok := msg.(string); ok && msg == "refresh" {
			Refresh <- struct{}{}
		} else if car, ok := msg.(modelscar.Car_Model); ok {
			Broadcast <- car
		}
	}
}

func HandleMessages() {
	for {
		select {
		case car := <-Broadcast:
			for client := range clients {
				if err := client.WriteJSON(car); err != nil {
					client.Close()
					delete(clients, client)
				}
			}

		case <-Refresh:
			for client := range clients {
				if err := client.WriteJSON("refresh"); err != nil {
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
