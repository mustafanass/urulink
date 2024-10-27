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
	"log"

	"github.com/redis/go-redis/v9"
	"urulink.go/message_service/env"
)

// RedisManager struct to manage Redis client and message channel
type RedisManager struct {
	Client    *redis.Client // Redis client for database operations
	messageCh chan []byte   // Channel for handling messages
}

// UruLinkInit initializes the Redis client and checks the connection
func UruLinkInit(env *env.EnvManger) *RedisManager {
	// Create a new Redis client with provided configuration from environment manager
	client := redis.NewClient(&redis.Options{
		Addr:     env.RedisHost + ":" + env.RedisPort, // Redis server address
		Password: env.RedisPassword,                   // Password for Redis server
		DB:       0,                                   // Default database to use
	})

	ctx := context.Background()         // Create a background context for Redis operations
	_, err := client.Ping(ctx).Result() // Ping the Redis server to test connection
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err) // Log fatal error if connection fails
	}

	messageCh := make(chan []byte)                             // Create a channel for message handling
	return &RedisManager{Client: client, messageCh: messageCh} // Return RedisManager instance
}
