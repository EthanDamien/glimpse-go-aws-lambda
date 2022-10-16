package user

import (
	"context"
	"database/sql"
	"fmt"
)

type AdminResponse struct { // not sure what this is !!!!!! but im using it since getadmin has it
	DESC string `json:"body"`
}
type GetUsersRequest struct { //eventhing this requst needs is declared here
	AdminID int `json:"adminID"`
}

type User struct {
	EmloyeeID string `json:"EmployeeID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	JobTitle  string `json:"JobTitle"`
}

const getAllAdminUsers = `
SELECT JSON_ARRAYAGG(JSON_OBJECT('EmployeeID', EmployeeID, 'FirstName', FirstName, 'LastName', 
LastName, 'JobTitle', JobTitle)) from Employees WHERE AdminID = "%d";
`

func GetAdmin(ctx context.Context, reqID string, req GetUsersRequest, db *sql.DB) (AdminResponse, error) {
	//validate JSON
	var query = ""
	if req.AdminID != 0 { // not empty int
		query = getAllAdminUsers
	}
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return AdminResponse{DESC: "Error Querying ALL Employees"}, fmt.Errorf("Could not query ALL Employees")
	}

	var res string
	rows.Next()
	if err := rows.Scan(&res); err != nil {
		return AdminResponse{DESC: "Error Converting to String"}, fmt.Errorf("SQL conversion to String error")
	}

	return AdminResponse{DESC: res}, nil
}
