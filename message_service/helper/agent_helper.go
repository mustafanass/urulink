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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
)

// AgentService makes an HTTP request based on the provided method type (get, post, put).
// It includes setting an authorization header and supports JSON payloads.
// Returns the status code, response body, and any error encountered.
func AgentService[t any](url string, c *fiber.Ctx, bodyInput t, httpMethods string) (int, []byte, error) {
	switch httpMethods {
	case "get":
		// Create a GET request to the specified URL
		agent := fiber.Get(url)
		// Set authorization header
		agent.Set("Authorization", c.Get("Authorization"))
		// Send request and check for errors
		statusCode, body, errs := agent.Bytes()
		if len(errs) > 0 {
			return 0, nil, errors.New("problem with errs in agent Byte (Get http methods)")
		}
		return statusCode, body, nil
	case "post":
		// Create a POST request with JSON body
		agent := fiber.Post(url)
		agent.Set("Authorization", c.Get("Authorization"))
		agent.JSON(bodyInput)
		statusCode, body, errs := agent.Bytes()
		if len(errs) > 0 {
			return 0, nil, errors.New("problem with errs in agent Byte (Post http methods)")
		}
		return statusCode, body, nil
	case "put":
		// Create a PUT request with JSON body
		agent := fiber.Put(url)
		agent.Set("Authorization", c.Get("Authorization"))
		agent.JSON(bodyInput)
		statusCode, body, errs := agent.Bytes()
		if len(errs) > 0 {
			return 0, nil, errors.New("problem with errs in agent Byte (Put http methods)")
		}
		return statusCode, body, nil
	}
	return 0, []byte{}, errors.New("http methods not allowed")
}

// AgentSendFile uploads a file to a specified URL with a POST request, including
// an authorization header and setting appropriate content type.
// Returns the status code, response body, and any error encountered.
func AgentSendFile(c *fiber.Ctx, url, filePath string) (int, []byte, error) {
	// Open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Create a form file for file upload
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content to form part
	if _, err := io.Copy(part, file); err != nil {
		return 0, nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the multipart writer to finalize the request body
	if err := writer.Close(); err != nil {
		return 0, nil, fmt.Errorf("failed to close writer: %w", err)
	}

	bodyBytes := body.Bytes()

	// Create a POST request with the constructed form data
	agent := fiber.Post(url)
	agent.Set("Authorization", c.Get("Authorization"))
	agent.Set("Content-Type", writer.FormDataContentType())
	agent.Body(bodyBytes)

	// Send request and handle errors
	statusCode, respBody, errs := agent.Bytes()
	if len(errs) > 0 {
		return 0, nil, fmt.Errorf("failed to send request: %v", errs)
	}

	return statusCode, respBody, nil
}

// AgentResponse unmarshals a JSON response body into the specified model type `t`.
// Returns the unmarshaled model and any error encountered during unmarshaling.
func AgentResponse[t any](c *fiber.Ctx, body []byte) (t, error) {
	var model t
	// Unmarshal JSON response body into the model
	err := json.Unmarshal(body, &model)
	if err != nil {
		return model, fmt.Errorf("failed to Unmarshal request: %v", err)
	}
	return model, nil
}
