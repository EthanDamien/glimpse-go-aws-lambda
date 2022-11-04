package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/admin"
	"github.com/EthanDamien/glimpse-go-aws-lambda/database"
	"github.com/EthanDamien/glimpse-go-aws-lambda/login"
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
	// EMPLOYEE
	case "createUser":
		var dest user.CreateUserRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return user.CreateUser(ctx, reqID, dest, db)
	case "getAllUsers":
		var dest user.GetAllUsersRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return user.GetAllUsers(ctx, reqID, dest, db)
	case "employeeLogin":
		var dest login.EmployeeLoginRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return login.EmployeeLogin(ctx, reqID, dest, db)
	case "createAdmin":
		var dest admin.CreateAdminRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return admin.CreateAdmin(ctx, reqID, dest, db)
	case "getAdmin":
		var dest admin.GetAdminRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return admin.GetAdmin(ctx, reqID, dest, db)
	case "getAdminByAdminID":
		var dest admin.GetAdminByAdminIDRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return admin.GetAdminByAdminID(ctx, reqID, dest, db)
	case "adminLogin":
		var dest login.AdminLoginRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return login.AdminLogin(ctx, reqID, dest, db)
	// SHIFT
	case "clockingAction":
		var dest shift.ClockingActionRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		// constant employeeID for now
		return shift.ClockingAction(ctx, reqID, dest, 1, db)
	case "createShift":
		var dest shift.CreateShiftRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return shift.CreateShift(ctx, reqID, dest, db)
	case "updateShift":
		var dest shift.UpdateShiftRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return shift.UpdateShift(ctx, reqID, dest, db)
	case "getMostRecentShifts":
		var dest shift.GetMostRecentShiftsRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return shift.GetMostRecentShifts(ctx, reqID, dest, db)
	case "getAllShifts":
		var dest shift.GetAllShiftsRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return shift.GetAllShifts(ctx, reqID, dest, db)
	case "getShift":
		var dest shift.GetShiftRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return shift.GetShift(ctx, reqID, dest, db)
	case "getEmployeeShifts":
		var dest shift.GetEmployeeShiftsRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return shift.GetEmployeeShifts(ctx, reqID, dest, db)
	}

	return HandleResponse{OK: false, ReqID: reqID}, fmt.Errorf("%s is an unknown event", req.Event)
}

func main() {
	lambda.Start(Handle)
}
