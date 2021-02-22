package datatype

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"strings"
	"time"
)

type Time time.Time

func (t *Time) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*t = Time(nullTime.Time)
	return
}

func (t Time) Value() (driver.Value, error) {
	y, m, d := time.Time(t).Date()
	h,M,s := time.Time(t).Clock()
	return time.Date(y, m, d, h, M, s, 0, time.Time(t).Location()), nil
}

// GormDataType gorm common data type
func (t Time) GormDataType() string {
	return "timestamp"
}

func (t Time) GobEncode() ([]byte, error) {
	return time.Time(t).GobEncode()
}

func (t *Time) GobDecode(b []byte) error {
	return (*time.Time)(t).GobDecode(b)
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	tx := time.Unix(i, 0)
	*t = Time(tx)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	tx := time.Time(t)
	return []byte(strconv.FormatInt(tx.Unix(), 10)), nil
}