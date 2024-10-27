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

// EnvManager struct holds configuration values loaded from environment variables
type EnvManger struct {
	FilesServiceUrl      string
	RabbitMQHost         string
	RabbitMQUser         string
	RabbitMQPassword     string
	RabbitMQPort         string
	RabbitMQExchangeName string
	RabbitMQQueueName    string
	DBHost               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBPort               string
	RedisHost            string
	RedisPort            string
	RedisPassword        string
}

// NewEnv initializes a new EnvManager instance, loading environment variables
// into their corresponding struct fields
func NewEnv() *EnvManger {
	env := &EnvManger{}

	// Helper function to load a single environment variable by name
	// Takes the variable name and a pointer to the target field in EnvManager
	loadEnv := func(envVar string, target *string) {
		// Lookup the environment variable value
		value, ok := os.LookupEnv(envVar)
		// Log a fatal error if the environment variable is not found
		if !ok {
			log.Fatalf("Environment variable %s not found", envVar)
		}
		// Set the target field in EnvManager to the variable's value
		*target = value
	}

	// Load RabbitMQ configuration values
	loadEnv("RABBITMQ_HOST", &env.RabbitMQHost)
	loadEnv("RABBITMQ_USER", &env.RabbitMQUser)
	loadEnv("RABBITMQ_PASSWORD", &env.RabbitMQPassword)
	loadEnv("RABBITMQ_PORT", &env.RabbitMQPort)
	loadEnv("RABBITMQ_EXCHANGE_NAME", &env.RabbitMQExchangeName)
	loadEnv("RABBITMQ_QUEUE_NAME", &env.RabbitMQQueueName)

	// Load Files Service URL
	loadEnv("URULINK_FILES_SERVICE", &env.DBHost)

	// Load Database configuration values
	loadEnv("DB_HOST", &env.DBHost)
	loadEnv("DB_USER", &env.DBUser)
	loadEnv("DB_PASSWORD", &env.DBPassword)
	loadEnv("DB_NAME", &env.DBName)
	loadEnv("DB_PORT", &env.DBPort)

	// Load Redis configuration values
	loadEnv("REDIS_HOST", &env.RedisHost)
	loadEnv("REDIS_PASSWORD", &env.RedisPassword)
	loadEnv("REDIS_PORT", &env.RedisPort)

	// Return the populated EnvManager instance
	return env
}
