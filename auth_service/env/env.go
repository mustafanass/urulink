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

// EnvManger holds the environment variables required for application configuration.
type EnvManger struct {
	AccessTokenKey string // Key used for access token encryption/decryption.
	DBHost         string // Database host address.
	DBUser         string // Database username.
	DBPassword     string // Database password.
	DBName         string // Name of the database.
	DBPort         string // Database port number.
}

// NewEnv initializes a new EnvManger instance and loads environment variables into it.
// Exits the program with an error if any required variable is not set.
func NewEnv() *EnvManger {
	env := &EnvManger{}

	// Helper function to load a single environment variable.
	// If the variable is not found, logs a fatal error and terminates the program.
	loadEnv := func(envVar string, target *string) {
		value, ok := os.LookupEnv(envVar)
		if !ok {
			log.Fatalf("Environment variable %s not found", envVar) // Log an error if envVar is missing.
		}
		*target = value // Set the value of the environment variable to the target field.
	}

	// Load each required environment variable into the EnvManger fields.
	loadEnv("ACCESS_TOKEN_KEY", &env.AccessTokenKey)
	loadEnv("DB_HOST", &env.DBHost)
	loadEnv("DB_USER", &env.DBUser)
	loadEnv("DB_PASSWORD", &env.DBPassword)
	loadEnv("DB_NAME", &env.DBName)
	loadEnv("DB_PORT", &env.DBPort)

	return env // Return the populated EnvManger instance.
}
