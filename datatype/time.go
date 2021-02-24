package datatype

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"strings"
	"time"
)

type Time sql.NullTime

func (t *Time) Scan(value interface{}) (err error) {
	nullTime := sql.NullTime{}
	err = nullTime.Scan(value)
	*t = Time(nullTime)
	return
}

func (t Time) Value() (driver.Value, error) {
	y, m, d := t.Time.Date()
	h,M,s := t.Time.Clock()
	return time.Date(y, m, d, h, M, s, 0, t.Time.Location()), nil
}

// GormDataType gorm common data type
func (t Time) GormDataType() string {
	return "timestamp"
}

func (t Time) GobEncode() ([]byte, error) {
	return t.Time.GobEncode()
}

func (t *Time) GobDecode(b []byte) error {
	return t.Time.GobDecode(b)
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	tx := time.Unix(i, 0)
	nullTime := sql.NullTime{}
	err = nullTime.Scan(tx)
	*t = Time(nullTime)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	tx := t.Time
	return []byte(strconv.FormatInt(tx.Unix(), 10)), nil
}