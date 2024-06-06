package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/khaledhikmat/threat-detection-shared/models"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

const (
	region       = "us-east-2"
	bucketPrefix = "threat-detection"
)

func NewAwsStorage(cfgsvc config.IService) IService {
	s3, err := makeS3Client(context.Background())
	if err != nil {
		panic(err)
	}

	return &awsStorage{
		S3Client: s3,
		CfgSvc:   cfgsvc,
	}
}

type awsStorage struct {
	S3Client *s3.Client
	CfgSvc   config.IService
}

// Not supported in AWS yet
func (s *awsStorage) StoreKeyValue(ctx context.Context, store, key, value string) error {
	return nil
}

func (s *awsStorage) StoreRecordingClip(ctx context.Context, clip models.RecordingClip) (string, error) {
	// Create a bucket for each camera if it doesn't exist
	bucketName := makeBucketName(clip.Camera)
	fmt.Printf("Bucket name: %s\n", bucketName)
	bucketExists, err := s.bucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}

	fmt.Printf("Bucket exists: %t\n", bucketExists)
	if !bucketExists {
		err = s.createBucket(ctx, bucketName, region)
		if err != nil {
			return "", err
		}
	}

	// Upoad the file to the bucket using a unique object key
	fmt.Printf("Uploading a file to: %s\n", bucketName)
	key := makeKeyName(clip.Camera, clip.ID)
	err = s.uploadFile(ctx, bucketName, key, clip.LocalReference)
	if err != nil {
		return "", err
	}

	// URL must be in the format: https://<bucket-name>.s3.<region>.amazonaws.com/<object-key>
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, key), nil
}

func (s *awsStorage) RetrieveRecordingClip(ctx context.Context, clip models.RecordingClip) ([]byte, error) {
	bucketName := makeBucketName(clip.Camera)
	key := makeKeyName(clip.Camera, clip.ID)

	result, err := s.S3Client.GetObject(ctx, &s3.GetObjectInput{
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

func (s *awsStorage) DownloadRecordingClip(ctx context.Context, clip models.RecordingClip) ([]byte, error) {
	bucketName := makeBucketName(clip.Camera)
	key := makeKeyName(clip.Camera, clip.ID)

	var partMiBs int64 = 10
	downloader := manager.NewDownloader(s.S3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
	})
	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("Couldn't download large object from %v:%v. Here's why: %v\n",
			bucketName, key, err)
	}

	return buffer.Bytes(), err
}

func (s *awsStorage) Finalize() {
}

// BucketExists checks whether a bucket exists in the current account.
func (s *awsStorage) bucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := s.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Printf("bucketExists - Error: %v\n", err)
		return false, nil // Assume the error is due to the bucket not existing
	}

	return true, err
}

// CreateBucket creates a bucket with the specified name in the specified Region.
func (s *awsStorage) createBucket(ctx context.Context, name string, region string) error {
	_, err := s.S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
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
func (s *awsStorage) uploadFile(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func makeS3Client(ctx context.Context) (*s3.Client, error) {
	var client *s3.Client

	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
	)

	if err != nil {
		return nil, err
	}

	client = s3.NewFromConfig(cfg)

	return client, nil
}

func makeBucketName(camera string) string {
	prefix := bucketPrefix
	if os.Getenv("AWS_BUCKET_PREFIX") != "" {
		prefix = os.Getenv("AWS_BUCKET_PREFIX")
	}
	return fmt.Sprintf("%s-%s", prefix, makeCameraName(camera))
}

func makeKeyName(camera, id string) string {
	return fmt.Sprintf("%s-%s", makeCameraName(camera), id)
}

// To make sure that the bucket name is unique, we remove spaces and convert to lowercase.
func makeCameraName(name string) string {
	lowercase := strings.ToLower(name)
	return strings.ReplaceAll(lowercase, " ", "")
}
