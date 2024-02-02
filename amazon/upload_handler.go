package amazon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Constants for AWS configuration
const (
	maxRetries         = 3
	awsAccessKeyID     = "abe15e537dac570be4d21b97ba40c3c2"
	awsSecretAccessKey = "54f96d2b27fd88abfbc7f0a5b07ce3f661bb30a8987f39f21d69940812a31676"
	awsBucketRegion    = "eu-east-1"
	awsBucketName      = "storage"
	awsAccountId       = "aaa780ca2d934ac0f129acd5a54e5c39"
	awsURL             = "https://aaa780ca2d934ac0f129acd5a54e5c39.r2.cloudflarestorage.com/storage"
)

// StreamUploadInit initializes the multipart upload process
func StreamUploadInit(context *context.Context, mimeType, filename string) (*s3.Client, *s3.CreateMultipartUploadOutput, error) {
	// Custom Endpoint Resolver for Cloudflare
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", awsAccountId),
		}, nil
	})

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(*context,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create S3 client
	svc := s3.NewFromConfig(cfg)

	// Set up parameters for multipart upload initialization
	path := filename
	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(awsBucketName),
		Key:         aws.String(path),
		ContentType: aws.String(mimeType),
	}

	// Initiate multipart upload
	resp, err := svc.CreateMultipartUpload(*context, input)
	if err != nil {
		return svc, resp, err
	}

	fmt.Println("Created multipart upload request")
	return svc, resp, nil
}

// StreamUpload uploads a part of the file in the multipart upload process
func StreamUpload(context *context.Context, svc *s3.Client, resp *s3.CreateMultipartUploadOutput, buffer []byte, partNumber int32) (*s3.UploadPartOutput, error) {
	var uploadResult *s3.UploadPartOutput
	var err error
	tries := 0

	// Retry loop for uploading a part
	for i := 0; i < maxRetries; i++ {
		partInput := &s3.UploadPartInput{
			Body:       bytes.NewReader(buffer),
			Bucket:     resp.Bucket,
			Key:        resp.Key,
			PartNumber: &partNumber,
			UploadId:   resp.UploadId,
		}
		uploadResult, err = svc.UploadPart(*context, partInput)
		if err != nil {
			// Retry if unsuccessful
			if tries < maxRetries {
				tries++
				continue
			}

			// Abort multipart upload in case of repeated failures
			aboInput := &s3.AbortMultipartUploadInput{
				Bucket:   resp.Bucket,
				Key:      resp.Key,
				UploadId: resp.UploadId,
			}
			_, aboErr := svc.AbortMultipartUpload(*context, aboInput)
			if aboErr != nil {
				println("Error(aborted): ", aboErr)
				return nil, aboErr
			}
			return nil, err
		}
	}

	fmt.Println("Uploaded part part number:", partNumber, "etag:", *uploadResult.ETag)

	return uploadResult, nil
}

// StreamDone completes the multipart upload process
func StreamDone(context *context.Context, svc *s3.Client, resp *s3.CreateMultipartUploadOutput, completedParts []types.CompletedPart) (string, error) {
	// Set up parameters for completing multipart upload
	compInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}

	// Complete multipart upload
	output, compErr := svc.CompleteMultipartUpload(*context, compInput)
	if compErr != nil {
		fmt.Println("Error: (completed parts= ", len(completedParts), "): ", compErr)
		return "", compErr
	}

	// Print JSON output
	json, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	fmt.Println(string(json))

	return *output.Location, nil
}

// DirectUpload uploads an object directly without using multipart upload
func DirectUpload(context *context.Context, mimeType, filename string, buffer []byte) (string, error) {
	// Custom Endpoint Resolver for Cloudflare
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", awsAccountId),
		}, nil
	})

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(*context,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create S3 client
	svc := s3.NewFromConfig(cfg)

	// Set up parameters for direct object upload
	path := filename
	input := &s3.PutObjectInput{
		Bucket:      aws.String(awsBucketName),
		Key:         aws.String(path),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(mimeType),
	}

	// Upload object directly
	resp, err := svc.PutObject(*context, input)
	if err != nil {
		return "", err
	}

	fmt.Println("Uploaded object successfully: ", resp)

	return awsURL + "/" + path, nil
}

// getUrlOfUploadedObject gets the URL of the uploaded object
func getUrlOfUploadedObject(context *context.Context, svc *s3.Client, awsBucketName, filename string) (string, error) {
	// List and get the URL of uploaded object
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(awsBucketName),
	}

	listObjectsOutput, err := svc.ListObjectsV2(*context, listObjectsInput)
	if err != nil {
		return "", err
	}

	fmt.Println("URLs of uploaded objects:")
	var url string
	for _, object := range listObjectsOutput.Contents {
		if *object.Key == filename {
			url = fmt.Sprintf("%s/%s", awsURL, *object.Key)
			fmt.Println(url)
			break
		}
	}
	return url, nil
}
