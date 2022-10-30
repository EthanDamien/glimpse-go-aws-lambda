package reko

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

func CreateRekognitionServiceInstance() *rekognition.Rekognition {
	sess := session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := rekognition.New(sess)
	return svc
}
