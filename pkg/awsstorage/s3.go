package awsstorage

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3cli handles S3 operations
type S3cli struct {
	session client.ConfigProvider
}

// Upload sent a io.Reader attachment to a S3 bucket with bucket and key name
func (s S3cli) Upload(bucket, key string, file io.Reader) (err error) {
	uploader := s3manager.NewUploader(s.session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return
}

// List is not complete implemented yet
func (s S3cli) List(bucket, _ string) error {
	svc := s3.New(s.session)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		return err
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	return err
}

// Get is not complete implemented yet
func (s S3cli) Get(bucket, key string) error {
	file, err := os.Create(key)
	if err != nil {
		return err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(s.session)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return err
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return nil
}
