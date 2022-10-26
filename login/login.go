package login

import (
	"context"
	"database/sql"
	"fmt"
)

type EmployeeLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

const adminLogin = `
SELECT JSON_ARRAYAGG(JSON_OBJECT('AdminID', AdminID, 'Email', Email, 'Password', 
Password, 'Company_Name', Company_Name, 'AdminPIN', AdminPIN)) from Admins WHERE Email = "%s" AND Password = "%s";
`
const employeeLogin = `
SELECT JSON_ARRAYAGG(JSON_OBJECT('EmployeeID', EmployeeID, 'AdminID', AdminID, 'Email', Email, 'Password', 
Password, 'FirstName', FirstName, 'LastName', LastName, 'Birthday', Birthday, 'JobTitle', JobTitle)) from Employees WHERE Email = "%s" AND Password = "%s";
`

func AdminLogin(ctx context.Context, reqID string, req AdminLoginRequest, db *sql.DB) (LoginResponse, error) {
	var query = fmt.Sprintf(adminLogin, req.Email, req.Password)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return LoginResponse{DESC: "Error Querying Admins Table", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("could not login as admin")
	}
	if rows == nil {
		return LoginResponse{DESC: "Error Querying Admins Table", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Invalid email or password/admin does not exist")
	}
	return LoginResponse{DESC: "Successful Admin Login"}, nil
}

func EmployeeLogin(ctx context.Context, reqID string, req EmployeeLoginRequest, db *sql.DB) (LoginResponse, error) {
	var query = fmt.Sprintf(employeeLogin, req.Email, req.Password)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return LoginResponse{DESC: "Error Querying Employees Table", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Could not login as employee")
	}
	if rows == nil {
		return LoginResponse{DESC: "Error Querying Employees Table", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Invalid email or password/employee does not exist")
	}
	return LoginResponse{DESC: "Successful Employee Login"}, nil
}
