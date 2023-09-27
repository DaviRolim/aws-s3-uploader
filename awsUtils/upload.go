package awsUtils

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToS3(bucketName string, fileHeader *multipart.FileHeader) error {
	// Open the image file
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new session with default session credentials
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(cfg))

	// S3 service client the Upload manager will use.
	svc := s3.New(sess)

	// Configure the S3 object input parameters
	input := &s3.PutObjectInput{
		Body:          file,
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileHeader.Filename),
		ContentType:   aws.String(fileHeader.Header["Content-Type"][0]),
		ContentLength: aws.Int64(fileHeader.Size),
	}

	fmt.Println("Created input")

	// Upload the image to S3
	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}

	fmt.Println("Uploaded image to s3")

	return nil
}
