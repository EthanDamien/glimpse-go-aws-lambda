package wage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// This query will get the first wage for the current interval (where it's time to set is <= clockIn Time)
const getWagesForEmployee = `
select WageEventID, EmployeeID, WagePerHour, TimeToSet from Wage 
where EmployeeId = %d 
order by TimeToSet desc;`

type WageInfo struct {
	WageEventID int       `json:"WageEventID"`
	EmployeeID  int       `json:"EmployeeID"`
	WagePerHour float32   `json:"WagePerHour"`
	TimeToSet   time.Time `json:"TimeToSet"`
}

type GetWageRequest struct {
	EmployeeID int `json:"EmployeeID"`
}

type GetWageResponse struct {
	RES  []WageInfo `json:"res"`
	DESC string     `json:"desc"`
	OK   bool       `json:"ok"`
}

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
