package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// credentials
const (
	AWS_S3_REGION         = "us-east-1" // Region
	AWS_S3_BUCKET         = "facefiles" // Bucket
	AWS_SECRET_ACCESS_KEY = "Dm5/nSIyg79NgkgTSxv45DBVieT6vM+gLm8JScd3"
	AWS_ACCESS_KEY_ID     = "AKIA42VSILKBI7OAOK4N"
)

// connect to AWS
func ConnectAws() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(AWS_S3_REGION),
			Credentials: credentials.NewStaticCredentials(
				AWS_ACCESS_KEY_ID,
				AWS_SECRET_ACCESS_KEY,
				"", // a token will be created when the session it's used.
			),
		})
	fmt.Print("Connecting to AWS S3")
	if err != nil {
		fmt.Print("PINGED WITH ERR")
		panic(err)
	}
	return sess
}
