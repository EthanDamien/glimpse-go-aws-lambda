package employeeTableData

const GetDataFromWeekTemplate = `
select EmployeeID, SUM(Minutes) as Minutes, SUM(Earnings) as Earnings
from (
	select EmployeeID, TIMESTAMPDIFF(MINUTE, ClockInTime, ClockOutTime) as Minutes, Earnings from Shift where 
	EmployeeID = %d and
	week(ClockInTime)=week(now()) 
	order by clockInTime asc
	) as WeekData
group by EmployeeID;`

const GetDataForMonthTemplate = `
select EmployeeID, SUM(Minutes) as Minutes, SUM(Earnings) as Earnings
from (
	select EmployeeID, TIMESTAMPDIFF(MINUTE, ClockInTime, ClockOutTime) as Minutes, Earnings from Shift where 
	EmployeeID = %d and
	month(ClockInTime)=month(now()) 
	order by clockInTime asc
	) as MonthData
group by EmployeeID;`

const GetDataForYearTemplate = `
select EmployeeID, SUM(Minutes) as Minutes, SUM(Earnings) as Earnings
from (
	select EmployeeID, TIMESTAMPDIFF(MINUTE, ClockInTime, ClockOutTime) as Minutes, Earnings from Shift where 
	EmployeeID = %d and
	ClockInTime >= DATE_ADD(CURRENT_TIMESTAMP(), INTERVAL -365 DAY) 
	order by clockInTime asc
    ) as YearData
group by EmployeeID;`

type employeeTableData struct {
	EmployeeID int     `json:"EmployeeID"`
	Minutes    int     `json:"Minutes"`
	Earnings   float64 `json:"Earnings"`
}
