package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

type UpdateUserPasswordRequest struct {
	EmployeeID  int    `json:"employeeID"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdateUserResponse struct {
	DESC string `json:"desc"`
	OK   bool   `json:"ok"`
}

const updateUserPwdTemplate = `
UPDATE Employees SET Password="%s" WHERE EmployeeID = "%d"`

const getUserPwdTemplate = `
SELECT Password FROM Employees WHERE EmployeeID = "%d"`

func UpdateUserPassword(ctx context.Context, reqID string, req UpdateUserPasswordRequest, db *sql.DB) (UpdateUserResponse, error) {

	if req.OldPassword == "" {
		return UpdateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing Old Password")
	}
	if req.NewPassword == "" {
		return UpdateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing New Password")
	}
	// get user's password
	var builtQuery = fmt.Sprintf(getUserPwdTemplate, req.EmployeeID)
	pwd, err := getPwdQueryRes(builtQuery, db)
	if err != nil {
		return UpdateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "Couldn't fetch user's original password")
	}
	if pwd != req.OldPassword {
		return UpdateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "Incorrect Original Password")
	}

	var updateQuery = fmt.Sprintf(updateUserPwdTemplate, req.NewPassword, req.EmployeeID)
	_, updateErr := db.ExecContext(ctx, updateQuery)

	if updateErr != nil {
		return UpdateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "Could not update user password")
	}
	return UpdateUserResponse{DESC: "Updated user password", OK: true}, nil
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
