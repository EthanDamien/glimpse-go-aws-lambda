package shift

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

type checkShift struct {
	ShiftEventID int       `json:"shiftEventID"`
	ClockInTime  time.Time `json:"clockInTime"`
}

const createShiftTemplate = `
Insert into Shift (ShiftEventID, EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated) 
values (NULL, %d, "%s", "%s", %f, "%s");`

const createShiftTemplateForClockIn = `
Insert into Shift (ShiftEventID, EmployeeID, ClockInTime, ClockOutTime, Earnings, LastUpdated) 
values (NULL, %d, "%s", NULL, NULL, "%s");`

const checkActiveShiftTemplate = `select ShiftEventID, ClockInTime from Shift where EmployeeID = %s and ClockOutTime is NULL;`

func CreateShift(ctx context.Context, reqID string, req CreateShiftRequest, db *sql.DB) (CreateShiftResponse, error) {

	if req.EmployeeID == 0 {
		return CreateShiftResponse{DESC: "CreateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing EmployeeID")
	}
	if req.ClockInTime.IsZero() {
		return CreateShiftResponse{DESC: "CreateShift err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing ClockInTime")
	}

	var builtQuery = fmt.Sprintf(createShiftTemplate, req.EmployeeID, req.ClockInTime, req.ClockOutTime, req.Earnings, time.Now())
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return CreateShiftResponse{DESC: "Could not insert into Shift Table", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}
	return CreateShiftResponse{DESC: "CreateShift success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

// This function will update the generate a shift if it doesn't exist, otherwise, it will modify the
// active shift
// return err if there are internal errors
// return string if shift clockin / clockout
func GenerateShiftorUpdate(ctx context.Context, employeeID string, db *sql.DB) (string, error) {
	//check if there's an active shift
	activeShiftExists, shiftID, clockIn, err := checkForActiveShift(ctx, employeeID, db)

	if err != nil {
		return "", err
	}

	//Shift does not exist
	if !activeShiftExists {
		//generate shift
		err := GenerateShiftForClockIn(ctx, employeeID, db)
		if err != nil {
			return "", err
		}
		return "clockin", nil
	}

	//shift exists
	clockedout, err := UpdateShiftForClockout(ctx, db, employeeID, strconv.Itoa(shiftID), clockIn)

	if err != nil {
		return "", err
	}

	if clockedout == true {
		return "clockedout", err
	}

	return "", nil
}

func GenerateShiftForClockIn(ctx context.Context, employeeID string, db *sql.DB) error {
	employeeIDAsInt, err := strconv.Atoi(employeeID)
	if err != nil {
		return err
	}
	var builtQuery = fmt.Sprintf(createShiftTemplateForClockIn, employeeIDAsInt, time.Now(), time.Now())
	_, errr := db.ExecContext(ctx, builtQuery)

	if errr != nil {
		return errr
	}
	return nil
}

// This function checks if there is an active shift for this employee
// If this is found, it will return that shiftID for editing
// If not found, return 0
// If not, it will return 0 with a nil err
// If input > 1, this will return an error, as this case should be handled elsewhere
func checkForActiveShift(ctx context.Context, employeeID string, db *sql.DB) (bool, int, time.Time, error) {
	var builtQuery = fmt.Sprintf(checkActiveShiftTemplate, employeeID)
	res, err := getCheckShiftQueryRes(builtQuery, db)

	if err != nil {
		return false, 0, time.Now(), err
	}

	if len(res) > 1 {
		return false, 0, time.Now(), fmt.Errorf("There's an issue here, this should never happen")
	}

	if len(res) == 1 {
		return true, res[0].ShiftEventID, res[0].ClockInTime, nil
	}

	return false, 0, time.Now(), nil

}

func getCheckShiftQueryRes(builtQuery string, db *sql.DB) ([]checkShift, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkShifts []checkShift

	for rows.Next() {
		var checkshift checkShift
		if err := rows.Scan(&checkshift.ShiftEventID, &checkshift.ClockInTime); err != nil {
			return checkShifts, err
		}
		checkShifts = append(checkShifts, checkshift)
	}
	return checkShifts, nil
}
