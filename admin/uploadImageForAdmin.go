package admin

import (
	"context"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/image"
	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// request format for uploading the picture
type UploadPictureRequest struct {
	AdminID       string `json:"AdminID"`
	EmployeeID    string `json:"EmployeeID"`
	PictureMeta64 string `json:"PictureMeta64"`
}

// This method uploads the picture
func UploadImageForAdmin(ctx context.Context, reqID string, req UploadPictureRequest) (AdminResponse, error) {
	if req.AdminID == "" {
		return AdminResponse{}, fmt.Errorf(statuscode.C500, "Missing AdminID")
	}
	if req.EmployeeID == "" {
		return AdminResponse{}, fmt.Errorf(statuscode.C500, "Missing EmployeeID")
	}
	if req.PictureMeta64 == "" {
		return AdminResponse{}, fmt.Errorf(statuscode.C500, "Missing PictureMetadata")
	}

	err := image.UploadImage(req.PictureMeta64, req.EmployeeID, "facefiles")

	if err != nil {
		return AdminResponse{}, fmt.Errorf(err.Error())
	}
	return AdminResponse{DESC: "Upload Picture Success"}, nil
}
