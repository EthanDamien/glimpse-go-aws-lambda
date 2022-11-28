package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

type CreateUserRequest struct {
	AdminID   int    `json:"adminID"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Birthday  string `json:"birthday"`
	JobTitle  string `json:"jobTitle"`
}

type CreateUserResponse struct {
	DESC string `json:"desc"`
	OK   bool   `json:"ok"`
}

const createUserTemplate = `
Insert into Employees (EmployeeID, AdminID, Email, Password, FirstName, LastName, Birthday, JobTitle) 
values (NULL, "%d", "%s", "%s", "%s", "%s", "%s", "%s");`

func CreateUser(ctx context.Context, reqID string, req CreateUserRequest, db *sql.DB) (CreateUserResponse, error) {
	if req.AdminID <= 0 {
		return CreateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "adminID is missing")
	}
	if req.Email == "" {
		return CreateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "email is missing")
	}
	if req.Password == "" {
		return CreateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "password is missing")
	}
	if req.FirstName == "" {
		return CreateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "first name is missing")
	}
	if req.LastName == "" {
		return CreateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "last name is missing")
	}
	if req.Birthday == "" {
		return CreateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "birthday is missing")
	}
	if len(req.Birthday) != 10 {
		return CreateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "birthday is in wrong format. Must be YYYY-MM-DD")
	}
	if req.JobTitle == "" {
		return CreateUserResponse{OK: false}, fmt.Errorf(statuscode.C500, "job title is missing")
	}

	var builtQuery = fmt.Sprintf(createUserTemplate, req.AdminID, req.Email, req.Password, req.FirstName, req.LastName, req.Birthday, req.JobTitle)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return CreateUserResponse{OK: false}, nil
	}
	return CreateUserResponse{DESC: "CreateUser success", OK: true}, nil
}
