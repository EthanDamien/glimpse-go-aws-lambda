package shift

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ClockingActionRequest struct {
	ClockTime time.Time `json:"clockTime"`
	Earnings  float32   `json:"earnings"`
}

type ClockingActionResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const getEmployeeMostRecentShiftTemplate = `SELECT * FROM Shift WHERE EmployeeID = %d ORDER BY LastUpdated DESC LIMIT 1;`

func ClockingAction(ctx context.Context, reqID string, req ClockingActionRequest, employeeID int, db *sql.DB) (ClockingActionResponse, error) {
	// most recent shift
	var builtQuery = fmt.Sprintf(getEmployeeMostRecentShiftTemplate, employeeID)
	res, err := getQueryRes(builtQuery, db)
	if err != nil {
		return ClockingActionResponse{DESC: "Could not get most recent shift", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}
	if len(res) == 0 || (!res[0].ClockInTime.IsZero() && !res[0].ClockOutTime.IsZero()) {
		_, csErr := CreateShift(ctx, reqID, CreateShiftRequest{EmployeeID: employeeID, ClockInTime: req.ClockTime, ClockOutTime: time.Time{}, Earnings: req.Earnings}, db)
		if csErr != nil {
			return ClockingActionResponse{DESC: "Could not create shift", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, csErr
		}
	} else {
		_, usErr := UpdateShift(ctx, reqID, UpdateShiftRequest{ShiftEventID: res[0].ShiftEventID, EmployeeID: employeeID, ClockInTime: res[0].ClockInTime, ClockOutTime: req.ClockTime, Earnings: req.Earnings}, db)
		if usErr != nil {
			return ClockingActionResponse{DESC: "Could not update shift", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, usErr
		}
	}
	return ClockingActionResponse{DESC: "Clock Action successful", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
