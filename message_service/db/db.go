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

package db

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"urulink.go/message_service/models"
)

// Database struct holds the GORM DB connection instance
type Database struct {
	Db *gorm.DB
}

// UruLinkInit initializes a database connection with specified settings
// Takes a DSN (Data Source Name) for MySQL, returns a Database instance or an error
func UruLinkInit(dsn string) (*Database, error) {
	// Open a new MySQL connection using the provided DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Retrieve the generic database object to configure connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set the maximum number of idle connections in the connection pool
	sqlDB.SetMaxIdleConns(15)
	// Set the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(30)
	// Set the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Return the Database struct with the configured GORM DB connection
	return &Database{
		Db: db,
	}, nil
}

// GetMessageByReceiverId retrieves all messages between the given sender and receiver IDs
// Orders messages by creation time in ascending order, returns a list of DirectMessage models
func (data Database) GetMessageByReceiverId(senderId, receiverId string) ([]models.DirectMessage, error) {
	var messages []models.DirectMessage
	// Query the database for messages between sender and receiver
	results := data.Db.Table("direct_message").
		Raw("SELECT * FROM direct_message WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) ORDER BY created_at ASC",
			senderId, receiverId, receiverId, senderId).Scan(&messages)
	return messages, results.Error
}

// CreateNewMsg creates a new message entry in the database within a transaction
// Rolls back the transaction if an error occurs or if a panic is recovered
func (data Database) CreateNewMsg(msg models.DirectMessage) error {
	// Begin a new transaction
	tx := data.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			// Log the panic and rollback the transaction if panic occurs
			log.Println("Recovered in CreateNewUser:", r)
			tx.Rollback()
		}
	}()

	// Attempt to create a new record in the 'direct_message' table
	result := tx.Table("direct_message").Create(&msg)
	if result.Error != nil {
		// Log and rollback the transaction if an error occurs
		log.Println("Error creating user:", result.Error)
		tx.Rollback()
		return result.Error
	}

	// Commit the transaction if no errors occur during message creation
	if err := tx.Commit().Error; err != nil {
		// Log and rollback the transaction if commit fails
		log.Println("Error committing transaction:", err)
		tx.Rollback()
		return err
	}

	// Return nil if the message was successfully created and transaction committed
	return nil
}
