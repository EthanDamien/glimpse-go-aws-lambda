package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CreateAdminRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Company_Name string `json:"companyName"`
	AdminPIN     string `json:"adminPin"`
}

type CreateAdminResponse struct {
	DESC  string `json:"desc"`
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

// Template to be used to insert to Table
const createAdminTemplate = `
Insert into Admins (AdminID, Email, Password, Company_Name, AdminPin) 
values (NULL, "%s", "%s", "%s", "%s");`

func CreateAdmin(ctx context.Context, reqID string, req CreateAdminRequest, db *sql.DB) (CreateAdminResponse, error) {

	//validate JSON
	if req.Email == "" {
		return CreateAdminResponse{DESC: "CreateAdmin err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("missing Email")
	}
	if req.Password == "" {
		return CreateAdminResponse{DESC: "CreateAdmin err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing Password")
	}
	if req.Company_Name == "" {
		return CreateAdminResponse{DESC: "CreateAdmin err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing Company Name")
	}
	if req.AdminPIN == "" {
		return CreateAdminResponse{DESC: "CreateAdmin err", OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("Missing AdminPIN")
	}

	//Use the template and fill in the blanks
	var builtQuery = fmt.Sprintf(createAdminTemplate, req.Email, req.Password, req.Company_Name, req.AdminPIN)
	_, err := db.ExecContext(ctx, builtQuery)

	//If this fails, send "error" response
	//TODO send actual error to Lambda
	if err != nil {
		return CreateAdminResponse{DESC: "Could not insert into Admin Table", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, nil
	}
	return CreateAdminResponse{DESC: "CreateAdmin success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
