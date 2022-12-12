package wage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// This query will get all the unique wages an employee ever had at the company
const getWagesForEmployee = `
select WageEventID, EmployeeID, WagePerHour, TimeToSet from Wage 
where EmployeeId = %d 
order by TimeToSet desc;`

// struct for getting wage information
type WageInfo struct {
	WageEventID int       `json:"WageEventID"`
	EmployeeID  int       `json:"EmployeeID"`
	WagePerHour float32   `json:"WagePerHour"`
	TimeToSet   time.Time `json:"TimeToSet"`
}

// request format for getting a wage
type GetWageRequest struct {
	EmployeeID int `json:"EmployeeID"`
}

// response format for getting a wage
type GetWageResponse struct {
	RES  []WageInfo `json:"res"`
	DESC string     `json:"desc"`
	OK   bool       `json:"ok"`
}

// returns status code 200 response and a GetWageResponse instance with all wage information for an employee
// returns status code 500 response and an error if there's an error
func GetWagesForEmployees(ctx context.Context, reqID string, req GetWageRequest, db *sql.DB) (GetWageResponse, error) {
	// Get all wages for an employee
	if req.EmployeeID == 0 {
		return GetWageResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing EmployeeID")
	}
	var builtQuery = fmt.Sprintf(getWagesForEmployee, req.EmployeeID)
	res, err := getQueryResult(builtQuery, db)
	if err != nil {
		return GetWageResponse{OK: false}, fmt.Errorf(statuscode.C500, "Couldn't get wages for employee")
	}
	if len(res) == 0 {
		return GetWageResponse{OK: false}, fmt.Errorf(statuscode.C500, "Wages for employee doesn't exist")
	}
	return GetWageResponse{RES: res, OK: true}, nil

}

// gets database query results for all wage information for an employee
// return an array of WageInfo instances if successful, else error
func getQueryResult(builtQuery string, db *sql.DB) ([]WageInfo, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wageInfo []WageInfo

	for rows.Next() {
		var wageRow WageInfo
		if err := rows.Scan(&wageRow.WageEventID, &wageRow.EmployeeID, &wageRow.WagePerHour, &wageRow.TimeToSet); err != nil {
			return wageInfo, err
		}
		wageInfo = append(wageInfo, wageRow)
	}
	return wageInfo, nil
}
