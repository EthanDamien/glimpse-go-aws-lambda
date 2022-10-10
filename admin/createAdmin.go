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

const createAdmin = `
Insert into Admins (AdminID, Email, Password, Company_Name, AdminPin) 
values (NULL, "a", "a", "a", "a");`

func CreateAdmin(ctx context.Context, reqID string, req CreateAdminRequest, db *sql.DB) (CreateAdminResponse, error) {
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

	db.Ping()
	_, err := db.ExecContext(ctx, createAdmin)
	if err != nil {
		return CreateAdminResponse{DESC: "Could not insert into Admin Table", OK: false, ID: time.Now().UnixNano(), ReqID: reqID}, nil
	}
	return CreateAdminResponse{DESC: "CreateAdmin success", OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}
