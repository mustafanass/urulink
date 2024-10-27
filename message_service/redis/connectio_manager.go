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

package redis

import (
	"context"

	"github.com/gofiber/websocket/v2"
	"urulink.go/message_service/helper"
)

// AddClient adds a new WebSocket client connection to Redis with a unique connection ID
func (cm *RedisManager) AddClient(ctx context.Context, userId string, conn *websocket.Conn) error {
	// Generate a unique connection ID for the client
	connectionID := helper.GenerateConnId()
	// Define the Redis key for storing the client's connection information
	connKey := "user:" + userId

	// Store the connection ID in Redis under the user's key
	err := cm.Client.Set(ctx, connKey, connectionID, 0).Err()
	if err != nil {
		return err // return error if setting the value in Redis fails
	}

	// Log the addition of the client to Redis for tracking purposes
	helper.LogInfo("Client added to Redis", map[string]interface{}{
		"userId":       userId,       // ID of the user being added
		"connectionID": connectionID, // unique connection ID generated
	})
	return nil
}

// IsClientConnected checks if a client is currently connected by looking up their connection key in Redis
func (cm *RedisManager) IsClientConnected(ctx context.Context, userId string) bool {
	// Define the Redis key for the user's connection information
	connKey := "user:" + userId
	// Attempt to get the connection ID from Redis
	_, err := cm.Client.Get(ctx, connKey).Result()
	// Return true if no error occurs, indicating the client is connected
	return err == nil
}

// RemoveClient removes a client's connection information from Redis
func (cm *RedisManager) RemoveClient(ctx context.Context, userId string) error {
	// Define the Redis key for the user's connection information
	connKey := "user:" + userId

	// Delete the user's connection information from Redis
	err := cm.Client.Del(ctx, connKey).Err()
	if err != nil {
		return err // return error if deletion from Redis fails
	}

	// Log the removal of the client from Redis for tracking purposes
	helper.LogInfo("Client removed from Redis", map[string]interface{}{
		"userId": userId, // ID of the user being removed
	})
	return nil
}
