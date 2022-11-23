package admin

import (
	"context"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/image"
)

type UploadPictureRequest struct {
	AdminID       string `json:"AdminID"`
	EmployeeID    string `json:"EmployeeID"`
	PictureMeta64 string `json:"PictureMeta64"`
}

// This method uploads the picture
func UploadImageForAdmin(ctx context.Context, reqID string, req UploadPictureRequest) (AdminResponse, error) {
	if req.AdminID == "" {
		return AdminResponse{DESC: "Upload Picture Error"}, fmt.Errorf("Status:500 Missing AdminID")
	}
	if req.EmployeeID == "" {
		return AdminResponse{DESC: "Upload Picture Error"}, fmt.Errorf("Status:500 Missing EmployeeID")
	}
	if req.PictureMeta64 == "" {
		return AdminResponse{DESC: "Upload Picture Error"}, fmt.Errorf("Status:500 Missing PictureMetadata")
	}

	err := image.UploadImage(req.PictureMeta64, req.EmployeeID, "facefiles")

	if err != nil {
		return AdminResponse{DESC: "Upload Picture Error"}, fmt.Errorf(err.Error())
	}
	return AdminResponse{DESC: "Upload Picture Success"}, nil
}
