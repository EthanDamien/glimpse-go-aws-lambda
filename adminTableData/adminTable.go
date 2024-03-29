package adminTableData

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

// request format for getting admin table data
type GetAdminTableDataReq struct {
	AdminID   int    `json:"adminID"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// response format for getting admin table data
type GetAdminTableDataRes struct {
	StatusCode string           `json:"StatusCode"`
	RES        []AdminTableData `json:"AdminTableData"`
}

// get data and statistics for admin dashboard
// return GetAdminTableDataRes instance if successful, else error
func GetAdminTableData(ctx context.Context, reqID string, req GetAdminTableDataReq, db *sql.DB) (GetAdminTableDataRes, error) {
	if req.AdminID == 0 {
		return GetAdminTableDataRes{}, fmt.Errorf(statuscode.C500, "AdminID Missing")
	}

	if req.StartDate == "" {
		return GetAdminTableDataRes{}, fmt.Errorf(statuscode.C500, "StartDate Missing")
	}

	if req.EndDate == "" {
		return GetAdminTableDataRes{}, fmt.Errorf(statuscode.C500, "EndDate Missing")
	}

	var builtQuery = fmt.Sprintf(GetDataForInterval, req.AdminID, req.StartDate, req.EndDate)

	res, err := getQueryRes(builtQuery, db)

	if err != nil {
		return GetAdminTableDataRes{}, fmt.Errorf(statuscode.C500, "Query Err")
	}

	return GetAdminTableDataRes{
		RES: res,
	}, nil
}

// perform query and return result
// return an array of AdminTableData instances if successful, else error
func getQueryRes(builtQuery string, db *sql.DB) ([]AdminTableData, error) {
	rows, err := db.Query(builtQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var AdminDataArr []AdminTableData

	for rows.Next() {
		var adminData AdminTableData
		if err := rows.Scan(
			&adminData.Date,
			&adminData.Earnings); err != nil {
			return AdminDataArr, err
		}
		AdminDataArr = append(AdminDataArr, adminData)
	}
	return AdminDataArr, nil
}
