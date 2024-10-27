/*
 * Copyright (C) 2024 Mustafa Naseer
 *
 * This file is part of urulink chat application.
 *
 * urulink is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation version 3 of the License .
 *
 *
 * urulink is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with urulink. If not, see <http://www.gnu.org/licenses/>.
 */

package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"urulink.go/message_service/helper"
	"urulink.go/message_service/models"
	"urulink.go/message_service/response"
)

// WebSocketHandler handles incoming WebSocket connections and manages real-time messaging
func (h Handler) WebSocketHandler(c *websocket.Conn) {
	defer c.Close()

	// Extract JWT information for the user and retrieve user ID
	userJwtInfo := c.Locals("userJwtInfo").(models.ClientsLoginResponse)
	userId := userJwtInfo.Uid
	receiverId := c.Query("receiver_id")
	if receiverId == "" {
		// Return an error if receiver ID is not provided
		response.HandleWebSocketError(c, "Receiver ID is required")
		return
	}

	// Create a cancellable context for managing connection lifetime
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		h.RedisClient.RemoveClient(ctx, userId)
		h.RabbitMQClient.CancelConsume(userId)
	}()

	// If client is already connected, cancel previous RabbitMQ consumption
	if h.RedisClient.IsClientConnected(ctx, userId) {
		h.RabbitMQClient.CancelConsume(userId)
	}

	// Add client to Redis, storing their connection information
	if err := h.RedisClient.AddClient(ctx, userId, c); err != nil {
		helper.LogError(c, "Failed to add client", err)
		response.HandleWebSocketError(c, "Failed to add client")
		return
	}

	// Retrieve message history for the user and send it over WebSocket
	messages, err := h.Database.GetMessageByReceiverId(userId, receiverId)
	if err != nil {
		helper.LogError(c, "Failed to retrieve message history", err)
		response.HandleWebSocketError(c, "Failed to retrieve message history")
		return
	}

	for _, msg := range messages {
		messageData, _ := json.Marshal(msg)
		if err := c.WriteMessage(websocket.TextMessage, messageData); err != nil {
			helper.LogError(c, "Failed to send historical messages", err)
			response.HandleWebSocketError(c, "Failed to send historical messages")
			return
		}
	}

	// Channel to queue incoming messages for processing by workers
	jobs := make(chan models.DirectMessageInput, 100)

	// Start worker goroutines to process messages
	for w := 0; w < h.MaxWorkers; w++ {
		go h.worker(ctx, jobs, userId, receiverId)
	}

	// Goroutine to listen for RabbitMQ messages and forward them to WebSocket
	go func() {
		h.RabbitMQClient.ListenForMessages(ctx, userId, func(message []byte) error {
			var msg models.DirectMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				helper.LogError(c, "Failed to unmarshal incoming message", err)
				return err
			}
			// Forward message to WebSocket if it is intended for the current user
			if msg.ReceiverID == userId {
				if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
					helper.LogError(c, "Failed to send message", err)
					return err
				}
			}
			return nil
		}, h.EnvManger)
	}()

	// Main loop to receive messages from the WebSocket client
	for {
		var msgInput models.DirectMessageInput
		if err := c.ReadJSON(&msgInput); err != nil {
			helper.LogError(c, "Failed to read WebSocket message", err)
			return
		}
		// Send received message to jobs channel for further processing
		jobs <- msgInput
	}
}

// worker function processes messages from the jobs channel
func (h Handler) worker(ctx context.Context, jobs <-chan models.DirectMessageInput, senderId, receiverId string) {
	for msgInput := range jobs {
		h.processMessage(ctx, msgInput, senderId, receiverId)
	}
}

// processMessage processes a single message by creating, saving, and optionally sending it
func (h Handler) processMessage(ctx context.Context, msgInput models.DirectMessageInput, senderId, receiverId string) {
	// Check if the receiver is online using Redis
	isOnline := h.RedisClient.IsClientConnected(ctx, receiverId)
	msgStatus := 2 // Offline by default
	if isOnline {
		msgStatus = 1 // Set status to online if receiver is connected
	}
	helper.LogInfo("Receiver online status", map[string]interface{}{
		"receiverId": receiverId,
		"isOnline":   isOnline,
	})

	// Create and save message in database
	msg, err := h.createMessage(msgInput, msgStatus, senderId, receiverId)
	if err != nil {
		helper.LogError(nil, "Failed to create and store message", err)
		return
	}

	// Marshal the message into JSON format for transmission
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		helper.LogError(nil, "Failed to marshal message", err)
		return
	}

	// Send message via RabbitMQ to the intended receiver
	if err := h.RabbitMQClient.PublishMessage(ctx, receiverId, msgBytes, h.EnvManger); err != nil {
		helper.LogError(nil, "Failed to send message via RabbitMQ", err)
	}
}

// createMessage constructs a message object and saves it to the database
func (h Handler) createMessage(msgInput models.DirectMessageInput, msgStatus int, userId, receiverId string) (models.DirectMessage, error) {
	var c *fiber.Ctx
	var conn *websocket.Conn

	helper.LogInfo("Creating message", map[string]interface{}{
		"senderId":   userId,
		"receiverId": receiverId,
		"status":     msgStatus,
	})

	var msg models.DirectMessage
	// If the message includes files, upload them and store the file path
	if msgInput.ContentType == "files" {
		statusCode, body, err := helper.AgentSendFile(c, h.EnvManger.FilesServiceUrl+"/upload", msgInput.Content)
		if err != nil {
			helper.LogError(conn, "Failed to upload files", err)
			return models.DirectMessage{}, errors.New("failed to upload files")
		}
		if statusCode != 200 {
			helper.LogError(conn, "File upload failed", errors.New("invalid status code"))
			return models.DirectMessage{}, errors.New("file upload failed")
		}
		// Decode and save the uploaded file path
		filePath, err := helper.AgentResponse[string](c, body)
		if err != nil {
			helper.LogError(conn, "Failed to decode file response", err)
			return models.DirectMessage{}, errors.New("failed to decode file response")
		}
		msg.FilePath = filePath
		helper.LogInfo("File upload processed successfully", nil)
	}

	// Populate message fields
	msg = models.DirectMessage{
		SenderID:    userId,
		ReceiverID:  receiverId,
		Content:     msgInput.Content,
		ContentType: msgInput.ContentType,
		Status:      msgStatus,
		CreatedAt:   time.Now().Unix(),
	}

	// Save message to database
	if err := h.Database.CreateNewMsg(msg); err != nil {
		helper.LogError(nil, "Failed to save message in database", err)
		return models.DirectMessage{}, errors.New("failed to save message in database")
	}
	helper.LogInfo("Message saved to database successfully", nil)
	return msg, nil
}
