package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type DeleteUserRequest struct {
	EmployeeIDs []int `json:"employeeIDs"`
}

type DeleteUserResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const deleteUserTemplate = `DELETE FROM Employees WHERE EmployeeID IN (%s);`

func DeleteUser(ctx context.Context, reqID string, req DeleteUserRequest, db *sql.DB) (DeleteUserResponse, error) {
	ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(req.EmployeeIDs)), ","), "[]")

	var builtQuery = fmt.Sprintf(deleteUserTemplate, ids)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return DeleteUserResponse{DESC: "Could not delete employee from the Employees Table", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, nil
	}
	return DeleteUserResponse{DESC: "Delete user success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
