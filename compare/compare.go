package compare

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
)

// struct for comparison request
type CompareReq struct {
	Loc1 string `json:"loc1"`
	Loc2 string `json:"loc2"`
}

// this compare function returns true if the picture matches, if not, returns false.
// All other errors also return false. The CompareFacesOutput is also returned for visibility.
// return a boolean, an error, and a CompareFacesOutput instance
func Compare(location1 string, location2 string) (bool, error, rekognition.CompareFacesOutput) {

	sess := session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	//Validate if location is an existing object
	keyExistsLoc1, err := keyExists("facefiles", location1, sess)

	if err != nil {
		return false, fmt.Errorf("Failure in KeyExists"), rekognition.CompareFacesOutput{}
	}

	keyExistsLoc2, err := keyExists("facefiles", location2, sess)

	if err != nil {
		return false, fmt.Errorf("Failure in KeyExists"), rekognition.CompareFacesOutput{}
	}

	if !keyExistsLoc1 || !keyExistsLoc2 {
		log.Printf("Appearance does not exist for: %s", location1)
		return false, nil, rekognition.CompareFacesOutput{}
	}

	svc := rekognition.New(sess)

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

// checks to see if the key exists in S3
// return true (found) or false (not found) if query performed successfully, else error
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
