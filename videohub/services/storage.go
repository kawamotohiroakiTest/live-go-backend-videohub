package services

import (
	"fmt"
	"live/common"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/minio/minio-go/v7"
	minioCredentials "github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
	Client      *s3.S3
	MinioClient *minio.Client
	Bucket      string
}

// Initialize S3 client for production
func NewStorageService() (*StorageService, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("AWS_REGION is not set")
	}

	bucketName := os.Getenv("STORAGE_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("STORAGE_BUCKET is not set")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKey == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID is not set")
	}

	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("AWS_SECRET_ACCESS_KEY is not set")
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize AWS session: %w", err)
	}

	client := s3.New(sess)

	return &StorageService{
		Client: client,
		Bucket: bucketName,
	}, nil
}

// Initialize MinIO client for local development
func InitMinioService() (*StorageService, error) {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("STORAGE_ENDPOINT is not set")
	}

	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("STORAGE_ACCESS_KEY is not set")
	}

	secretKey := os.Getenv("STORAGE_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("STORAGE_SECRET_KEY is not set")
	}

	bucketName := os.Getenv("STORAGE_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("STORAGE_BUCKET is not set")
	}

	useSSL := os.Getenv("STORAGE_USE_SSL") == "true"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  minioCredentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize MinIO client: %w", err)
	}

	return &StorageService{
		MinioClient: minioClient,
		Bucket:      bucketName,
	}, nil
}

func (s *StorageService) GetVideoPresignedURL(videoPath string) (string, error) {

	if s.MinioClient != nil {
		minioEndpoint := os.Getenv("MINIO_ENDPOINT")
		if minioEndpoint == "" {
			return "", fmt.Errorf("MINIO_ENDPOINT is not set")
		}

		urlStr := fmt.Sprintf("http://%s/%s/%s", "localhost:9000", s.Bucket, videoPath)

		return urlStr, nil
	} else if s.Client != nil {
		req, _ := s.Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(videoPath),
		})

		presignedURL, err := req.Presign(24 * time.Hour)
		if err != nil {
			common.LogVideoHubInfo(fmt.Sprintf("presignedURL: %s", presignedURL))

			return "", fmt.Errorf("Failed to generate presigned URL for video: %w", err)
		}

		return presignedURL, nil
	} else {
		return "", fmt.Errorf("Storage service is not initialized")
	}
}
