package image

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// this deletes an image with the location name in s3
func DeleteImage(location string) {
	sess := session.New()
	svc := s3.New(sess)

	deleteParameters := &s3.DeleteObjectInput{
		Bucket: aws.String("facefiles"),
		Key:    aws.String(location),
	}
	result, err := svc.DeleteObject(deleteParameters)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	fmt.Printf(result.GoString())
}
