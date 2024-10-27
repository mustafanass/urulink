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

// Define a constant string containing allowed characters for random string generation
const (
	number = "0123456789abcdefghijklmnobqsrwxyz" // Allowed characters (numbers and lowercase letters)
)

// generateRandomString generates a random string of specified length from the provided characters
func generateRandomString(chars string, length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	var result string // Initialize an empty string to hold the result
	for i := 0; i < length; i++ {
		index := randomGenerator.Intn(len(chars))
		result += string(chars[index])
	}
	return result // Return the generated random string
}

// GenerateFilesName generates a random file name of fixed length (20 characters)
func GenerateFilesName() string {
	const uidLength = 20                           // Define the length of the generated file name
	return generateRandomString(number, uidLength) // Call the helper function to generate the file name
}
