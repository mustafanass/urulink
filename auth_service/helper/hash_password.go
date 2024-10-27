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

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the provided plaintext password using bcrypt and returns the hashed string.
// It takes in a string password and returns the hashed password or an error if hashing fails.
func HashPassword(password string) (string, error) {
	// Generate hash from password with bcrypt's default cost level
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err // Return empty string and error if hashing fails
	}
	return string(hash), nil // Return the hashed password as a string
}

// CheckPassword compares a plaintext password with a hashed password and returns true if they match.
// It takes a plaintext password and a hashed password, and uses bcrypt's comparison function.
func CheckPassword(password, hash string) bool {
	// Compare the plaintext password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Return true if there's no error, meaning passwords match
}
