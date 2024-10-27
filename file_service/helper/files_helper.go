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
	"errors"
	"mime/multipart"
	"path/filepath"

	"urulink.com/file_service/models" // Importing the models package for FileTypeConfig
)

// Define constants for maximum file sizes
const (
	ImageLimit = 20 * 1024 * 1024  // Maximum limit for image files (20 MB)
	VideoLimit = 100 * 1024 * 1024 // Maximum limit for video files (100 MB)
	FileLimit  = 50 * 1024 * 1024  // Maximum limit for general files (50 MB)
)

// Mapping of file extensions to their configurations, including maximum file sizes
var fileTypeMap = map[string]models.FileTypeConfig{
	".jpeg": {MaxSize: ImageLimit}, // JPEG image configuration
	".jpg":  {MaxSize: ImageLimit}, // JPG image configuration
	".png":  {MaxSize: ImageLimit}, // PNG image configuration
	".mp4":  {MaxSize: VideoLimit}, // MP4 video configuration
	".pdf":  {MaxSize: FileLimit},  // PDF file configuration
	".txt":  {MaxSize: FileLimit},  // TXT file configuration
	".zip":  {MaxSize: FileLimit},  // ZIP file configuration
}

// ValidateFile checks the file type and size for a given multipart file header
func ValidateFile(fileHeader *multipart.FileHeader) error {
	fileExt := filepath.Ext(fileHeader.Filename)  // Get the file extension
	fileSize := fileHeader.Size                   // Get the file size
	fileConfig, supported := fileTypeMap[fileExt] // Check if the file extension is supported

	// Return an error if the file type is unsupported
	if !supported {
		return errors.New("file type is unsupported")
	}

	// Return an error if the file size exceeds the allowed limit
	if fileSize > fileConfig.MaxSize {
		return errors.New("file size exceeds the allowed limit")
	}
	return nil // Return nil if the file is valid
}
