package shift

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
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
	RES  []Shift `json:"RES"`
	DESC string  `json:"DESC"`
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
SELECT ShiftEventID, s.EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated, FirstName, LastName FROM Shift s INNER JOIN Employees e ON e.EmployeeID = s.EmployeeID WHERE s.EmployeeID = %d ORDER BY ClockInTime DESC;`

const getShiftTemplate = `
SELECT ShiftEventID, s.EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated, FirstName, LastName FROM Shift s INNER JOIN Employees e ON e.EmployeeID = s.EmployeeID WHERE ShiftEventID = %d ORDER BY LastUpdated DESC;`

const getShiftsBetweenDatesTemplate = `
Select ShiftEventID, s.EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated, FirstName, LastName 
from Shift where EmployeeID = "7" and 
ClockInTime >= CAST("%s" as DATETIME) and 
ClockInTime <= CAST("%s" as DATETIME) and
ClockOutTime >= CAST("%s" as DATETIME) and 
ClockOutTime <= CAST("%s" as DATETIME) 
ORDER BY ShiftEventID DESC;`

func GetMostRecentShifts(ctx context.Context, reqID string, req GetMostRecentShiftsRequest, db *sql.DB) (GetShiftResponse, error) {
	if req.AdminID == 0 {
		return GetShiftResponse{}, fmt.Errorf(statuscode.C500, "Missing AdminID")
	}
	var builtQuery = fmt.Sprintf(getMostRecentShifts, req.AdminID)
	res, err := getQueryRes(builtQuery, db)
	if err != nil {
		return GetShiftResponse{
			DESC: "Could not get shifts",
		}, fmt.Errorf(statuscode.C500, "RecentShifts Err")
	}

	return GetShiftResponse{
		RES: res,
	}, nil
}

func GetAllShifts(ctx context.Context, reqID string, req GetAllShiftsRequest, db *sql.DB) (GetShiftResponse, error) {

	if req.FromDate.IsZero() {
		return GetShiftResponse{}, fmt.Errorf(statuscode.C500, "Missing FromDate")
	}
	if req.ToDate.IsZero() {
		return GetShiftResponse{}, fmt.Errorf(statuscode.C500, "Missing ToDate")
	}

	var builtQuery = fmt.Sprintf(getAllShiftsTemplate, req.FromDate, req.ToDate, req.AdminID)
	res, err := getQueryRes(builtQuery, db)
	if err != nil {
		return GetShiftResponse{}, err
	}

	return GetShiftResponse{
		RES: res,
	}, nil
}

func GetEmployeeShifts(ctx context.Context, reqID string, req GetEmployeeShiftsRequest, db *sql.DB) (GetShiftResponse, error) {

	if req.EmployeeID == 0 {
		return GetShiftResponse{}, fmt.Errorf(statuscode.C500, "Missing EmployeeID")
	}

	var builtQuery = fmt.Sprintf(getEmployeeShiftsTemplate, req.EmployeeID)
	res, err := getQueryRes(builtQuery, db)

	if err != nil {
		return GetShiftResponse{}, fmt.Errorf(statuscode.C500, "GetEmployeeShifts err")
	}
	return GetShiftResponse{
		RES: res,
	}, nil
}

func GetShift(ctx context.Context, reqID string, req GetShiftRequest, db *sql.DB) (GetShiftResponse, error) {

	if req.ShiftEventID == 0 {
		return GetShiftResponse{}, fmt.Errorf(statuscode.C500, "Missing ShiftEventID")
	}

	var builtQuery = fmt.Sprintf(getShiftTemplate, req.ShiftEventID)
	res, err := getQueryRes(builtQuery, db)

	if err != nil {
		return GetShiftResponse{}, fmt.Errorf(statuscode.C500, "GetShift err")
	}
	return GetShiftResponse{
		RES: res,
	}, nil
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
