package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const createUserTemplate = `
Insert into Employees (EmployeeID, AdminID, Email, Password, FirstName, LastName, Birthday, JobTitle) 
values (NULL, "%d", "%s", "%s", "%s", "%s", "%s", "%s");`

func CreateUser(ctx context.Context, reqID string, req CreateUserRequest, db *sql.DB) (CreateUserResponse, error) {
	if req.AdminID <= 0 {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("adminID is missing")
	}
	if req.Email == "" {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("email is missing")
	}
	if req.Password == "" {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("password is missing")
	}
	if req.FirstName == "" {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("first name is missing")
	}
	if req.LastName == "" {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("last name is missing")
	}
	if req.Birthday == "" {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("birthday is missing")
	}
	if len(req.Birthday) != 10 {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("birthday is in wrong format. Must be YYYY-MM-DD")
	}
	if req.JobTitle == "" {
		return CreateUserResponse{DESC: "CreateUser err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("job title is missing")
	}

	var builtQuery = fmt.Sprintf(createUserTemplate, req.AdminID, req.Email, req.Password, req.FirstName, req.LastName, req.Birthday, req.JobTitle)
	_, err := db.ExecContext(ctx, builtQuery)

	if err != nil {
		return CreateUserResponse{DESC: "Could not insert into Employees Table", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, nil
	}
	return CreateUserResponse{DESC: "CreateUser success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
