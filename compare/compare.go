package compare

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type CompareReq struct {
	Loc1 string `json:"loc1"`
	Loc2 string `json:"loc2"`
}

func TestCompare(ctx context.Context, reqID string, req CompareReq, db *sql.DB) (CompareResponse, error) {
	//validate JSON
	if req.Loc1 == "" {
		return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf("Missing Location 1")
	}
	if req.Loc2 == "" {
		return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf("Missing Location 2")
	}

	sess := session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := rekognition.New(sess)

	input := &rekognition.CompareFacesInput{
		SimilarityThreshold: aws.Float64(99),
		SourceImage: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String("facefiles"),
				Name:   aws.String(req.Loc1),
			},
		},
		TargetImage: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String("facefiles"),
				Name:   aws.String(req.Loc2),
			},
		},
	}

	result, err := svc.CompareFaces(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidParameterException:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(rekognition.ErrCodeInvalidParameterException)
			case rekognition.ErrCodeInvalidS3ObjectException:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(rekognition.ErrCodeInvalidS3ObjectException)
			case rekognition.ErrCodeImageTooLargeException:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(rekognition.ErrCodeImageTooLargeException)
			case rekognition.ErrCodeAccessDeniedException:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(rekognition.ErrCodeAccessDeniedException)
			case rekognition.ErrCodeInternalServerError:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(rekognition.ErrCodeInternalServerError)
			case rekognition.ErrCodeThrottlingException:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(rekognition.ErrCodeThrottlingException)
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(rekognition.ErrCodeProvisionedThroughputExceededException)
			case rekognition.ErrCodeInvalidImageFormatException:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(rekognition.ErrCodeInvalidImageFormatException)
			default:
				return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return CompareResponse{DESC: "CompareFaces err"}, fmt.Errorf(err.Error())
		}
	}
	fmt.Println(result)
	return CompareResponse{DESC: "Face Match success"}, nil
}
