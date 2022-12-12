package login

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// struct for an Admin
type Admin struct {
	AdminID      string `json:"AdminID"`
	Email        string `json:"Email"`
	Company_Name string `json:"Company_Name"`
	AdminPIN     string `json:"AdminPIN"`
}

// struct for an Employee
type Employee struct {
	EmployeeID string `json:"employeeID"`
	AdminID    string `json:"AdminID"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Birthday   string `json:"birthday"`
	JobTitle   string `json:"jobTitle"`
}

// request format for logging in an employee
type EmployeeLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// request format for logging in an admin
type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// response format for logging in an employee
type EmployeeLoginResponse struct {
	RES  []Employee `json:"res"`
	DESC string     `json:"desc"`
	OK   bool       `json:"ok"`
}

// response format for logging in an admin
type AdminLoginResponse struct {
	RES  []Admin `json:"res"`
	DESC string  `json:"desc"`
	OK   bool    `json:"ok"`
}

const employeeLogin = `SELECT EmployeeID, AdminID, Email, FirstName, LastName, Birthday, JobTitle FROM Employees WHERE Email = "%s" AND Password = "%s";`

const adminLogin = `SELECT AdminID, Email, Company_Name, AdminPIN FROM Admins WHERE Email = "%s" AND Password = "%s";`

// this function returns a status code 200 response and an AdminLoginResponse instance if admin successfully logs in
// returns status code 500 for all errors and an error with an appropriate error message
func AdminLogin(ctx context.Context, reqID string, req AdminLoginRequest, db *sql.DB) (AdminLoginResponse, error) {
	if req.Email == "" {
		return AdminLoginResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing Email")
	}
	if req.Password == "" {
		return AdminLoginResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing Password")
	}

	var query = fmt.Sprintf(adminLogin, req.Email, req.Password)
	res, err := getQueryResAdmin(query, db)
	if res == nil {
		return AdminLoginResponse{OK: false}, fmt.Errorf(statuscode.C500, "Incorrect Email and/or Password")
	}
	if err != nil {
		return AdminLoginResponse{OK: false}, fmt.Errorf(statuscode.C500, "Could not log in admin.")
	}
	return AdminLoginResponse{RES: res, OK: true}, nil

}

// returns a status code 200 response and an EmployeeLoginResponse instance if employee successfully logs in
// returns status code 500 and an error for all errors with an appropriate error message
func EmployeeLogin(ctx context.Context, reqID string, req EmployeeLoginRequest, db *sql.DB) (EmployeeLoginResponse, error) {
	if req.Email == "" {
		return EmployeeLoginResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing Email")
	}
	if req.Password == "" {
		return EmployeeLoginResponse{OK: false}, fmt.Errorf(statuscode.C500, "Missing Password")
	}

	var query = fmt.Sprintf(employeeLogin, req.Email, req.Password)
	res, err := getQueryResEmployee(query, db)
	if res == nil {
		return EmployeeLoginResponse{OK: false}, fmt.Errorf(statuscode.C500, "Incorrect Email and/or Password")
	}
	if err != nil {
		return EmployeeLoginResponse{OK: false}, fmt.Errorf(statuscode.C500, "Could not log in employee")
	}
	return EmployeeLoginResponse{RES: res, OK: true}, nil

}

// gets database query results for admin login query
// return an array of Admin objects if successful, else error
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

// gets database results for employee login query
// return an array of Employee instances if successful, else error
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
