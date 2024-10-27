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

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"urulink.com/routes"
)

// main function is the entry point of the application.
func main() {
	// Create a new Fiber application instance.
	app := fiber.New()

	// Set up the application routes by initializing them with the app instance.
	routes.SetRoutes(app)

	// Start the Fiber application and listen on port 8081, logging any fatal errors.
	log.Fatal(app.Listen(":8081"))
}
