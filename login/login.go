package login

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Admin struct {
	AdminID      string `json:"AdminID"`
	Email        string `json:"Email"`
	Company_Name string `json:"Company_Name"`
	AdminPIN     string `json:"AdminPIN"`
}

type Employee struct {
	EmployeeID string `json:"employeeID"`
	AdminID    string `json:"AdminID"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Birthday   string `json:"birthday"`
	JobTitle   string `json:"jobTitle"`
}

type EmployeeLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmployeeLoginResponse struct {
	RES   []Employee `json:"res"`
	DESC  string     `json:"desc"`
	OK    bool       `json:"ok"`
	ID    int64      `json:"id"`
	ReqID string     `json:"req_id"`
}

type AdminLoginResponse struct {
	RES   []Admin `json:"res"`
	DESC  string  `json:"desc"`
	OK    bool    `json:"ok"`
	ID    int64   `json:"id"`
	ReqID string  `json:"req_id"`
}

const employeeLogin = `SELECT EmployeeID, AdminID, Email, FirstName, LastName, Birthday, JobTitle FROM Employees WHERE Email = "%s" AND Password = "%s";`

const adminLogin = `SELECT AdminID, Email, Company_Name, AdminPIN FROM Admins WHERE Email = "%s" AND Password = "%s";`

func AdminLogin(ctx context.Context, reqID string, req AdminLoginRequest, db *sql.DB) (AdminLoginResponse, error) {
	if req.Email == "" {
		return AdminLoginResponse{DESC: "Could not get admin - missing Email", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Missing Email")
	}
	if req.Password == "" {
		return AdminLoginResponse{DESC: "Could not get admin - missing Password", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Missing Password")
	}

	var query = fmt.Sprintf(adminLogin, req.Email, req.Password)
	res, err := getQueryResAdmin(query, db)
	if res == nil {
		return AdminLoginResponse{DESC: "Could not get admin - incorrect email/password", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Incorrect Email and/or Password")
	}
	if err != nil {
		return AdminLoginResponse{DESC: "Could not get admin", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}
	return AdminLoginResponse{RES: res, OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil

}

func EmployeeLogin(ctx context.Context, reqID string, req EmployeeLoginRequest, db *sql.DB) (EmployeeLoginResponse, error) {
	if req.Email == "" {
		return EmployeeLoginResponse{DESC: "Could not get employee - missing Email", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Missing Email")
	}
	if req.Password == "" {
		return EmployeeLoginResponse{DESC: "Could not get employee - missing Password", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Missing Password")
	}

	var query = fmt.Sprintf(employeeLogin, req.Email, req.Password)
	res, err := getQueryResEmployee(query, db)
	if res == nil {
		return EmployeeLoginResponse{DESC: "Could not get employee - incorrect email/password", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, fmt.Errorf("Incorrect Email and/or Password")
	}
	if err != nil {
		return EmployeeLoginResponse{DESC: "Could not get admin", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, err
	}
	return EmployeeLoginResponse{RES: res, OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil

}

func getQueryResAdmin(builtQuery string, db *sql.DB) ([]Admin, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adminInfo []Admin

	for rows.Next() {
		var admin Admin
		if err := rows.Scan(&admin.AdminID, &admin.Email, &admin.Company_Name, &admin.AdminPIN); err != nil {
			return adminInfo, err
		}
		adminInfo = append(adminInfo, admin)
	}
	return adminInfo, nil
}

func getQueryResEmployee(builtQuery string, db *sql.DB) ([]Employee, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employeeInfo []Employee

	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.EmployeeID, &employee.AdminID, &employee.Email, &employee.FirstName, &employee.LastName, &employee.Birthday, &employee.JobTitle); err != nil {
			return employeeInfo, err
		}
		employeeInfo = append(employeeInfo, employee)
	}
	return employeeInfo, nil
}
