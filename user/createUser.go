package user

import (
	"context"
	"fmt"
	"time"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateUserResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

func CreateUser(_ context.Context, reqID string, req CreateUserRequest) (CreateUserResponse, error) {
	if req.FirstName == "" {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("the first_name is missing")
	}
	if req.LastName == "" {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("the last_name is missing")
	}

	return CreateUserResponse{DESC: "CreateUser success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
