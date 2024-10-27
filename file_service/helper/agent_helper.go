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

package helper

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// AgentService is a generic function that sends HTTP requests (GET, POST, PUT) to a specified URL
func AgentService[t any](url string, c *fiber.Ctx, bodyInput t, httpMethods string) (int, []byte, error) {
	switch httpMethods {
	case "get":
		accessToken := c.Get("Authorization")
		agent := fiber.Get(url)
		agent.Set("Authorization", accessToken)
		statusCode, body, errs := agent.Bytes()
		if len(errs) > 0 {
			return 0, nil, errors.New("problem with errs in agent Byte (Get http methods)") // Handle any errors during the GET request
		}
		return statusCode, body, nil

	case "post":
		accessToken := c.Get("Authorization")
		agent := fiber.Post(url)
		agent.Set("Authorization", accessToken)
		agent.JSON(bodyInput)
		statusCode, body, errs := agent.Bytes()
		if len(errs) > 0 {
			return 0, nil, errors.New("problem with errs in agent Byte (Post http methods)") // Handle any errors during the POST request
		}
		return statusCode, body, nil

	case "put":
		accessToken := c.Get("Authorization")
		agent := fiber.Put(url)
		agent.Set("Authorization", accessToken)
		agent.JSON(bodyInput)
		statusCode, body, errs := agent.Bytes()
		if len(errs) > 0 {
			return 0, nil, errors.New("problem with errs in agent Byte (Put http methods)") // Handle any errors during the PUT request
		}
		return statusCode, body, nil
	}

	// Return an error if the provided HTTP method is not allowed
	return 0, []byte{}, errors.New("http methods not allowed")
}

// AgentResponse handles the response from an HTTP request, unmarshalling the JSON body into a specified model type
func AgentResponse[t any](c *fiber.Ctx, body []byte) (t, error) {
	var model t                         // Declare a variable of the generic type
	err := json.Unmarshal(body, &model) // Unmarshal the JSON body into the model variable
	if err != nil {
		return model, errors.New("error with AgentResponse in Unmarshal functions") // Handle any errors during unmarshalling
	}
	return model, nil // Return the unmarshalled model
}
