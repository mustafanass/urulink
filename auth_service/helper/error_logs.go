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
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

var logWriter io.Writer

func init() {
	f, err := os.OpenFile("service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logWriter = io.MultiWriter(f, os.Stdout) // Log to both file and stdout
	log.SetOutput(logWriter)
}

// LogError logs an error with structured information for Fiber context
func LogError(c *fiber.Ctx, message string, err error) {
	logData := fmt.Sprintf("path: %s, method: %s, clientIp: %s, query: %s", c.Path(), c.Method(), c.IP(), c.OriginalURL())
	if err != nil {
		logData += fmt.Sprintf(", error: %v", err)
	} else {
		logData += ", error: No error information"
	}
	log.Printf("[ERROR] %s: %s", message, logData)
}

// LogInfo logs an informational message with structured data for Fiber context
func LogInfo(c *fiber.Ctx, message string, data map[string]interface{}) {
	logData := fmt.Sprintf("path: %s, method: %s, clientIp: %s, query: %s", c.Path(), c.Method(), c.IP(), c.OriginalURL())
	for k, v := range data {
		logData += fmt.Sprintf(", %s: %v", k, v)
	}
	log.Printf("[INFO] %s: %s", message, logData)
}
