package shift

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
	"github.com/EthanDamien/glimpse-go-aws-lambda/wage"
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

const updateShiftTemplate = `
UPDATE Shift SET ClockInTime="%s", ClockOutTime="%s", Earnings="%f", LastUpdated="%s" where ShiftEventID = %d;`

const updateShiftForClockoutTemplate = `
UPDATE Shift SET ClockOutTime="%s", Earnings="%f", LastUpdated="%s" where ShiftEventID = %s;`

func UpdateShift(ctx context.Context, reqID string, req UpdateShiftRequest, db *sql.DB) (UpdateShiftResponse, error) {

	if req.ShiftEventID == 0 {
		return UpdateShiftResponse{OK: false, ID: 0, ReqID: reqID}, fmt.Errorf(statuscode.C500, "Missing ShiftEventID")
	}
	if req.ClockInTime.IsZero() {
		return UpdateShiftResponse{OK: false, ID: 0, ReqID: reqID}, fmt.Errorf(statuscode.C500, "Missing ClockInTime")
	}
	if req.ClockOutTime.IsZero() {
		return UpdateShiftResponse{OK: false, ID: 0, ReqID: reqID}, fmt.Errorf(statuscode.C500, "Missing ClockOutTime")
	}
	if req.ClockInTime.After(req.ClockOutTime) {
		return UpdateShiftResponse{OK: false, ID: 0, ReqID: reqID}, fmt.Errorf(statuscode.C500, "ClockInTime must be earlier than ClockOutTime")
	}

	var builtQuery = fmt.Sprintf(updateShiftTemplate, req.ClockInTime, req.ClockOutTime, req.Earnings, time.Now(), req.ShiftEventID)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return UpdateShiftResponse{OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf(statuscode.C500, "UpdateShift Err")
	}
	return UpdateShiftResponse{OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

// This updates the shift if it exists, it takes in the employeeID, ShiftEventId, and ClockIn
// returns true if successful, false if not
func UpdateShiftForClockout(ctx context.Context, db *sql.DB, employeeID string, ShiftEventID string, ClockIn time.Time) (bool, error) {
	//Get Shift clockInTime
	//Get Valid Wage
	wagePerHour, err := wage.GetWageForCurrentInterval(ctx, db, employeeID, ClockIn)

	if err != nil {
		return false, err
	}
	//Calculate Earnings
	now := time.Now()
	earnings := GetEarnings(ClockIn, now, wagePerHour)
	var builtQuery = fmt.Sprintf(updateShiftForClockoutTemplate, now, earnings, now, ShiftEventID)
	_, errr := db.ExecContext(ctx, builtQuery)

	if errr != nil {
		return false, errr
	}

	return true, nil
}
