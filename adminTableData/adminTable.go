package adminTableData

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EthanDamien/glimpse-go-aws-lambda/statuscode"
)

type GetAdminTableDataReq struct {
	AdminID   int    `json:"adminID"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type GetAdminTableDataRes struct {
	StatusCode string           `json:"StatusCode"`
	RES        []AdminTableData `json:"AdminTableData"`
}

func GetAdminTableData(ctx context.Context, reqID string, req GetAdminTableDataReq, db *sql.DB) (GetAdminTableDataRes, error) {
	if req.AdminID == 0 {
		return GetAdminTableDataRes{
			StatusCode: statuscode.C500,
		}, fmt.Errorf("AdminID Missing")
	}

	if req.StartDate == "" {
		return GetAdminTableDataRes{
			StatusCode: statuscode.C500,
		}, fmt.Errorf("StartDate Missing")
	}

	if req.EndDate == "" {
		return GetAdminTableDataRes{
			StatusCode: statuscode.C500,
		}, fmt.Errorf("EndDate Missing")
	}

	var builtQuery = fmt.Sprintf(GetDataForInterval, req.AdminID, req.StartDate, req.EndDate)

	res, err := getQueryRes(builtQuery, db)

	if err != nil {
		return GetAdminTableDataRes{
			StatusCode: statuscode.C500,
		}, fmt.Errorf("Query Err")
	}

	return GetAdminTableDataRes{
		StatusCode: statuscode.C200,
		RES:        res,
	}, nil
}

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
