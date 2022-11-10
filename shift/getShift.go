package shift

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type GetMostRecentShiftsRequest struct {
	AdminID int `json:"adminID"`
}

type GetAllShiftsRequest struct {
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
	AdminID  int       `json:"adminID"`
}

type GetEmployeeShiftsRequest struct {
	EmployeeID int `json:"employeeID"`
}

type GetShiftRequest struct {
	ShiftEventID int `json:"shiftEventID"`
}

type GetShiftResponse struct {
	RES   []Shift `json:"res"`
	DESC  string  `json:"desc"`
	OK    bool    `json:"ok"`
	ID    int64   `json:"id"`
	ReqID string  `json:"req_id"`
}

type Shift struct {
	ShiftEventID int       `json:"shiftEventID"`
	EmployeeID   int       `json:"employeeID"`
	ClockInTime  time.Time `json:"clockInTime"`
	ClockOutTime time.Time `json:"clockOutTime"`
	Earnings     float32   `json:"earnings"`
	LastUpdated  time.Time `json:"lastUpdated"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
}

const getMostRecentShifts = `
SELECT ShiftEventID, s.EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated, FirstName, LastName FROM Shift s INNER JOIN Employees e ON e.EmployeeID = s.EmployeeID WHERE e.AdminID = %d ORDER BY LastUpdated DESC LIMIT 20;`

const getAllShiftsTemplate = `
SELECT ShiftEventID, s.EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated, FirstName, LastName FROM Shift s INNER JOIN Employees e ON e.EmployeeID = s.EmployeeID WHERE ClockInTime >= "%s" AND ClockOutTime <= "%s" AND e.AdminID = %d
ORDER BY LastUpdated;`

const getEmployeeShiftsTemplate = `
SELECT ShiftEventID, s.EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated, FirstName, LastName FROM Shift s INNER JOIN Employees e ON e.EmployeeID = s.EmployeeID WHERE EmployeeID = %d ORDER BY LastUpdated DESC;`

const getShiftTemplate = `
SELECT ShiftEventID, s.EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated, FirstName, LastName FROM Shift s INNER JOIN Employees e ON e.EmployeeID = s.EmployeeID WHERE ShiftEventID = %d ORDER BY LastUpdated DESC;`

func GetMostRecentShifts(ctx context.Context, reqID string, req GetMostRecentShiftsRequest, db *sql.DB) (GetShiftResponse, error) {

	var builtQuery = fmt.Sprintf(getMostRecentShifts, req.AdminID)
	res, err := getQueryRes(builtQuery, db)
	if err != nil {
		return GetShiftResponse{DESC: "Could not get shifts", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}

	return GetShiftResponse{RES: res, OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

func GetAllShifts(ctx context.Context, reqID string, req GetAllShiftsRequest, db *sql.DB) (GetShiftResponse, error) {

	if req.FromDate.IsZero() {
		return GetShiftResponse{DESC: "Could not get shifts - missing FromDate", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Missing FromDate")
	}
	if req.ToDate.IsZero() {
		return GetShiftResponse{DESC: "Could not get shifts - missing ToDate", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Missing ToDate")
	}

	var builtQuery = fmt.Sprintf(getAllShiftsTemplate, req.FromDate, req.ToDate, req.AdminID)
	res, err := getQueryRes(builtQuery, db)
	if err != nil {
		return GetShiftResponse{DESC: "Could not get shifts", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}

	return GetShiftResponse{RES: res, OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

func GetEmployeeShifts(ctx context.Context, reqID string, req GetEmployeeShiftsRequest, db *sql.DB) (GetShiftResponse, error) {

	if req.EmployeeID == 0 {
		return GetShiftResponse{DESC: "Could not get shifts - missing EmployeeID", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Missing EmployeeID")
	}

	var builtQuery = fmt.Sprintf(getEmployeeShiftsTemplate, req.EmployeeID)
	res, err := getQueryRes(builtQuery, db)

	if err != nil {
		return GetShiftResponse{DESC: "Could not get shifts", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}
	return GetShiftResponse{RES: res, OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

func GetShift(ctx context.Context, reqID string, req GetShiftRequest, db *sql.DB) (GetShiftResponse, error) {

	if req.ShiftEventID == 0 {
		return GetShiftResponse{DESC: "Could not get shifts - missing ShiftEventID", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Missing ShiftEventID")
	}

	var builtQuery = fmt.Sprintf(getShiftTemplate, req.ShiftEventID)
	res, err := getQueryRes(builtQuery, db)

	if err != nil {
		return GetShiftResponse{DESC: "Could not get shifts", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}
	return GetShiftResponse{RES: res, OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

func getQueryRes(builtQuery string, db *sql.DB) ([]Shift, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []Shift

	for rows.Next() {
		var shift Shift
		if err := rows.Scan(&shift.ShiftEventID, &shift.EmployeeID, &shift.ClockInTime,
			&shift.ClockOutTime, &shift.Earnings, &shift.LastUpdated, &shift.FirstName, &shift.LastName); err != nil {
			return shifts, err
		}
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func GetEarnings(pastTime time.Time, forwardTime time.Time, wage float64) float64 {
	diffInHours := forwardTime.Sub(pastTime).Hours()
	return wage * diffInHours
}
