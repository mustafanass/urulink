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
	"fmt"

	"urulink.com/db"
	"urulink.com/env"
	"urulink.com/jwt"
)

// Handler struct aggregates dependencies including Database, JWT manager, and environment variables manager
type Handler struct {
	Database  *db.Database
	JWT       *jwt.JWTManager
	EnvManger *env.EnvManger
}

// Init initializes the Handler with all necessary dependencies, including database connection and JWT manager.
func Init() Handler {
	var err error
	var handlers_data Handler

	// Initialize environment variables manager
	env := env.NewEnv()

	// Construct DSN (Data Source Name) for MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		env.DBUser,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	// Initialize and connect to the database with the constructed DSN
	handlers_data.Database, err = db.UruLinkInit(dsn)
	if err != nil {
		// Panic if database connection fails
		panic("failed to connect to the database:")
	}

	// Assign the environment manager to the Handler
	handlers_data.EnvManger = env

	// Initialize JWT manager with the environment manager
	handlers_data.JWT = &jwt.JWTManager{EnvMange: handlers_data.EnvManger}

	// Return the fully initialized Handler
	return handlers_data
}
