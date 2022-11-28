package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

type DeleteUserRequest struct {
	EmployeeIDs []int `json:"employeeIDs"`
}

type DeleteUserResponse struct {
	DESC string `json:"desc"`
	OK   bool   `json:"ok"`
}

const deleteUserTemplate = `DELETE FROM Employees WHERE EmployeeID IN (%s);`

func DeleteUser(ctx context.Context, reqID string, req DeleteUserRequest, db *sql.DB) (DeleteUserResponse, error) {
	ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(req.EmployeeIDs)), ","), "[]")

	var builtQuery = fmt.Sprintf(deleteUserTemplate, ids)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return DeleteUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "Delete employee failed.")
	}
	return DeleteUserResponse{DESC: "Delete user success", OK: true}, nil
}
