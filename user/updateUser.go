package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UpdateUserPasswordRequest struct {
	EmployeeID  int    `json:"employeeID"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdateUserResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const updateUserPwdTemplate = `
UPDATE Employees SET Password="%s" WHERE EmployeeID = "%d"`

const getUserPwdTemplate = `
SELECT Password FROM Employees WHERE EmployeeID = "%d"`

func UpdateUserPassword(ctx context.Context, reqID string, req UpdateUserPasswordRequest, db *sql.DB) (UpdateUserResponse, error) {

	if req.OldPassword == "" {
		return UpdateUserResponse{DESC: "UpdateUserPassword err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing Old Password")
	}
	if req.NewPassword == "" {
		return UpdateUserResponse{DESC: "UpdateUserPassword err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing New Password")
	}
	// get user's password
	var builtQuery = fmt.Sprintf(getUserPwdTemplate, req.EmployeeID)
	pwd, err := getPwdQueryRes(builtQuery, db)
	if err != nil {
		return UpdateUserResponse{DESC: "UpdateUserPassword err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Couldn't fetch user's original password")
	}
	if pwd != req.OldPassword {
		return UpdateUserResponse{DESC: "UpdateUserPassword err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Incorrect Original Password")
	}

	var updateQuery = fmt.Sprintf(updateUserPwdTemplate, req.NewPassword, req.EmployeeID)
	_, updateErr := db.ExecContext(ctx, updateQuery)

	if updateErr != nil {
		return UpdateUserResponse{DESC: "Could not update user password", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, updateErr
	}
	return UpdateUserResponse{DESC: "Updated user password", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

func getPwdQueryRes(builtQuery string, db *sql.DB) (string, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	var pwd string

	for rows.Next() {
		var password string
		if err := rows.Scan(&password); err != nil {
			return password, err
		}
		pwd = password
	}

	return pwd, nil
}
