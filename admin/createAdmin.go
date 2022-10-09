package admin

import (
	"context"
	"fmt"
	"time"
)

type CreateAdminRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateAdminResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

func CreateAdmin(_ context.Context, reqID string, req CreateAdminRequest) (CreateAdminResponse, error) {
	if req.FirstName == "" {
		return CreateAdminResponse{DESC: "CreateAdmin err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("the first_name is missing")
	}
	if req.LastName == "" {
		return CreateAdminResponse{DESC: "CreateAdmin err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("the last_name is missing")
	}

	return CreateAdminResponse{DESC: "CreateAdmin success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
