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

package env

import (
	"log"
	"os"
)

// EnvManger is a struct that stores environment configurations
// for connecting to MinIO storage service.
type EnvManger struct {
	MinioHost   string // MinIO server host address
	MinioKey    string // MinIO access key
	MinioSecret string // MinIO secret key
	MinioBucket string // MinIO bucket name
}

// NewEnv initializes a new EnvManger instance and loads environment variables.
// If any variable is missing, it logs a fatal error.
func NewEnv() *EnvManger {
	env := &EnvManger{} // Initialize EnvManger struct

	// Helper function to load environment variables into struct fields
	loadEnv := func(envVar string, target *string) {
		// Lookup environment variable
		value, ok := os.LookupEnv(envVar)
		if !ok {
			log.Fatalf("Environment variable %s not found", envVar) // Log fatal error if variable is missing
		}
		*target = value // Set the value to the target field
	}

	// Load MinIO-specific environment variables into struct fields
	loadEnv("MINIO_HOST", &env.MinioHost)
	loadEnv("MINIO_KEY", &env.MinioKey)
	loadEnv("MINIO_SECRET", &env.MinioSecret)
	loadEnv("MINIO_BUCKET", &env.MinioBucket)

	return env // Return populated EnvManger instance
}
