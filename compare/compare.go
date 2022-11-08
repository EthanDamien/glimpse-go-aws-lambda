package compare

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
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

// this compare function returns true if the picture matches, if not, returns false.
// All other errors also return false. The CompareFacesOutput is also returned for visibility.
func Compare(location1 string, location2 string) (bool, error, rekognition.CompareFacesOutput) {

	sess := session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	//Validate if location is an existing object
	keyExists, err := keyExists("facefiles", location1, sess)

	if err != nil {
		return false, fmt.Errorf("Failure in KeyExists"), rekognition.CompareFacesOutput{}
	}

	if !keyExists {
		log.Printf("Appearance does not exist for: %s", location1)
		return false, nil, rekognition.CompareFacesOutput{}
	}

	svc := rekognition.New(sess)

	if err != nil {
		return false, fmt.Errorf("Failure Decoding Image"), rekognition.CompareFacesOutput{}
	}

	log.Printf("Comparing %s to %s", location1, location2)
	input := &rekognition.CompareFacesInput{
		SimilarityThreshold: aws.Float64(99),
		SourceImage: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String("facefiles"),
				Name:   aws.String(location1),
			},
		},
		TargetImage: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String("facefiles"),
				Name:   aws.String(location2),
			},
		},
	}

	result, err := svc.CompareFaces(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidParameterException:
				return false, fmt.Errorf(rekognition.ErrCodeInvalidParameterException), rekognition.CompareFacesOutput{}
			case rekognition.ErrCodeInvalidS3ObjectException:
				return false, fmt.Errorf(rekognition.ErrCodeInvalidS3ObjectException), rekognition.CompareFacesOutput{}
			case rekognition.ErrCodeImageTooLargeException:
				return false, fmt.Errorf(rekognition.ErrCodeImageTooLargeException), rekognition.CompareFacesOutput{}
			case rekognition.ErrCodeAccessDeniedException:
				return false, fmt.Errorf(rekognition.ErrCodeAccessDeniedException), rekognition.CompareFacesOutput{}
			case rekognition.ErrCodeInternalServerError:
				return false, fmt.Errorf(rekognition.ErrCodeInternalServerError), rekognition.CompareFacesOutput{}
			case rekognition.ErrCodeThrottlingException:
				return false, fmt.Errorf(rekognition.ErrCodeThrottlingException), rekognition.CompareFacesOutput{}
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				return false, fmt.Errorf(rekognition.ErrCodeProvisionedThroughputExceededException), rekognition.CompareFacesOutput{}
			case rekognition.ErrCodeInvalidImageFormatException:
				return false, fmt.Errorf(rekognition.ErrCodeInvalidImageFormatException), rekognition.CompareFacesOutput{}
			default:
				return false, fmt.Errorf(aerr.Error()), rekognition.CompareFacesOutput{}
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return false, fmt.Errorf(err.Error()), rekognition.CompareFacesOutput{}
		}
	}
	if len(result.FaceMatches) == 0 {
		return false, nil, *result
	}
	fmt.Println(result)
	return true, nil, *result
}

func keyExists(bucket string, key string, sess *session.Session) (bool, error) {
	svc := s3.New(sess)

	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}
