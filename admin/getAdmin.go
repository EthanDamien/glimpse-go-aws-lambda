package admin

import (
	"context"
	"database/sql"
	"fmt"
)

type GetAdminRequest struct {
	Email string `json:"email"`
}

type GetAdminByAdminIDRequest struct {
	AdminID int `json:"adminID"`
}

type Admin struct {
	AdminID      string `json:"AdminID"`
	AdminPIN     string `json:"AdminPIN"`
	Company_Name string `json:"Company_Name"`
	Email        string `json:"Email"`
	Password     string `json:"Password"`
}

const getAllAdmins = `
SELECT JSON_ARRAYAGG(JSON_OBJECT('AdminID', AdminID, 'Email', Email, 'Password', 
Password, 'Company_Name', Company_Name, 'AdminPIN', AdminPIN)) from Admins;
`

const getSpecificAdmin = `
SELECT JSON_ARRAYAGG(JSON_OBJECT('AdminID', AdminID, 'Email', Email, 'Password', 
Password, 'Company_Name', Company_Name, 'AdminPIN', AdminPIN)) from Admins WHERE Email = "%s";
`
const getSpecificAdminByID = `
SELECT JSON_ARRAYAGG(JSON_OBJECT('AdminID', AdminID, 'Email', Email, 'Password', 
Password, 'Company_Name', Company_Name, 'AdminPIN', AdminPIN)) from Admins WHERE AdminID = "%d";
`

func GetAdmin(ctx context.Context, reqID string, req GetAdminRequest, db *sql.DB) (AdminResponse, error) {
	//validate JSON
	var query = ""
	if req.Email == "" {
		query = getAllAdmins
	} else {
		query = fmt.Sprintf(getSpecificAdmin, req.Email)
	}
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return AdminResponse{DESC: "Error Querying Admins Table"}, fmt.Errorf("Could not query Admins")
	}

	var res string
	rows.Next()
	if err := rows.Scan(&res); err != nil {
		return AdminResponse{DESC: "Error Converting to String"}, fmt.Errorf("SQL conversion to String error")
	}

	return AdminResponse{DESC: res}, nil
}

func GetAdminByAdminID(ctx context.Context, reqID string, req GetAdminByAdminIDRequest, db *sql.DB) (AdminResponse, error) {
	//validate JSON
	var query = ""
	query = fmt.Sprintf(getSpecificAdminByID, req.AdminID)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return AdminResponse{DESC: "Error Querying Admins Table"}, fmt.Errorf("Could not query Admins")
	}

	var res string
	rows.Next()
	if err := rows.Scan(&res); err != nil {
		return AdminResponse{DESC: "Error Converting to String"}, fmt.Errorf("SQL conversion to String error")
	}

	return AdminResponse{DESC: res}, nil
}
