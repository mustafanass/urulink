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

	"urulink.com/file_service/env"     // Package to manage environment variables
	"urulink.com/file_service/storage" // Package for MinIO storage operations
)

// Handler struct stores environment configuration, MinIO storage instance, and context
type Handler struct {
	EnvManger *env.EnvManger        // Environment manager for accessing MinIO configurations
	Minio     *storage.MinioStorage // Instance of MinIO storage to handle file operations
	Ctx       context.Context       // Context for handling request lifetimes
}

// Init initializes the Handler struct and connects to the MinIO server
func Init() Handler {
	var err error
	var handlers_data Handler
	env := env.NewEnv()

	handlers_data.Ctx = context.Background() // Initialize a background context for the handler

	// Initialize MinIO storage instance with environment variables and context
	handlers_data.Minio, err = storage.InitMinio(env.MinioHost, env.MinioKey, env.MinioSecret, env.MinioBucket, handlers_data.Ctx)
	if err != nil {
		fmt.Println(err)                                 // Print the error if initialization fails
		panic("failed to connect to the Minio server :") // Panic with an error message
	}
	handlers_data.EnvManger = env

	return handlers_data
}
