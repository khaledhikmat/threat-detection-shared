package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/khaledhikmat/threat-detection-shared/equates"
)

const (
	region       = "us-east-2"
	bucketPrefix = "threat-detection"
)

var s3Client *s3.Client

func storeClipViaAWS(ctx context.Context, clip equates.RecordingClip) (string, error) {
	err := makeS3Client(ctx)
	if err != nil {
		return "", err
	}

	// Create a bucket for each camera if it doesn't exist
	bucketName := makeBucketName(clip.Camera)
	fmt.Printf("Bucket name: %s\n", bucketName)
	bucketExists, err := bucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}

	fmt.Printf("Bucket exists: %t\n", bucketExists)
	if !bucketExists {
		err = createBucket(ctx, bucketName, region)
		if err != nil {
			return "", err
		}
	}

	// Upoad the file to the bucket using a unique object key
	fmt.Printf("Uploading a file to: %s\n", bucketName)
	key := makeKeyName(clip.Camera, clip.ID)
	err = uploadFile(ctx, bucketName, key, clip.LocalReference)
	if err != nil {
		return "", err
	}

	// URL must be in the format: https://<bucket-name>.s3.<region>.amazonaws.com/<object-key>
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, key), nil
}

func retrieveClipFromAWS(ctx context.Context, clip equates.RecordingClip) ([]byte, error) {
	err := makeS3Client(ctx)
	if err != nil {
		return []byte{}, err
	}

	bucketName := makeBucketName(clip.Camera)
	key := makeKeyName(clip.Camera, clip.ID)

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, key, err)
		return []byte{}, err
	}
	defer result.Body.Close()

	// Read the object's body into the clip's body
	b, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Printf("Couldn't read object %v:%v. Here's why: %v\n", bucketName, key, err)
		return []byte{}, err
	}

	return b, nil
}

func downloadClipFromAWS(ctx context.Context, clip equates.RecordingClip) ([]byte, error) {
	err := makeS3Client(ctx)
	if err != nil {
		return []byte{}, err
	}

	bucketName := makeBucketName(clip.Camera)
	key := makeKeyName(clip.Camera, clip.ID)

	var partMiBs int64 = 10
	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
	})
	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("Couldn't download large object from %v:%v. Here's why: %v\n",
			bucketName, key, err)
	}
	return buffer.Bytes(), err
}

// BucketExists checks whether a bucket exists in the current account.
func bucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Printf("bucketExists - Error: %v\n", err)
		return false, nil // Assume the error is due to the bucket not existing
	}

	return true, err
}

// CreateBucket creates a bucket with the specified name in the specified Region.
func createBucket(ctx context.Context, name string, region string) error {
	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func uploadFile(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func makeS3Client(ctx context.Context) error {
	if s3Client == nil {
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
		)

		if err != nil {
			return err
		}

		s3Client = s3.NewFromConfig(cfg)
	}

	return nil
}

func makeBucketName(camera string) string {
	return fmt.Sprintf("%s-%s", bucketPrefix, strings.ToLower(camera))
}

func makeKeyName(camera, id string) string {
	return fmt.Sprintf("%s-%s", camera, id)
}
