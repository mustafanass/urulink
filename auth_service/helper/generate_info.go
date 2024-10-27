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
	"math/rand"
	"time"
)

// Define the set of characters to be used for generating numeric strings
const (
	number = "0123456789"
)

// generateRandomString creates a random string of specified length from the provided character set
func generateRandomString(chars string, length int) string {
	// Seed the random generator with current time to ensure randomness
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	// Initialize the result string
	var result string
	for i := 0; i < length; i++ {
		// Select a random index from the character set and append it to the result
		index := randomGenerator.Intn(len(chars))
		result += string(chars[index])
	}
	return result
}

// GenerateUid generates a unique identifier of fixed length using numeric characters
func GenerateUid() string {
	const uidLength = 10 // Define the length of the UID
	return generateRandomString(number, uidLength)
}
