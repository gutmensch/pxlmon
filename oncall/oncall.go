/*
  Combined Monitoring Tool for mobile.de / ebayK / motortalk

  (c) Robert Schumann <roschumann@ebay.com>, mobile.de GmbH

  License: BSD
*/
package oncall


import (
    //"fmt"
    "time"
    //"strconv"
    "os"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)

var location = "Europe/Berlin"


func GetDutyOfficer() string {

    officer := queryInfraDb("select login as result from " + GetShiftType() + " join persons using (person_id) where duty_date=CURDATE()")

    return officer
}

func inTimeSpan(start, end, check time.Time) bool {
    return check.After(start) && check.Before(end)
}

func isNonBusinessDay() bool {
    return (len(queryInfraDb("select nonbusiness_day as result from nonbusiness_days where nonbusiness_day = CURDATE()")) != 0)
}

func GetShiftType() string {

    const defaultShift string = "duty"
    const dayShift string = "irq"

    loc, _ := time.LoadLocation(location)
    t := time.Now().In(loc)

    dayShiftStart := time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, loc)
    dayShiftEnd := time.Date(t.Year(), t.Month(), t.Day(), 18, 0, 0, 0, loc)
    day := time.Now().In(loc).Weekday()

    if day != time.Saturday &&
       day != time.Sunday &&
       !isNonBusinessDay() &&
       inTimeSpan(dayShiftStart, dayShiftEnd, time.Now()) {
         return dayShift
    }

    return defaultShift
}

func queryInfraDb(cmd string) string {

    defaultResult := ""
    dbConnection := os.Getenv("ONCALLDB")

    db, err := sql.Open("mysql", dbConnection)
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }

    var result sql.NullString
    row := db.QueryRow(cmd)
    if err := row.Scan(&result); err == nil {
        return result.String;
    } else if err == sql.ErrNoRows {
      // noop
    } else {
        panic(err.Error())
    }

    return defaultResult
}
