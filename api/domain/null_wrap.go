package domain

import (
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"bytes"
	"time"
)

var nullLiteral = []byte("null")
type NullInt64 struct {
	sql.NullInt64
}
func NewNullInt(ni int64) NullInt64 {
	return NullInt64{NullInt64:sql.NullInt64{ni, ni != 0}}
}
func (i NullInt64) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int64)
	} else {
		return json.Marshal(nil)
	}
}
func (i *NullInt64) UnmarshalJSON(data []byte) error {
	var num int64
	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}
	i.Int64 = num
	i.Valid = num != 0
	return nil
}


type NullTime struct {
	mysql.NullTime
}
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return nt.Time.MarshalJSON()
	} else {
		return nullLiteral, nil
	}
}
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if bytes.Compare(data, []byte("null")) == 0 {
		nt.Valid = false
		return nil
	}
	now := time.Now()
	err := nt.UnmarshalJSON(data)

	if err != nil {
		return err
	}
	nt.Valid = true
	nt.Time = now

	return nil
}

