package db

import (
	"database/sql"
	"database/sql/driver"
)

type FakeSQLX struct {
	GetReturn    interface{}
	SelectReturn interface{}
}

func NewSetFakeSQLX(get interface{}, selReturn interface{}) *FakeSQLX {
	return &FakeSQLX{GetReturn: get, SelectReturn: selReturn}
}

func (f *FakeSQLX) Get(dest interface{}, query string, args ...interface{}) error {
	dest = f.GetReturn
	return nil
}

func (f *FakeSQLX) Select(dest interface{}, query string, args ...interface{}) error {
	dest = f.SelectReturn
	return nil
}

func (f *FakeSQLX) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return driver.RowsAffected(1), nil
}
