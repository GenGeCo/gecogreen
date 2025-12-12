package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// R2Storage handles file uploads to Cloudflare R2
type R2Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

// NewR2Storage creates a new R2 storage client
func NewR2Storage(accountID, accessKeyID, secretKey, bucketName string) (*R2Storage, error) {
	// R2 endpoint
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	// Create custom resolver for R2
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// Public URL for accessing files
	publicURL := fmt.Sprintf("https://pub-%s.r2.dev", accountID)

	return &R2Storage{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicURL,
	}, nil
}

// Upload uploads a file to R2 and returns the public URL
func (s *R2Storage) Upload(ctx context.Context, file io.Reader, filename string, contentType string, folder string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload to R2
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(uniqueName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return public URL
	return fmt.Sprintf("%s/%s", s.publicURL, uniqueName), nil
}

// Delete deletes a file from R2
func (s *R2Storage) Delete(ctx context.Context, fileURL string) error {
	// Extract key from URL
	key := strings.TrimPrefix(fileURL, s.publicURL+"/")

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	return err
}

// GeneratePresignedURL generates a presigned URL for direct upload
func (s *R2Storage) GeneratePresignedURL(ctx context.Context, filename string, contentType string, folder string) (string, string, error) {
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(uniqueName),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s", s.publicURL, uniqueName)
	return request.URL, publicURL, nil
}

// IsValidImageType checks if the content type is a valid image
func IsValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}

// MaxFileSize is the maximum file size for uploads (5MB)
const MaxFileSize = 5 * 1024 * 1024
