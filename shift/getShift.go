package shift

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type GetAllShiftsRequest struct {
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

type GetEmployeeShiftsRequest struct {
	EmployeeID int `json:"employeeID"`
}

type GetShiftRequest struct {
	ShiftEventID int `json:"shiftEventID"`
}

type GetShiftResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

type Shift struct {
	ShiftEventID int       `json:"shiftEventID"`
	EmployeeID   int       `json:"employeeID"`
	ClockInTime  time.Time `json:"clockInTime"`
	ClockOutTime time.Time `json:"clockOutTime"`
	Earnings     float32   `json:"earnings"`
}

const getAllShiftsTemplate = `
SELECT * FROM Shift WHERE ClockInTime >= "%s" AND ClockOutTime <= "%s";`

const getEmployeeShiftsTemplate = `
SELECT * FROM Shift WHERE EmployeeID = %d;`

const getShiftTemplate = `
SELECT * FROM Shift WHERE ShiftEventID = %d;`

func GetAllShifts(ctx context.Context, reqID string, req GetAllShiftsRequest, db *sql.DB) (string, error) {
	//validate JSON
	if req.FromDate.IsZero() {
		return "", fmt.Errorf("Missing FromDate")
	}
	if req.ToDate.IsZero() {
		return "", fmt.Errorf("Missing ToDate")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(getAllShiftsTemplate, req.FromDate, req.ToDate)
	res, err := getJSON(builtQuery, db)
	return res, err
}

func GetEmployeeShifts(ctx context.Context, reqID string, req GetEmployeeShiftsRequest, db *sql.DB) (sql.Result, error) {

	//validate JSON
	if req.EmployeeID == 0 {
		return nil, fmt.Errorf("Missing EmployeeID")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(getEmployeeShiftsTemplate, req.EmployeeID)
	response, err := db.ExecContext(ctx, builtQuery)

	//If this fails, send "error" response
	//TODO send actual error to Lambda
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetShift(ctx context.Context, reqID string, req GetShiftRequest, db *sql.DB) (sql.Result, error) {

	//validate JSON
	if req.ShiftEventID == 0 {
		return nil, fmt.Errorf("Missing ShiftEventID")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(getShiftTemplate, req.ShiftEventID)
	response, err := db.ExecContext(ctx, builtQuery)

	//If this fails, send "error" response
	//TODO send actual error to Lambda
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getJSON(builtQuery string, db *sql.DB) (string, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}
