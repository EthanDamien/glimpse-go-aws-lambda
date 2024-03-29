package wage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// This query will get the first wage for the current interval (where it's time to set is <= clockIn Time)
const getWageForCurrentIntervalTemplate = `
select WagePerHour from Wage 
where EmployeeId = %s and TimeToSet <= CAST("%s" as DATETIME) 
order by TimeToSet desc limit 1;`

// struct for getting wage per hour
type WagePerHour struct {
	WagePerHour float64 `json:"WagePerHour"`
}

// gets the wage (float) for a given interval, error otherwise
func GetWageForCurrentInterval(ctx context.Context, db *sql.DB, EmployeeID string, ClockIn time.Time) (float64, error) {
	//Get Shift clockInTime
	//Get Valid Wage
	//Calculate Earnings
	var builtQuery = fmt.Sprintf(getWageForCurrentIntervalTemplate, EmployeeID, ClockIn)
	res, err := getQueryRes(builtQuery, db)

	if err != nil {
		return 0, err
	}

	if len(res) != 1 {
		return 0, fmt.Errorf(statuscode.C500, "No Valid Wage Detected")
	}

	return float64(res[0].WagePerHour), nil
}

// gets database query results for wage given employee and time interval
// returns an array of WagePerHour instances, else error
func getQueryRes(builtQuery string, db *sql.DB) ([]WagePerHour, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wagesPerHour []WagePerHour

	for rows.Next() {
		var checkshift WagePerHour
		if err := rows.Scan(&checkshift.WagePerHour); err != nil {
			return wagesPerHour, err
		}
		wagesPerHour = append(wagesPerHour, checkshift)
	}
	return wagesPerHour, nil
}
