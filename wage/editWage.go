package wage

type EditWageRequest struct {
	WageEventID string `json:WageEventID`
	EmployeeID  string `json:"EmployeeID"`
	WagePerHour string `json:"WagePerHour"`
	TimeToSet   string `json:"TimeToSet"`
}

type EditWageResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const editWageTemplate = `Insert into Wage (WageEventID, EmployeeID, WagePerHour, TimeToSet) values (NULL, %s, %s, "%s"); `

//TODO: Edit wage

// func EditWage(ctx context.Context, reqID string, req EditWageRequest, db *sql.DB) (EditWageResponse, error) {
// 	if req.EmployeeID == "" {
// 		return EditWageResponse{DESC: "EditWage err"}, fmt.Errorf("Missing EmployeeID")
// 	}
// 	if req.WagePerHour == "" {
// 		return EditWageResponse{DESC: "EditWage err"}, fmt.Errorf("Missing WagePerHour")
// 	}
// 	if req.TimeToSet == "" {
// 		return EditWageResponse{DESC: "EditWage err"}, fmt.Errorf("Missing TimeToSet")
// 	}

// 	var builtQuery = fmt.Sprintf(createWageTemplate, req.EmployeeID, req.WagePerHour, req.TimeToSet)
// 	_, err := db.ExecContext(ctx, builtQuery)

// 	if err != nil {
// 		return CreateWageResponse{DESC: "CreateWage err"}, fmt.Errorf("Missing Password")
// 	}

// 	return CreateWageResponse{DESC: fmt.Sprintf("Wage Created with values EmployeeID: %s, Wage %s, TimeToSet %s",
// 		req.EmployeeID, req.WagePerHour, req.TimeToSet)}, nil
// }
