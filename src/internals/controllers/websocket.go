package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var rooms = make(map[string]map[string]*websocket.Conn)

func HandleWebSocket(c *fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		roomID := c.Params("room_id")
		userID := c.Params("user_id")

		// Initialize the room if it doesn't exist
		if rooms[roomID] == nil {
			rooms[roomID] = make(map[string]*websocket.Conn)
		}

		// Store the connection
		rooms[roomID][userID] = c
		defer delete(rooms[roomID], userID)

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			// Broadcast the message to all other users in the same room
			for uid, conn := range rooms[roomID] {
				if uid != userID {
					if err := conn.WriteMessage(mt, msg); err != nil {
						log.Println("write:", err)
						break
					}
				}
			}
		}
	})(c)
}
