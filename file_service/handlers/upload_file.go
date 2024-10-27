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
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"urulink.com/file_service/helper"
	"urulink.com/file_service/models"
	"urulink.com/file_service/response"
)

// UploadFile handles file uploads through a multipart form
func (h *Handler) UploadFile(c *fiber.Ctx) error {
	// Parse the multipart form from the incoming request
	form, err := c.MultipartForm()
	if err != nil {
		helper.LogError(c, "Failed to parse multipart form", err)
		return c.Status(400).SendString(err.Error())
	}

	// Retrieve the list of uploaded files from the form
	files := form.File["files"]
	if len(files) == 0 {
		helper.LogError(c, "No files found in request", fmt.Errorf("files not found"))
		return c.Status(400).SendString("files not found")
	}

	var filesInfo []models.FileSender // Slice to hold information about uploaded files
	for _, file := range files {      // Iterate over each file in the uploaded files
		// Validate the file's type and size
		err := helper.ValidateFile(file)
		if err != nil {
			helper.LogError(c, "File validation failed", err)
			return c.Status(405).SendString(err.Error())
		}

		// Generate a random file name with the original file extension
		randomFileName := fmt.Sprintf("%s%s", helper.GenerateFilesName(), filepath.Ext(file.Filename))

		// Open the file for reading
		fileData, err := file.Open()
		if err != nil {
			helper.LogError(c, "Failed to open file", err) // Log error if file opening fails
			return c.Status(500).SendString(err.Error())   // Return 500 status with error message
		}
		defer fileData.Close() // Ensure the file is closed after processing

		// Upload the file to MinIO storage
		if err := h.Minio.UploadFile(h.Ctx, fileData, randomFileName, file.Size); err != nil {
			helper.LogError(c, "Failed to upload file to Minio", err) // Log error if upload fails
			return c.Status(500).SendString(err.Error())              // Return 500 status with error message
		}

		// Generate a presigned URL for the uploaded file, valid for 48 hours
		presignedURL, err := h.Minio.GeneratePresignedURL(h.Ctx, randomFileName, 48*time.Hour)
		if err != nil {
			helper.LogError(c, "Failed to generate presigned URL", err) // Log error if URL generation fails
			return c.Status(400).SendString(err.Error())                // Return 400 status with error message
		}

		// Append file information (name and URL) to the slice
		filesInfo = append(filesInfo, models.FileSender{
			FileName: randomFileName,
			FileUrl:  presignedURL,
		})
	}

	// Log successful upload with the count of files
	helper.LogInfo(c, "Files uploaded successfully", map[string]interface{}{
		"filesCount": len(filesInfo), // Include the count of uploaded files in the log
	})

	// Return a response with a 200 status and the information of uploaded files
	return response.HandleInformation(c, 200, filesInfo)
}
