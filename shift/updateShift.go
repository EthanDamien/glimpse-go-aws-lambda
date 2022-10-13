package shift

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UpdateShiftRequest struct {
	ShiftEventID int       `json:"shiftEventID"`
	EmployeeID   int       `json:"employeeID"`
	ClockInTime  time.Time `json:"clockInTime"`
	ClockOutTime time.Time `json:"clockOutTime"`
	Earnings     float32   `json:"earnings"`
}

type UpdateShiftResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

// Template to be used to insert to Table
const updateShiftTemplate = `
UPDATE Shift SET ClockInTime="%s", ClockOutTime="%s", Earnings="%f" where ShiftEventID = %d;`

func UpdateShift(ctx context.Context, reqID string, req UpdateShiftRequest, db *sql.DB) (UpdateShiftResponse, error) {

	//validate JSON
	if req.ShiftEventID == 0 {
		return UpdateShiftResponse{DESC: "UpdateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing ShiftEventID")
	}
	if req.ClockInTime.IsZero() {
		return UpdateShiftResponse{DESC: "UpdateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing ClockInTime")
	}
	if req.ClockOutTime.IsZero() {
		return UpdateShiftResponse{DESC: "UpdateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing ClockOutTime")
	}
	if req.ClockInTime.After(req.ClockOutTime) {
		return UpdateShiftResponse{DESC: "UpdateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("ClockInTime must be earlier than ClockOutTime")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(updateShiftTemplate, req.ClockInTime, req.ClockOutTime, req.Earnings, req.ShiftEventID)
	_, err := db.ExecContext(ctx, builtQuery)

	//If this fails, send "error" response
	//TODO send actual error to Lambda
	if err != nil {
		return UpdateShiftResponse{DESC: "Could not update Shift Table", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}
	return UpdateShiftResponse{DESC: "UpdateShift success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
