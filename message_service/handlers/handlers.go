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
	"fmt"

	"urulink.go/message_service/db"
	"urulink.go/message_service/env"
	"urulink.go/message_service/rabbitmq"
	"urulink.go/message_service/redis"
)

// Handler struct contains references to various services and clients needed by the application
type Handler struct {
	Database       *db.Database              // Database connection instance
	EnvManger      *env.EnvManger            // Environment manager for config values
	RedisClient    *redis.RedisManager       // Redis client manager instance
	Ctx            context.Context           // Context for managing request lifetimes
	RabbitMQClient *rabbitmq.RabbitMQManager // RabbitMQ client manager instance
	MaxWorkers     int                       // Maximum number of worker goroutines
}

// Init initializes the Handler with necessary service connections and configurations
func Init() Handler {
	var err error
	var handlers_data Handler

	// Load environment configurations
	env := env.NewEnv()

	// Construct the DSN (Data Source Name) for MySQL using environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		env.DBUser,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	// Initialize the database connection using the constructed DSN
	handlers_data.Database, err = db.UruLinkInit(dsn)
	if err != nil {
		// Panic if there is an error connecting to the database
		panic("failed to connect to the database!")
	}

	// Assign the loaded environment configuration to the Handler struct
	handlers_data.EnvManger = env

	// Initialize the Redis client with environment configurations
	handlers_data.RedisClient = redis.UruLinkInit(env)

	// Initialize the RabbitMQ client with environment configurations
	handlers_data.RabbitMQClient, err = rabbitmq.UruLinkInit(env)
	if err != nil {
		// Panic if there is an error connecting to RabbitMQ
		panic("failed to connect to rabbitmq server!")
	}

	// Set the maximum number of workers for concurrent processing
	handlers_data.MaxWorkers = 10

	// Create a background context for handling request lifetimes
	handlers_data.Ctx = context.Background()

	// Return the populated Handler instance
	return handlers_data
}
