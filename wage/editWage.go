package wage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

type EditWageRequest struct {
	WageEventID int       `json:"WageEventID"`
	WagePerHour float32   `json:"WagePerHour"`
	TimeToSet   time.Time `json:"TimeToSet"`
}

type EditWageResponse struct {
	DESC string `json:"desc"`
	OK   bool   `json:"ok"`
}

const editWageTemplate = `UPDATE Wage SET WagePerHour = %f, TimeToSet = "%s" WHERE WageEventID = %d;`

// returns a status code 200 response if the wage for a given employee
// (based on wageEventID) was successfully updated with a description
// returns a status code 500 response if there was an error
func EditWage(ctx context.Context, reqID string, req EditWageRequest, db *sql.DB) (EditWageResponse, error) {
	var builtQuery = fmt.Sprintf(editWageTemplate, req.WagePerHour, req.TimeToSet, req.WageEventID)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return EditWageResponse{OK: false}, fmt.Errorf(statuscode.C500, "Couldn't edit wage")
	}

	return EditWageResponse{DESC: fmt.Sprintf("Wage successfully updated"), OK: true}, nil
}
