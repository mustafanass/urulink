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
	"strings"

	"github.com/gofiber/fiber/v2"
	"urulink.com/helper"
	"urulink.com/models"
	"urulink.com/response"
)

// Register handles user registration by validating input, checking username availability,
// hashing the password, and creating a new user in the database.
func (h Handler) Register(c *fiber.Ctx) error {
	var userInfoInput models.UsersInfoInput
	// Parse JSON request body into userInfoInput struct
	if err := c.BodyParser(&userInfoInput); err != nil {
		helper.LogError(c, "Failed to parse request body in Register", err)
		return c.SendStatus(400) // Bad Request if parsing fails
	}

	// Check if the username already exists in the database
	id, err := h.Database.CheckUsername(userInfoInput.Username)
	if err != nil {
		helper.LogError(c, "Failed to check username in database", err)
		return c.SendStatus(500) // Internal Server Error if database check fails
	}

	if id != 0 {
		// Username is already taken, return a 403 Forbidden status
		helper.LogInfo(c, "Username already exists", map[string]interface{}{"username": userInfoInput.Username})
		return c.SendStatus(403)
	}

	// Hash the user's password before storing it
	password, err := helper.HashPassword(userInfoInput.Password)
	if err != nil {
		helper.LogError(c, "Failed to hash password", err)
		return c.SendStatus(500)
	}

	// Create a UsersInfo object with hashed password and generated UID
	userInfo := models.UsersInfo{
		Uid:      helper.GenerateUid(),
		Username: userInfoInput.Username,
		Password: password,
		Name:     userInfoInput.Name,
	}

	// Attempt to create the new user in the database
	if err = h.Database.CreateNewUser(userInfo); err != nil {
		helper.LogError(c, "Failed to create new user in database", err)
		return c.SendStatus(500)
	}

	// Log success and return 200 OK status
	helper.LogInfo(c, "User registered successfully", map[string]interface{}{"username": userInfo.Username})
	return c.SendStatus(200)
}

// Login handles user login by validating credentials, checking password, and generating JWT.
func (h Handler) Login(c *fiber.Ctx) error {
	var userLoginInfo models.UserLoginInfo
	// Parse JSON request body into userLoginInfo struct
	if err := c.BodyParser(&userLoginInfo); err != nil {
		helper.LogError(c, "Failed to parse request body in Login", err)
		return c.SendStatus(400) // Bad Request if parsing fails
	}

	// Retrieve user info from the database by username
	userInfo, err := h.Database.GetUserInfoByUsername(userLoginInfo.Username)
	if err != nil {
		helper.LogError(c, "Failed to get user info from database", err)
		return c.SendStatus(500) // Internal Server Error if retrieval fails
	}
	if userInfo.Id == 0 {
		// Username not found, return 403 Forbidden status
		helper.LogInfo(c, "User not found", map[string]interface{}{"username": userLoginInfo.Username})
		return c.SendStatus(403)
	}

	// Check if the provided password matches the stored hashed password
	hashPasswordCheck := helper.CheckPassword(userLoginInfo.Password, userInfo.Password)
	if !hashPasswordCheck {
		// Invalid password, return 403 Forbidden status
		helper.LogInfo(c, "Invalid password", map[string]interface{}{"username": userLoginInfo.Username})
		return c.SendStatus(403)
	}

	// Generate JWT for the user upon successful login
	urufiAccessToken, err := h.JWT.CreateUsersJwt(userInfo.Uid, userInfo.Username)
	if err != nil {
		helper.LogError(c, "Failed to create JWT", err)
		return c.SendStatus(403) // Return 403 if JWT generation fails
	}

	// Log success and return the generated JWT
	helper.LogInfo(c, "User logged in successfully", map[string]interface{}{"username": userInfo.Username})
	return response.HandleInformation(c, 200, urufiAccessToken)
}

// CheckLogin verifies the provided access token and returns user info if valid.
func (h Handler) CheckLogin(c *fiber.Ctx) error {
	// Extract the access token from the Authorization header
	authHeader := c.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate the access token and retrieve UID and username
	uid, username, err := h.JWT.CheckAccessToken(accessToken)
	if err != nil || uid == "" || username == "" {
		// Invalid or expired token, return 401 Unauthorized status
		helper.LogError(c, "Invalid access token", err)
		return c.SendStatus(401)
	}

	// Create response structure with UID and username
	jwtUserInfo := models.ClientsLoginResponse{
		Uid:      uid,
		Username: username,
	}

	// Log success and return the user info
	helper.LogInfo(c, "User access token validated", map[string]interface{}{"uid": uid})
	return response.HandleInformation(c, 200, jwtUserInfo)
}
