package s3

import (
	"fmt"
	// "os"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadObject(AWS_S3_BUCKET string, w http.ResponseWriter, r *http.Request, sess *session.Session) error {

	/*
		// Open file to upload
		file, err := os.Open(filePath)
		if err != nil {
			logger.Error("Unable to open file %q, %v", zap.Error(err))
			return err
		}
		defer file.Close()
	*/
	r.ParseMultipartForm(10 << 20)

	// From form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Something went wrong retrieving the file from the form", http.StatusInternalServerError)
		return err
	}
	defer file.Close()

	fileName := header.Filename

	// Upload to s3
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		http.Error(w, "Something went wrong uploading the file", http.StatusInternalServerError)
		return err
	}

	fmt.Printf("Successfully uploaded %q to %q\n", fileName, AWS_S3_BUCKET)
	return nil
}
