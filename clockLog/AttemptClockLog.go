package clockLog

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/EthanDamien/glimpse-go-aws-lambda/compare"
	"github.com/EthanDamien/glimpse-go-aws-lambda/shift"
	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

type AttemptClockLogRes struct {
	StatusCode string `json:"StatusCode"`
	Event      string `json:"Event"`
	EmployeeID int    `json:"EmployeeID"`
}

type AttemptClockLogReq struct {
	AdminID       string `json:"AdminID"`
	PictureMeta64 string `json:"PictureMeta64"`
}

func AttemptClockLog(ctx context.Context, reqID string, req AttemptClockLogReq, db *sql.DB) (AttemptClockLogRes, error) {
	//validate JSON
	if req.AdminID == "" {
		return AttemptClockLogRes{
			Event: "MissingAdminID",
		}, fmt.Errorf(statuscode.C500, "Missing AdminID")
	}
	if req.PictureMeta64 == "" {
		return AttemptClockLogRes{
			Event: "MissingImage",
		}, fmt.Errorf(statuscode.C500, "Missing Image")
	}
	employeeID, err := compare.FindMatchingEmployee(req.AdminID, req.PictureMeta64, db)
	if err != nil {
		if err.Error() == "Employee Not Found" {
			//Return
			return AttemptClockLogRes{
				Event: "Employee Not Found",
			}, err
		}
		return AttemptClockLogRes{
			Event: "Find Matching Employee Error",
		}, err
	}

	//clock in/out
	clockLog, err := shift.GenerateShiftorUpdate(ctx, strconv.Itoa(employeeID), db)
	if err != nil {
		return AttemptClockLogRes{
			Event: "Error when Generating/updating shift",
		}, err
	}

	return AttemptClockLogRes{
		Event:      clockLog,
		EmployeeID: employeeID,
	}, nil

}
