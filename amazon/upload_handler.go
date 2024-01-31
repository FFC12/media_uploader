package amazon

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func s3StreamUpload(bucket, region, filename string, data []byte) error {
	// create a reader from data in memory
	reader := bytes.NewReader(data)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		// the aws sdk will manage all the memory and file reading for us
		Body: reader,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
	return nil
}
