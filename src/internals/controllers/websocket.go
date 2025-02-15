package controllers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Message represents different types of WebSocket messages
type Message struct {
	Type      string      `json:"type"`
	UserId    string      `json:"userId"`
	RoomId    string      `json:"roomId"`
	Content   string      `json:"content,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Offer     interface{} `json:"offer,omitempty"`
	Answer    interface{} `json:"answer,omitempty"`
	Candidate interface{} `json:"candidate,omitempty"`
}

var rooms = make(map[string]map[string]*websocket.Conn)

// HandleWebSocket handles WebSocket connections for specific room and user
func HandleWebSocket(c *fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		roomID := conn.Params("room_id")
		userID := conn.Params("user_id")

		if roomID == "" || userID == "" {
			log.Println("Invalid room or user parameters")
			return
		}

		// Initialize room if it doesn't exist
		if rooms[roomID] == nil {
			rooms[roomID] = make(map[string]*websocket.Conn)
		}

		rooms[roomID][userID] = conn
		log.Printf("User %s joined room %s", userID, roomID)

		broadcastMessage(Message{
			Type:      "user-joined",
			UserId:    userID,
			RoomId:    roomID,
			Timestamp: time.Now(),
		}, roomID, userID)

		defer func() {
			delete(rooms[roomID], userID)
			if len(rooms[roomID]) == 0 {
				delete(rooms, roomID)
			}
			log.Printf("User %s left room %s", userID, roomID)

			broadcastMessage(Message{
				Type:      "user-left",
				UserId:    userID,
				RoomId:    roomID,
				Timestamp: time.Now(),
			}, roomID, userID)
		}()

		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				break
			}

			var message Message
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("parse error: %v", err)
				continue
			}

			message.Timestamp = time.Now()
			message.UserId = userID
			message.RoomId = roomID

			switch message.Type {
			case "chat":
				broadcastMessage(message, roomID, "")
			case "offer", "answer", "ice-candidate":
				if recipientConn, exists := rooms[roomID][message.UserId]; exists {
					data, _ := json.Marshal(message)
					if err := recipientConn.WriteMessage(messageType, data); err != nil {
						log.Printf("write error: %v", err)
					}
				}
			default:
				log.Printf("Unknown message type: %v", message.Type)
			}
		}
	})(c)
}

// broadcastMessage sends a message to all users in a room except the sender
func broadcastMessage(message Message, roomID string, excludeUserID string) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	for userID, conn := range rooms[roomID] {
		if userID != excludeUserID {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("write error for user %s: %v", userID, err)
			}
		}
	}
}
