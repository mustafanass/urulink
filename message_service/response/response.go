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

package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func HandleWebSocketError(c *websocket.Conn, errMsg string) error {
	return c.WriteJSON(fiber.Map{
		"error": errMsg,
	})
}

func HandleError(c *fiber.Ctx, statusCode int, err string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"error": err,
	})
}

func HandleInformation(c *fiber.Ctx, statusCode int, data interface{}) error {
	c.Status(statusCode)
	return c.JSON(data)
}
