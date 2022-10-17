package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type GetAllUsersRequest struct {
	AdminID int `json:"adminID"`
}

type GetUserResponse struct {
	RES   []User `json:"res"`
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

type User struct {
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	JobTitle   string `json:"jobTitle"`
	EmployeeID string `json:"employeeID"`
}

const getAllUsers = `SELECT e.EmployeeID, e.FirstName, e.LastName, e.JobTitle, e.Email FROM Employees e WHERE e.AdminID = %d;`

func GetAllUsers(ctx context.Context, reqID string, req GetAllUsersRequest, db *sql.DB) (GetUserResponse, error) {

	var builtQuery = fmt.Sprintf(getAllUsers, req.AdminID)
	res, err := getQueryRes(builtQuery, db)
	if err != nil {
		return GetUserResponse{DESC: "Could not get shifts", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}

	return GetUserResponse{RES: res, OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

func getQueryRes(builtQuery string, db *sql.DB) ([]User, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Email, &user.EmployeeID, &user.FirstName,
			&user.LastName, &user.JobTitle); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
