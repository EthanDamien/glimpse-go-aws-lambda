package wage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// request format for create wage
type CreateWageRequest struct {
	EmployeeID  int       `json:"EmployeeID"`
	WagePerHour float64   `json:"WagePerHour"`
	TimeToSet   time.Time `json:"TimeToSet"`
}

// response format for create wage
type CreateWageResponse struct {
	DESC string `json:"desc"`
	OK   bool   `json:"ok"`
}

const createWageTemplate = `Insert into Wage (WageEventID, EmployeeID, WagePerHour, TimeToSet) values (NULL, %d, %f, "%s"); `

// Create wage for an employee
func CreateWage(ctx context.Context, reqID string, req CreateWageRequest, db *sql.DB) (CreateWageResponse, error) {
	if req.EmployeeID == 0 {
		return CreateWageResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing EmployeeID")
	}
	if req.WagePerHour == 0 {
		return CreateWageResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing WagePerHour")
	}
	if req.TimeToSet.IsZero() {
		return CreateWageResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing TimeToSet")
	}

	var builtQuery = fmt.Sprintf(createWageTemplate, req.EmployeeID, req.WagePerHour, req.TimeToSet)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return CreateWageResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing Password")
	}

	return CreateWageResponse{DESC: fmt.Sprintf("Wage Created with values EmployeeID: %d, Wage %f, TimeToSet %s",
		req.EmployeeID, req.WagePerHour, req.TimeToSet), OK: true}, nil
}
