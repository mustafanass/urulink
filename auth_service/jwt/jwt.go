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

package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"urulink.com/env"
)

// JWTManager is responsible for managing JWT creation and verification.
// It relies on environment variables loaded into EnvManger for the secret key.
type JWTManager struct {
	EnvMange *env.EnvManger // Holds environment manager to access JWT secret key
}

// CheckAccessToken verifies and parses the access token to extract the UID and username.
// It returns the UID, username, and an error if parsing or validation fails.
func (jm *JWTManager) CheckAccessToken(accessToken string) (string, string, error) {
	clientsSecretKey := jm.EnvMange.AccessTokenKey // Retrieve the secret key from environment manager

	// Parse the JWT token with claims
	token, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(clientsSecretKey), nil // Provide the secret key for signature verification
	})
	if err != nil {
		return "", "", errors.New("failed to parse access token") // Return error if token parsing fails
	}

	// Type assertion to access claims within the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("failed to assert claims") // Return error if claims cannot be asserted
	}

	// Extract UID from claims, ensuring it is a string
	uid, ok := claims["uid"].(string)
	if !ok {
		return "", "", errors.New("invalid UID in claims") // Return error if UID is missing or invalid
	}

	// Extract username from claims, ensuring it is a string
	username, ok := claims["username"].(string)
	if !ok {
		return "", "", errors.New("invalid Username in claims") // Return error if username is missing or invalid
	}

	return uid, username, nil // Return UID and username if parsing succeeds
}

// CreateUsersJwt generates a JWT token with the user's UID and username as claims.
// The token includes an expiration time set to 72 hours from creation.
func (jm *JWTManager) CreateUsersJwt(uid, username string) (string, error) {
	// Set up claims with UID, username, and an expiration time of 72 hours
	claims := jwt.MapClaims{
		"uid":      uid,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create a new JWT with HS256 signing method and the defined claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and return the signed token string
	tokenString, err := jwtToken.SignedString([]byte(jm.EnvMange.AccessTokenKey))
	return tokenString, err
}
