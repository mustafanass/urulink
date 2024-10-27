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

package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"urulink.go/message_service/handlers"
	"urulink.go/message_service/helper"
	"urulink.go/message_service/models"
)

// WebSocketConnection validates a WebSocket connection request with authorization check before allowing the upgrade.
// It retrieves the user information and stores it in the context for downstream access.
func WebSocketConnection(h *handlers.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the request is a WebSocket upgrade request
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		// Retrieve the authorization service URL from environment variables
		authUrl, ok := os.LookupEnv("URULINK_AUTH_SERVICE")
		if !ok {
			// Panic if the authorization service environment variable is not set
			panic("URULINK_AUTH_SERVICE Env not found !")
		}

		// Get the access token from the request headers
		accessToken := c.Get("Authorization")

		// Set up a POST request to the auth service to verify the user's login status
		agent := fiber.Post(authUrl + "/check-login")
		agent.Set("Authorization", accessToken) // Add Authorization header with the access token

		// Execute the request and capture the response
		statusCode, body, errs := agent.Bytes()
		// If there are any errors or the status code is not 200, respond with "Unauthorized"
		if len(errs) > 0 || statusCode != 200 {
			fmt.Println(errs, statusCode)
			return c.Status(401).SendString("Unauthorized requests 1")
		}

		// Parse the response to extract the user's JWT information
		userJwtInfo, err := helper.AgentResponse[models.ClientsLoginResponse](c, body)
		if err != nil {
			// If parsing fails, respond with "Unauthorized"
			return c.Status(401).SendString("Unauthorized requests 2")
		}

		// Store the user information in the context for access in WebSocket handlers
		c.Locals("userJwtInfo", userJwtInfo)
		c.Locals("allowed", true) // Mark the connection as authorized

		// Proceed with the next middleware or route handler
		return c.Next()
	}
}
