package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// request format for getting all employees
type GetAllUsersRequest struct {
	AdminID int `json:"adminID"`
}

// response format for getting all employees
type GetUserResponse struct {
	RES  []User `json:"res"`
	DESC string `json:"desc"`
	OK   bool   `json:"ok"`
}

// format of a User object
type User struct {
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	JobTitle   string `json:"jobTitle"`
	EmployeeID string `json:"employeeID"`
}

const getAllUsers = `SELECT e.Email, e.FirstName, e.LastName, e.JobTitle, e.EmployeeID FROM Employees e WHERE e.AdminID = %d;`

// Get all users in the employee table
// returns GetUserResponse if successful, else error
func GetAllUsers(ctx context.Context, reqID string, req GetAllUsersRequest, db *sql.DB) (GetUserResponse, error) {

	var builtQuery = fmt.Sprintf(getAllUsers, req.AdminID)
	res, err := getQueryRes(builtQuery, db)
	if err != nil {
		return GetUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "Could not retrieve all users.")
	}

	return GetUserResponse{RES: res, OK: true}, nil
}

// Perform the query and return results
// return an array of User instances, else error
func getQueryRes(builtQuery string, db *sql.DB) ([]User, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Email, &user.FirstName, &user.LastName,
			&user.JobTitle, &user.EmployeeID); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
