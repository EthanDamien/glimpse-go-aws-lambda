package compare

import (
	"bytes"
	"context"
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

type FindMatchingReq struct {
	AdminID       string `json:"AdminID"`
	PictureMeta64 string `json:"PictureMeta64"`
}

// Query AdminID
const findEmployeesPerAdminIDTemplate = `SELECT EmployeeID from Employees where AdminID = "%s";`

// Using json input, find a matching employee for the given adminID, and image.
func FindMatchingEmployee(ctx context.Context, reqID string, req FindMatchingReq, db *sql.DB) (CompareResponse, error) {
	if req.AdminID == "" {
		return CompareResponse{DESC: "Find Matching Employee Error"}, fmt.Errorf("Missing AdminID")
	}
	if req.PictureMeta64 == "" {
		return CompareResponse{DESC: "Find Matching Employee Error"}, fmt.Errorf("Missing PictureMetadata")
	}

	//create random image name and upload it
	tempImageNum := rand.Intn(10000)
	var tempImageName = strconv.Itoa(tempImageNum)

	image.UploadImage(req.PictureMeta64, tempImageName, "facefiles")

	//Do the query
	var builtQuery = fmt.Sprintf(findEmployeesPerAdminIDTemplate, req.AdminID)
	employeeIDs, err := getQueryRes(builtQuery, db)
	if err != nil {
		// image.DeleteImage(tempImageName)
		return CompareResponse{DESC: "Find Matching Query Err"}, err
	}
	for _, id := range employeeIDs {
		var tempImageNameLoc = fmt.Sprintf("%s.jpg", tempImageName)
		var idAsNameLoc = fmt.Sprintf("%s.jpg", strconv.Itoa(id.EmployeeID))
		log.Printf("Checking %s", idAsNameLoc)
		isMatch, err, _ := Compare(idAsNameLoc, tempImageNameLoc)
		if err != nil {
			// image.DeleteImage(tempImageName)
			return CompareResponse{DESC: "Compare Err"}, err
		}
		if isMatch {
			// image.DeleteImage(tempImageName)
			return CompareResponse{DESC: strconv.Itoa(id.EmployeeID)}, nil
		}
	}

	// image.DeleteImage(tempImageName)
	return CompareResponse{DESC: "Employee not found"}, nil
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