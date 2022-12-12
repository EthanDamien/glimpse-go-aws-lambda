package adminTableData

const GetDataForInterval = `
select CAST(ClockOutTime as DATE) as "Date", SUM(Earnings) as "Earnings" 
from (
select ClockOutTime, Earnings 
from Shift 
join Employees on Shift.EmployeeID = Employees.EmployeeID where
AdminID = %d
and
ClockOutTime >= CAST("%s" as DATE) and
ClockOutTime <= CAST("%s" as DATE)
)as IntervalSpecified
group by day(ClockOutTime);`

// structure for the admin table data
type AdminTableData struct {
	Date     string  `json:"Date"`
	Earnings float64 `json:"Earnings"`
}
