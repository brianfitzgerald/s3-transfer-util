package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	sourceBucket := flag.String("source-bucket", "", "bucket to copy")
	destBucket := flag.String("dest-bucket", "", "bucket to copy to")

	sourceProfile := flag.String("source-creds", "", "creds for copying from the source bucket")
	destProfile := flag.String("source-creds", "", "creds for copying to the destination bucket")

	if *sourceBucket == "" || *destBucket == "" {
		panic("Please specify source and destination buckets")
	}

	if *sourceProfile == "" || *destProfile == "" {
		panic("Please specify source and destination profiles")
	}

	sourceBucketSession, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", *sourceProfile),
	})

	sourceBucketClient := s3.New(sourceBucketSession)

	sourceBucketObjectsList := &s3.ListObjectsInput{
		Bucket: sourceBucket,
	}

	result, err := sourceBucketClient.ListObjects(sourceBucketObjectsList)

	if err != nil {
		panic(err)
	}

	destBucketSession, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", *sourceProfile),
	})

	destBucketClient := s3.New(destBucketSession)

	for _, item := range result.Contents {
		copyOutput, err := destBucketClient.CopyObject(&s3.CopyObjectInput{
			Bucket: destBucket,
			Key:    item.Key,
		})
		println(copyOutput.String)
		if err != nil {
			log.Fatal(err)
		}
	}

	println("Copy complete.")

}
