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
	"os" // Package for accessing environment variables

	"github.com/gofiber/fiber/v2"
	"urulink.com/file_service/handlers"
	"urulink.com/file_service/helper"
	"urulink.com/file_service/models"
)

// HttpAuth is a middleware function that checks the authorization of incoming requests
func HttpAuth(h *handlers.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Look up the authentication URL from environment variables
		authUrl, ok := os.LookupEnv("URUFI_AUTH_URL")
		if !ok {
			panic("URUFI_AUTH_URL Env not found !")
		}
		accessToken := c.Get("Authorization") // Retrieve the Authorization header from the request

		// Prepare a request to the authentication service to check login status
		agent := fiber.Post(authUrl + "/check-login")
		agent.Set("Authorization", accessToken) // Set the Authorization header for the request

		// Send the request and capture the status code, response body, and errors
		statusCode, body, errs := agent.Bytes()
		if len(errs) > 0 || statusCode != 200 {
			return c.Status(401).SendString("Unauthorized requests")
		}

		// Parse the response body into the ClientsLoginResponse model
		userJwtInfo, err := helper.AgentResponse[models.ClientsLoginResponse](c, body)
		if err != nil {
			return c.Status(401).SendString("Unauthorized requests")
		}

		// Store the parsed user JWT info in the context for later use
		c.Locals("userJwtInfo", userJwtInfo)
		return c.Next() // Proceed to the next middleware or handler
	}
}
