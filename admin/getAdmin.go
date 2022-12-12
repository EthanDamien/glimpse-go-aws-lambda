package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// request format for getting an admin by their email
type GetAdminRequest struct {
	Email string `json:"email"`
}

// request format for getting an admin by their ID
type GetAdminByAdminIDRequest struct {
	AdminID int `json:"adminID"`
}

// format for an Admin object
type Admin struct {
	AdminID      string `json:"AdminID"`
	Email        string `json:"Email"`
	Password     string `json:"Password"`
	Company_Name string `json:"Company_Name"`
	AdminPIN     string `json:"AdminPIN"`
}

// response format for getting an admin by their ID
type AdminIDResponse struct {
	RES   Admin  `json:"res"`
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
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
SELECT * from Admins WHERE AdminID = "%d";
`

// get an admin by their email
// return AdminResponse instance if successful, else error
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
		return AdminResponse{}, fmt.Errorf(statuscode.C500, "Could not query Admins")
	}

	var res string
	rows.Next()
	if err := rows.Scan(&res); err != nil {
		return AdminResponse{}, fmt.Errorf(statuscode.C500, "SQL conversion to String error")
	}

	return AdminResponse{DESC: res}, nil
}

// get an admin by their ID
// return AdminIDResponse instance if successful, else error
func GetAdminByAdminID(ctx context.Context, reqID string, req GetAdminByAdminIDRequest, db *sql.DB) (AdminIDResponse, error) {
	//validate JSON
	if req.AdminID == 0 {
		return AdminIDResponse{}, fmt.Errorf(statuscode.C500, "No AdminID")
	}

	var query = ""
	query = fmt.Sprintf(getSpecificAdminByID, req.AdminID)
	res, err := getQueryRes(query, db)
	if err != nil {
		return AdminIDResponse{}, fmt.Errorf(statuscode.C500, "Could not query Admins")
	}

	return AdminIDResponse{RES: res, OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}

// Perform query
// return Admin instance if successful, else error
func getQueryRes(builtQuery string, db *sql.DB) (Admin, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return Admin{}, err
	}
	defer rows.Close()

	var adminFound Admin

	for rows.Next() {
		var admin Admin
		if err := rows.Scan(&admin.AdminID, &admin.Email, &admin.Password, &admin.Company_Name, &admin.AdminPIN); err != nil {
			return admin, err
		}
		adminFound = admin
	}
	return adminFound, nil
}
