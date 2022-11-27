package wage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CreateWageRequest struct {
	EmployeeID  int       `json:"EmployeeID"`
	WagePerHour float64   `json:"WagePerHour"`
	TimeToSet   time.Time `json:"TimeToSet"`
}

type CreateWageResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const createWageTemplate = `Insert into Wage (WageEventID, EmployeeID, WagePerHour, TimeToSet) values (NULL, %d, %f, "%s"); `

func CreateWage(ctx context.Context, reqID string, req CreateWageRequest, db *sql.DB) (CreateWageResponse, error) {
	if req.EmployeeID == 0 {
		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing EmployeeID")
	}
	if req.WagePerHour == 0 {
		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing WagePerHour")
	}
	if req.TimeToSet.IsZero() {
		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing TimeToSet")
	}

	var builtQuery = fmt.Sprintf(createWageTemplate, req.EmployeeID, req.WagePerHour, req.TimeToSet)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing Password")
	}

	return CreateWageResponse{DESC: fmt.Sprintf("Wage Created with values EmployeeID: %d, Wage %f, TimeToSet %s",
		req.EmployeeID, req.WagePerHour, req.TimeToSet)}, nil
}
