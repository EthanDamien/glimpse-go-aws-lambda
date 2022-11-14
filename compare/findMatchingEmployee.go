package compare

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/EthanDamien/glimpse-go-aws-lambda/image"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type EmployeeID struct {
	EmployeeID int `json:"employeeID"`
}

// Query AdminID
const findEmployeesPerAdminIDTemplate = `SELECT EmployeeID from Employees where AdminID = "%s";`

// Using json input, find a matching employee for the given adminID, and image.
func FindMatchingEmployee(AdminId string, PictureMeta64 string, db *sql.DB) (int, error) {

	//create random image name and upload it
	tempImageNum := rand.Intn(10000)
	var tempImageName = strconv.Itoa(tempImageNum)

	image.UploadImage(PictureMeta64, tempImageName, "facefiles")
	var tempImageNameLoc = fmt.Sprintf("%s.jpg", tempImageName)
	//Do the query
	var builtQuery = fmt.Sprintf(findEmployeesPerAdminIDTemplate, AdminId)
	employeeIDs, err := getQueryRes(builtQuery, db)
	if err != nil {
		// image.DeleteImage(tempImageName)
		return 0, err
	}
	for _, id := range employeeIDs {
		var idAsNameLoc = fmt.Sprintf("%s.jpg", strconv.Itoa(id.EmployeeID))
		log.Printf("Checking %s", idAsNameLoc)
		isMatch, err, _ := Compare(idAsNameLoc, tempImageNameLoc)
		if err != nil {
			image.DeleteImage(tempImageNameLoc)
			return 0, err
		}
		if isMatch {
			image.DeleteImage(tempImageNameLoc)
			return id.EmployeeID, nil
		}
	}

	image.DeleteImage(tempImageNameLoc)
	return 0, fmt.Errorf("Employee Not Found")
}

func getQueryRes(builtQuery string, db *sql.DB) ([]EmployeeID, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employeeIDs []EmployeeID

	for rows.Next() {
		var employeeID EmployeeID
		if err := rows.Scan(&employeeID.EmployeeID); err != nil {
			return employeeIDs, err
		}
		employeeIDs = append(employeeIDs, employeeID)
	}
	return employeeIDs, nil
}

func uploadPicture(base64Meta string, location string, bucket string) (bool, error) {
	base64data := base64Meta[strings.IndexByte(base64Meta, ',')+1:]
	decodedImage, err := base64.StdEncoding.DecodeString(base64data)

	if err != nil {
		return false, fmt.Errorf(err.Error())
	}
	sess := session.New()

	//To ensure that only one picture is added for each employee, delete, for each picture upload.
	deleteObject(sess, location)

	uploader := s3manager.NewUploader(sess)
	uploadParameters := &s3manager.UploadInput{
		Bucket: aws.String("facefiles"),
		Key:    aws.String(fmt.Sprintf("%s.jpg", location)),
		Body:   bytes.NewReader(decodedImage),
	}
	_, err = uploader.Upload(uploadParameters)

	if err != nil {
		return false, fmt.Errorf(err.Error())
	}

	return true, nil
}

func deleteObject(sess *session.Session, location string) {
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
