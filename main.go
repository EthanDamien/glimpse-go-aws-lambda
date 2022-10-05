package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type HandleResponse struct {
	OK    bool   `json:"ok"`
	ReqID string `json:"req_id"`
}

type HandleRequest struct {
	Event string          `json:"event"`
	Body  json.RawMessage `json:"body"`
}

// Handle the calls
func Handle(ctx context.Context, req HandleRequest) (interface{}, error) {
	var reqID string
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		reqID = lc.AwsRequestID
	}

	select {
	case <-ctx.Done():
		return HandleResponse{OK: false, ReqID: reqID}, fmt.Errorf("request timeout: %w", ctx.Err())
	default:
	}

	switch req.Event {
	case "createUser":
		var dest CreateUserRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return CreateUser(ctx, reqID, dest)
	case "createAdmin":
		var dest CreateAdminRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return CreateAdmin(ctx, reqID, dest)
	}

	return HandleResponse{OK: false, ReqID: reqID}, fmt.Errorf("%s is an unknown event", req.Event)
}

func main() {
	lambda.Start(Handle)
}
