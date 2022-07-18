package db

import (
	"database/sql"
	"fmt"
	"github.com/jamesineda/reschedular/app/utils"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
	"strings"
	"unicode"
)

const (
	select_    = "SELECT"
	insertInto = "INSERT INTO"
	from       = "FROM"
	values     = "VALUES"
	where_     = "WHERE"
	and        = "AND"
	in         = "IN"
	null       = "NULL"
)

type SQLXClient interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type Client interface {
	GetById(id string, table interface{}) (interface{}, error)
	GetList(rows interface{}, filters Filters) error
	Create(object interface{}) error
}

type DatabaseConn struct {
	SQLXClient
}

func NewFakeDatabaseConn(fake *FakeSQLX) (Client, error) {
	return &DatabaseConn{fake}, nil
}

func NewDatabaseConn(config *utils.DatabaseConfig) (Client, error) {
	switch config.Driver {
	case "fake":
		return &DatabaseConn{&FakeSQLX{}}, nil

	default:
		db, err := sqlx.Connect(config.Driver, config.Dsn)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}

		return &DatabaseConn{db}, nil
	}
}

/*
	A more polished solution could be something like a Filters class with a pre-defined set of supported
	comparators and values, which would be validated and handled in the database client class, as to avoid SQL errors.
	For the sake of time, I'm just implementing a basic slice of arguments that will be expected like:

	[
		["attr_a", "=", "foo"],
		["attr_b", "<", 100],
		["attr_c", "IN", "(1,2,3)"],
	]
*/
type Filters [][]interface{}

func (f *Filters) Values() (values []interface{}) {
	for _, filter := range *f {
		values = append(values, filter[2])
	}
	return
}

func (db *DatabaseConn) GetById(id string, table interface{}) (interface{}, error) {
	tableName, selectFields := getSelectOptions(table)

	// I'd prefer an incremental ID on the database table, as well as created_at/ updated_at timestamps, which
	// would be used for ordering in all queries.
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id=$1 LIMIT 1", selectFields, tableName)
	if err := db.Get(table, query, id); err != nil {
		return nil, err
	}

	return table, nil
}

func (db *DatabaseConn) GetList(table interface{}, filters Filters) error {
	tableName, selectFields := getSelectOptions(table)
	query := generateSelectQuery(tableName, selectFields, filters)
	if err := db.Select(table, query, filters.Values()...); err != nil {
		return err
	}

	return nil
}

func (db *DatabaseConn) Create(object interface{}) error {
	tableName, selectFields := getSelectOptions(object)
	tags := getTags(object, "db")
	fieldNames := ":" + strings.Join(tags, ",:")
	query := generateInsertQuery(tableName, selectFields, fieldNames)
	_, err := db.NamedExec(query, object)
	return err
}

func generateInsertQuery(tableName, fieldNames, namedExecColName string) string {
	return strings.Join([]string{insertInto, tableName, "(", fieldNames, ")", values, "(", namedExecColName, ")"}, " ")
}

func generateSelectQuery(tableName, fieldNames string, filters Filters) string {
	query := strings.Join([]string{select_, fieldNames, from, tableName}, " ")
	if filters != nil && len(filters) > 0 {
		for fi, filter := range filters {

			whereOrAnd := and
			if fi == 0 {
				whereOrAnd = where_
			}

			var q string
			if filter[2] != nil || filter[2] == null {
				q = "?"
			} else {
				q = null
			}

			if filter[1] == in {
				q = "("
				v := reflect.ValueOf(filter[2])
				for i := 0; i < v.Len(); i++ {
					switch v.Index(i).Type().Kind() {
					case reflect.String:
						q += "'" + fmt.Sprint(v.Index(i).Interface()) + "',"
					default:
						q += fmt.Sprint(v.Index(i).Interface()) + ","
					}
				}
				q = q[:len(q)-1] + ")"
			}

			query = strings.Join([]string{query, whereOrAnd, fmt.Sprintf("%v", filter[0]), fmt.Sprintf("%v", filter[1]), q}, " ")
		}
	}
	return query
}

func getSelectOptions(table interface{}) (tn, fields string) {
	tn = getTableName(table)
	tags := getTags(table, "db")
	fields = strings.Join(tags, ",")
	return
}

func getTableName(obj interface{}) string {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	names := strings.Split(t.Name(), ".")
	name := names[len(names)-1]
	return toSnakecase(name) + "s"
}

func toSnakecase(in string) string {
	runes := []rune(in)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) ||
			unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

func getTags(s interface{}, tk string) []string {
	tags := make([]string, 0)
	v := reflect.TypeOf(s)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get(tk)
		tags = append(tags, strings.Split(tag, ",")[0])
	}
	return tags
}
