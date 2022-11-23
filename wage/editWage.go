package wage

import (
	"context"
	"database/sql"
	"fmt"
)

type EditWageRequest struct {
	WageEventID int `json:"WageEventID"`
	WagePerHour int `json:"WagePerHour"`
}

type EditWageResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const editWageTemplate = `UPDATE Wage SET WagePerHour = %d WHERE WageEventID = %d;`

//TODO: Edit wage

func EditWage(ctx context.Context, reqID string, req EditWageRequest, db *sql.DB) (EditWageResponse, error) {
	var builtQuery = fmt.Sprintf(editWageTemplate, req.WagePerHour, req.WageEventID)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return EditWageResponse{DESC: "EditWage err"}, fmt.Errorf("Coudln't edit wage")
	}

	return EditWageResponse{DESC: fmt.Sprintf("Wage successfully updated")}, nil
}
