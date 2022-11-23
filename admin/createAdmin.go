package admin

import (
	"context"
	"database/sql"
	"fmt"
)

type CreateAdminRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Company_Name string `json:"companyName"`
	AdminPIN     string `json:"adminPin"`
}

// Template to be used to insert to Table
const createAdminTemplate = `
Insert into Admins (AdminID, Email, Password, Company_Name, AdminPin) 
values (NULL, "%s", "%s", "%s", "%s");`

func CreateAdmin(ctx context.Context, reqID string, req CreateAdminRequest, db *sql.DB) (AdminResponse, error) {

	//validate JSON
	if req.Email == "" {
		return AdminResponse{DESC: "CreateAdmin err"}, fmt.Errorf("Status:500 Missing Email")
	}
	if req.Password == "" {
		return AdminResponse{DESC: "CreateAdmin err"}, fmt.Errorf("Status:500 Missing Password")
	}
	if req.Company_Name == "" {
		return AdminResponse{DESC: "CreateAdmin err"}, fmt.Errorf("Status:500 Missing Company Name")
	}
	if req.AdminPIN == "" {
		return AdminResponse{DESC: "CreateAdmin err"}, fmt.Errorf("Status:500 Missing AdminPIN")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(createAdminTemplate, req.Email, req.Password, req.Company_Name, req.AdminPIN)
	_, err := db.ExecContext(ctx, builtQuery)

	//If this fails, send "error" response
	//TODO send actual error to Lambda
	if err != nil {
		return AdminResponse{DESC: "Could not insert into Admin Table"}, fmt.Errorf("Status:500 Internal server error")
	}
	return AdminResponse{DESC: "Inserted into table"}, nil
}
