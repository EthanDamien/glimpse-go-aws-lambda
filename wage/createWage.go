package wage

import (
	"context"
	"database/sql"
	"fmt"
)

type CreateWageRequest struct {
	EmployeeID  string `json:"EmployeeID"`
	WagePerHour string `json:"WagePerHour"`
	TimeToSet   string `json:"TimeToSet"`
}

type CreateWageResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const createWageTemplate = `Insert into Wage (WageEventID, EmployeeID, WagePerHour, TimeToSet) values (NULL, %s, %s, "%s"); `

func CreateWage(ctx context.Context, reqID string, req CreateWageRequest, db *sql.DB) (CreateWageResponse, error) {
	if req.EmployeeID == "" {
		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing EmployeeID")
	}
	if req.WagePerHour == "" {
		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing WagePerHour")
	}
	if req.TimeToSet == "" {
		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing TimeToSet")
	}

	var builtQuery = fmt.Sprintf(createWageTemplate, req.EmployeeID, req.WagePerHour, req.TimeToSet)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing Password")
	}

	return CreateWageResponse{DESC: fmt.Sprintf("Wage Created with values EmployeeID: %s, Wage %s, TimeToSet %s",
		req.EmployeeID, req.WagePerHour, req.TimeToSet)}, nil
}
