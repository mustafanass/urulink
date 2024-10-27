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

const (
	// numericChars defines the set of characters to use for generating random strings.
	numericChars = "0123456789abcdefghijklmnopqrstuvwxyz"
)

// generateRandomString creates a random string of specified length using the provided characters.
// chars: string of characters to choose from.
// length: desired length of the generated string.
func generateRandomString(chars string, length int) string {
	// Initialize a new random source based on the current time in nanoseconds for unique randomness.
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	// Build the result string by randomly selecting characters from the input set.
	var result string
	for i := 0; i < length; i++ {
		index := randomGenerator.Intn(len(chars)) // Select a random index
		result += string(chars[index])            // Append character at the index to the result
	}
	return result
}

// GenerateConnId generates a unique connection ID of predefined length using numeric and alphabetic characters.
// Returns a random string ID of length 14.
func GenerateConnId() string {
	const uidLength = 14
	return generateRandomString(numericChars, uidLength)
}
