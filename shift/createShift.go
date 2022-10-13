package shift

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CreateShiftRequest struct {
	EmployeeID   int       `json:"employeeID"`
	ClockInTime  time.Time `json:"clockInTime"`
	ClockOutTime time.Time `json:"clockOutTime"`
	Earnings     float32   `json:"earnings"`
}

type CreateShiftResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const createShiftTemplate = `
Insert into Shift (ShiftEventID, EmployeeID, ClockInTime, ClockOutTime, Earnings) 
values (NULL, %d, "%s", "%s", %f);`

func CreateShift(ctx context.Context, reqID string, req CreateShiftRequest, db *sql.DB) (CreateShiftResponse, error) {

	if req.EmployeeID == 0 {
		return CreateShiftResponse{DESC: "CreateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing EmployeeID")
	}
	if req.ClockInTime.IsZero() {
		return CreateShiftResponse{DESC: "CreateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing ClockInTime")
	}
	if req.ClockOutTime.IsZero() {
		return CreateShiftResponse{DESC: "CreateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing ClockOutTime")
	}
	if req.ClockInTime.After(req.ClockOutTime) {
		return CreateShiftResponse{DESC: "CreateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("ClockInTime must be earlier than ClockOutTime")
	}

	var builtQuery = fmt.Sprintf(createShiftTemplate, req.EmployeeID, req.ClockInTime, req.ClockOutTime, req.Earnings)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return CreateShiftResponse{DESC: "Could not insert into Shift Table", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}
	return CreateShiftResponse{DESC: "CreateShift success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
