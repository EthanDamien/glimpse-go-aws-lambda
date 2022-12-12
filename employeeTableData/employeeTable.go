package employeeTableData

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// This struct is the request for the Employee Table Data
type GetEmployeeTableDataReq struct {
	EmployeeID int `json:"employeeID"`
}

// this is the struct for the Employee Table Data Result
type GetEmployeeTableDataRes struct {
	MinutesForWeek   int     `json:"MinutesForWeek"`
	MinutesForMonth  int     `json:"MinutesForMonth"`
	MinutesForYear   int     `json:"MinutesForYear"`
	EarningsForWeek  float64 `json:"EarningsForWeek"`
	EarningsForMonth float64 `json:"EarningsForMonth"`
	EarningsForYear  float64 `json:"EarningsForYear"`
}

// This function gets the Employee Table Data
// return GetEmployeeTableDataRes instance if successful, else error
func GetEmployeeTableData(ctx context.Context, reqID string, req GetEmployeeTableDataReq, db *sql.DB) (GetEmployeeTableDataRes, error) {
	if req.EmployeeID == 0 {
		return GetEmployeeTableDataRes{}, fmt.Errorf(statuscode.C500, "EmployeeID Missing")
	}

	var weekQuery = fmt.Sprintf(GetDataFromWeekTemplate, req.EmployeeID)
	var monthQuery = fmt.Sprintf(GetDataForMonthTemplate, req.EmployeeID)
	var yearQuery = fmt.Sprintf(GetDataForYearTemplate, req.EmployeeID)

	resWeek, err := getQueryRes(weekQuery, db)

	var weekMinutes = 0
	var weekEarnings = 0.0
	if len(resWeek) != 0 {
		weekMinutes = resWeek[0].Minutes
		weekEarnings = resWeek[0].Earnings
	}

	if err != nil {
		return GetEmployeeTableDataRes{}, fmt.Errorf(statuscode.C500, "Week Query Err")
	}

	resMonth, err := getQueryRes(monthQuery, db)
	var monthMinutes = 0
	var monthEarnings = 0.0
	if len(resMonth) != 0 {
		monthMinutes = resMonth[0].Minutes
		monthEarnings = resMonth[0].Earnings
	}

	if err != nil {
		return GetEmployeeTableDataRes{}, fmt.Errorf(statuscode.C500, "Month Query Err")
	}

	resYear, err := getQueryRes(yearQuery, db)
	var yearMinutes = 0
	var yearEarnings = 0.0
	if len(resYear) != 0 {
		yearMinutes = resYear[0].Minutes
		yearEarnings = resYear[0].Earnings
	}

	if err != nil {
		return GetEmployeeTableDataRes{}, fmt.Errorf(statuscode.C500, "Year Query Err")
	}

	return GetEmployeeTableDataRes{
		MinutesForWeek:   weekMinutes,
		MinutesForMonth:  monthMinutes,
		MinutesForYear:   yearMinutes,
		EarningsForWeek:  weekEarnings,
		EarningsForMonth: monthEarnings,
		EarningsForYear:  yearEarnings,
	}, nil
}

// get query result for EmployeeTable Data
// return array of employeeTableData objects if successful, else error
func getQueryRes(builtQuery string, db *sql.DB) ([]employeeTableData, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employeeDataArr []employeeTableData

	for rows.Next() {
		var employeeData employeeTableData
		if err := rows.Scan(&employeeData.EmployeeID,
			&employeeData.Minutes,
			&employeeData.Earnings); err != nil {
			return employeeDataArr, err
		}
		employeeDataArr = append(employeeDataArr, employeeData)
	}
	return employeeDataArr, nil
}
