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
	"urulink.com/models"
)

// Database struct to hold the GORM DB instance for database operations.
type Database struct {
	Db *gorm.DB
}

// UruLinkInit initializes a new database connection using the provided DSN (Data Source Name).
// Sets connection pool parameters such as max idle connections, max open connections, and max lifetime.
func UruLinkInit(dsn string) (*Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure the database connection pool settings.
	sqlDB.SetMaxIdleConns(15)           // Maximum number of idle connections.
	sqlDB.SetMaxOpenConns(30)           // Maximum number of open connections.
	sqlDB.SetConnMaxLifetime(time.Hour) // Lifetime of each connection.

	return &Database{
		Db: db,
	}, nil
}

// CreateNewUser creates a new user in the database within a transaction to ensure data integrity.
// If any error occurs during the process, it rolls back the transaction to avoid partial data writes.
func (data Database) CreateNewUser(userInfo models.UsersInfo) error {
	// Begin a new transaction.
	tx := data.Db.Begin()
	defer func() {
		// Recover in case of panic and roll back the transaction.
		if r := recover(); r != nil {
			log.Println("Recovered in CreateNewUser:", r)
			tx.Rollback()
		}
	}()

	// Attempt to create the new user record in the users_info table.
	result := tx.Table("users_info").Create(&userInfo)
	if result.Error != nil {
		log.Println("Error creating user:", result.Error)
		tx.Rollback() // Rollback transaction in case of error.
		return result.Error
	}

	// Commit the transaction if no error occurs.
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		tx.Rollback() // Rollback in case of commit error.
		return err
	}

	return nil
}

// CheckUsername verifies if a username exists in the users_info table and returns the user's ID if found.
func (data Database) CheckUsername(username string) (int, error) {
	var userInfo int
	// Execute a raw SQL query to find a user ID based on the given username.
	result := data.Db.Table("users_info").Raw("SELECT id FROM users_info WHERE username = ?", username).Scan(&userInfo)
	return userInfo, result.Error
}

// CheckUserUid verifies if a user UID exists in the users_info table and returns the user's ID if found.
func (data Database) CheckUserUid(uid string) (int, error) {
	var userInfo int
	// Execute a raw SQL query to find a user ID based on the given UID.
	result := data.Db.Table("users_info").Raw("SELECT id FROM users_info WHERE uid = ?", uid).Scan(&userInfo)
	return userInfo, result.Error
}

// GetUserInfoByUsername retrieves full user information from the users_info table based on the provided username.
func (data Database) GetUserInfoByUsername(username string) (models.UsersInfo, error) {
	var userInfo models.UsersInfo
	// Execute a raw SQL query to fetch all user details for the specified username.
	result := data.Db.Table("users_info").Raw("SELECT * FROM users_info WHERE username = ?", username).Scan(&userInfo)
	return userInfo, result.Error
}
