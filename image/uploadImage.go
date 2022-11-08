package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadImage(base64Meta string, location string, bucket string) error {
	base64data := base64Meta[strings.IndexByte(base64Meta, ',')+1:]
	decodedImage, err := base64.StdEncoding.DecodeString(base64data)

	if err != nil {
		return fmt.Errorf(err.Error())
	}
	sess := session.New()

	//To ensure that only one picture is added for each employee, delete, for each picture upload.
	deleteObjectForUpload(sess, location)

	uploader := s3manager.NewUploader(sess)
	uploadParameters := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s.jpg", location)),
		Body:   bytes.NewReader(decodedImage),
	}
	_, err = uploader.Upload(uploadParameters)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func deleteObjectForUpload(sess *session.Session, location string) {
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
