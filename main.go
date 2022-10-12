package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/admin"
	"github.com/EthanDamien/glimpse-go-aws-lambda/database"
	"github.com/EthanDamien/glimpse-go-aws-lambda/s3"
	"github.com/EthanDamien/glimpse-go-aws-lambda/shift"
	"github.com/EthanDamien/glimpse-go-aws-lambda/user"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type HandleResponse struct {
	OK    bool   `json:"ok"`
	ReqID string `json:"req_id"`
}

type HandleRequest struct {
	Event string          `json:"event"`
	Body  json.RawMessage `json:"body"`
}

var logger *zap.Logger
var db *sql.DB
var awsSession *session.Session

// This function initializes the database connection
func initDatabaseConnection() {
	l, _ := zap.NewProduction()
	logger = l
	logger.Info("Getting DB connection")

	dbConnection, err := database.GetConnection()
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		panic(err)
	}

	logger.Info("Pinging Database")
	err = dbConnection.Ping()
	if err != nil {
		logger.Error("error pinging database", zap.Error(err))
		panic(err)
	}

	// Set global var
	db = dbConnection
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

	//Initialize Database
	initDatabaseConnection()

	//Connect to s3
	awsSession := s3.ConnectAws()
	logger.Info(*awsSession.Config.Region)

	//This is the first row in the json request and will do certain things based on this variable
	switch req.Event {
	case "createUser":
		var dest user.CreateUserRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return user.CreateUser(ctx, reqID, dest)
	case "createAdmin":
		var dest admin.CreateAdminRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return admin.CreateAdmin(ctx, reqID, dest, db)
	case "createShift":
		var dest shift.CreateShiftRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return shift.CreateShift(ctx, reqID, dest, db)
	}

	return HandleResponse{OK: false, ReqID: reqID}, fmt.Errorf("%s is an unknown event", req.Event)
}

func main() {
	lambda.Start(Handle)
}
