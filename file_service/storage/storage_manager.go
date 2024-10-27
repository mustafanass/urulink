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

package storage

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

// UploadFile uploads a file to the MinIO storage using the provided context, file data, object name, and file size.
func (ms MinioStorage) UploadFile(ctx context.Context, fileData io.Reader, objectName string, fileSize int64) error {
	// PutObject uploads the file to the specified bucket with the provided object name and file data.
	_, err := ms.Client.PutObject(ctx, ms.BucketName, objectName, fileData, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DownloadFile downloads a file from MinIO storage to the specified destination path.
func (ms *MinioStorage) DownloadFile(ctx context.Context, objectName string, destPath string) error {
	// FGetObject retrieves the object from the specified bucket and saves it to the local file system at the destination path.
	err := ms.Client.FGetObject(ctx, ms.BucketName, objectName, destPath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

// GeneratePresignedURL generates a presigned URL for accessing the specified object in MinIO for a limited time.
func (ms *MinioStorage) GeneratePresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := url.Values{} // Initialize URL parameters for the presigned URL
	// PresignedGetObject generates a presigned URL to access the object for the specified expiry duration.
	presignedURL, err := ms.Client.PresignedGetObject(ctx, ms.BucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
