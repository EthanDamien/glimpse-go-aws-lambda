package clockLog

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/EthanDamien/glimpse-go-aws-lambda/compare"
	"github.com/EthanDamien/glimpse-go-aws-lambda/shift"
)

type AttemptClockLogRes struct {
	DESC string `json:"body"`
}

type AttemptClockLogReq struct {
	AdminID       string `json:"AdminID"`
	PictureMeta64 string `json:"PictureMeta64"`
}

func AttemptClockLog(ctx context.Context, reqID string, req AttemptClockLogReq, db *sql.DB) (AttemptClockLogRes, error) {
	//validate JSON
	if req.AdminID == "" {
		return AttemptClockLogRes{DESC: "AttemptClockLog err"}, fmt.Errorf("Missing AdminID")
	}
	if req.PictureMeta64 == "" {
		return AttemptClockLogRes{DESC: "AttemptClockLog err"}, fmt.Errorf("Missing Image")
	}
	employeeID, err := compare.FindMatchingEmployee(req.AdminID, req.PictureMeta64, db)
	if err != nil {
		if err.Error() == "Employee Not Found" {
			//Return
			return AttemptClockLogRes{DESC: "Employee Not Found err"}, err
		}
		return AttemptClockLogRes{DESC: "AttemptClockLogErr"}, err
	}

	//clock in/out
	clockLog, err := shift.GenerateShiftorUpdate(ctx, strconv.Itoa(employeeID), db)
	if err != nil {
		return AttemptClockLogRes{DESC: "Error when Generating/updating shift"}, err
	}

	return AttemptClockLogRes{DESC: fmt.Sprintf("%s, %d", clockLog, employeeID)}, nil

}
