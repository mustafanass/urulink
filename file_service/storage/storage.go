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

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage defines a struct for MinIO storage, holding the client and bucket name.
type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

// InitMinio initializes a MinIO client and creates a bucket if it doesn't exist.
func InitMinio(url, accessKey, secretKey, bucketName string, ctx context.Context) (*MinioStorage, error) {
	// Create a new MinIO client with the provided URL and credentials
	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""), // Use static credentials
		Secure: false,                                             // Set to true for HTTPS
	})
	if err != nil {
		return nil, err
	}

	// Check if the specified bucket exists
	existsk, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	// If the bucket does not exist, create it
	if !existsk {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"}) // Default region
		if err != nil {
			return nil, err
		}
	}

	// Return a pointer to MinioStorage with the initialized client and bucket name
	return &MinioStorage{
		Client:     minioClient,
		BucketName: bucketName,
	}, nil
}
