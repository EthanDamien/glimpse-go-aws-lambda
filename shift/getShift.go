package shift

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type GetAllShiftsRequest struct {
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

type GetEmployeeShiftsRequest struct {
	EmployeeID int `json:"employeeID"`
}

type GetShiftRequest struct {
	ShiftEventID int `json:"shiftEventID"`
}

const getAllShiftsTemplate = `
SELECT * FROM Shift WHERE ClockInTime >= "%s" AND ClockOutTime <= "%s");`

const getEmployeeShiftsTemplate = `
SELECT * FROM Shift WHERE EmployeeID == "%d");`

const getShiftTemplate = `
SELECT * FROM Shift WHERE ShiftEventID == "%d");`

func GetAllShifts(ctx context.Context, reqID string, req GetAllShiftsRequest, db *sql.DB) (sql.Result, error) {

	//validate JSON
	if req.FromDate.IsZero() {
		return nil, fmt.Errorf("Missing FromDate")
	}
	if req.ToDate.IsZero() {
		return nil, fmt.Errorf("Missing ToDate")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(getAllShiftsTemplate, req.FromDate, req.ToDate)
	response, err := db.ExecContext(ctx, builtQuery)

	//If this fails, send "error" response
	//TODO send actual error to Lambda
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetEmployeeShifts(ctx context.Context, reqID string, req GetEmployeeShiftsRequest, db *sql.DB) (sql.Result, error) {

	//validate JSON
	if req.EmployeeID == 0 {
		return nil, fmt.Errorf("Missing EmployeeID")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(getEmployeeShiftsTemplate, req.EmployeeID)
	response, err := db.ExecContext(ctx, builtQuery)

	//If this fails, send "error" response
	//TODO send actual error to Lambda
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetShift(ctx context.Context, reqID string, req GetShiftRequest, db *sql.DB) (sql.Result, error) {

	//validate JSON
	if req.ShiftEventID == 0 {
		return nil, fmt.Errorf("Missing ShiftEventID")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(getShiftTemplate, req.ShiftEventID)
	response, err := db.ExecContext(ctx, builtQuery)

	//If this fails, send "error" response
	//TODO send actual error to Lambda
	if err != nil {
		return nil, err
	}
	return response, nil
}
