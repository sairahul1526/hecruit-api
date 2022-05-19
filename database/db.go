package database

import (
	"database/sql"
	CONFIG "hecruit-backend/config"
	CONSTANT "hecruit-backend/constant"
	"time"

	LOGGER "hecruit-backend/logger"
	UTIL "hecruit-backend/util"
	"strconv"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // for postgres driver
)

var db *sql.DB
var err error

// ConnectDatabase - connect to mysql database with given configuration
func ConnectDatabase() error {
	db, err = sql.Open("postgres", CONFIG.DBConfig)
	if err != nil {
		LOGGER.Log("ConnectDatabase", CONFIG.DBConfig, err)
		return err
	}

	// database conection pooling
	db.SetMaxOpenConns(CONFIG.DBConnectionPool)
	db.SetMaxIdleConns(CONFIG.DBConnectionPool)
	db.SetConnMaxLifetime(time.Hour)

	return nil
}

// database utils

// InsertWithUniqueID - insert data into table with unique id
func InsertWithUniqueID(table string, body map[string]string, key string) (string, sql.Result, error) {
	var (
		result sql.Result
		err    error
	)
	for i := 0; i < CONSTANT.NumberOfTimesUniqueInserts; i++ { // try to insert with unqiue id for certain number of times; if no limit, server crashes in certain conditions
		body[key] = generateRandomID()
		result, err = InsertSQL(table, body)
		if err == nil {
			break
		}
	}
	if err != nil {
		LOGGER.Log("InsertWithUniqueID", table, body, key, err)
		return "", nil, err
	}
	return body[key], result, nil
}

func generateRandomID() string {
	id := uuid.New()
	return id.String()
}

// RowCount - get number of items in database with specified query
func RowCount(tableName string, where string, args ...interface{}) (int, error) {
	data, err := SelectProcess("select count(*) as ctn from "+tableName+" where "+where, args...)
	if err != nil {
		LOGGER.Log("RowCount", tableName, where, args, err)
		return 0, err
	}
	if len(data) == 0 {
		return 0, nil
	}
	count, _ := strconv.Atoi(data[0]["ctn"])
	return count, nil
}

// CheckIfExists - check if data exists in table
func CheckIfExists(table string, params map[string]string) error {
	data, err := SelectSQL(table, []string{"1"}, params)
	if err != nil {
		LOGGER.Log("CheckIfExists", table, params, err)
		return err
	}
	if len(data) > 0 {
		return nil
	}
	return CONSTANT.SQLCheckIfExistsEmptyError
}

// sql wrapper functions

// ExecuteSQL - execute statement with defined values, with all the input as params to prevent sql injection
func ExecuteSQL(SQLQuery string, params ...interface{}) (sql.Result, error) {
	LOGGER.Log("ExecuteSQL", SQLQuery, params)
	result, err := db.Exec(SQLQuery, params...)
	if err != nil {
		LOGGER.Log("ExecuteSQL", SQLQuery, params, err)
		return nil, err
	}

	return result, nil
}

// QueryRowSQL - get single data with defined values, with all the input as params to prevent sql injection
func QueryRowSQL(SQLQuery string, params ...interface{}) (string, error) {
	var value string

	LOGGER.Log("QueryRowSQL", SQLQuery, params)
	err := db.QueryRow(SQLQuery, params...).Scan(&value)
	if err != nil {
		LOGGER.Log("QueryRowSQL", SQLQuery, params, err)
		return "", err
	}

	return value, nil
}

// UpdateSQL - update data with defined values
func UpdateSQL(tableName string, params map[string]string, body map[string]string) (sql.Result, error) {
	if len(body) == 0 {
		LOGGER.Log("UpdateSQL", tableName, params, body, CONSTANT.SQLUpdateBodyEmptyError)
		return nil, CONSTANT.SQLUpdateBodyEmptyError
	}

	args := []interface{}{}
	SQLQuery := "update " + tableName + " set "

	init := false
	i := 1
	for key, val := range body {
		if init {
			SQLQuery += ","
		}
		SQLQuery += `"` + key + `" = $` + strconv.Itoa(i)
		args = append(args, val)
		init = true
		i++
	}
	// add updated_at
	SQLQuery += `, "updated_at" = $` + strconv.Itoa(i)
	args = append(args, UTIL.GetCurrentTime())
	i++

	SQLQuery += " where "
	init = false
	for key, val := range params {
		if init {
			SQLQuery += " and "
		}
		SQLQuery += `"` + key + `" = $` + strconv.Itoa(i)
		args = append(args, val)
		init = true
		i++
	}

	LOGGER.Log("UpdateSQL", SQLQuery, args)
	result, err := db.Exec(SQLQuery, args...)
	if err != nil {
		LOGGER.Log("UpdateSQL", tableName, params, body, err)
		return nil, err
	}

	return result, nil
}

// DeleteSQL - delete data with defined values
func DeleteSQL(tableName string, params map[string]string) (sql.Result, error) {
	if len(params) == 0 {
		// atleast one value should be specified for deleting, cannot delete all values
		LOGGER.Log("DeleteSQL", tableName, params, CONSTANT.SQLDeleteAllNotAllowedError)
		return nil, CONSTANT.SQLDeleteAllNotAllowedError
	}

	args := []interface{}{}
	SQLQuery := "delete from " + tableName + " where "

	init := false
	i := 1
	for key, val := range params {
		if init {
			SQLQuery += " and "
		}
		SQLQuery += `"` + key + `" = $` + strconv.Itoa(i)
		args = append(args, val)
		init = true
		i++
	}

	LOGGER.Log("DeleteSQL", SQLQuery, args)
	result, err := db.Exec(SQLQuery, args...)
	if err != nil {
		LOGGER.Log("DeleteSQL", tableName, params, err)
		return nil, err
	}
	return result, nil
}

// InsertSQL - insert data with defined values
func InsertSQL(tableName string, body map[string]string) (sql.Result, error) {
	if len(body) == 0 {
		LOGGER.Log("InsertSQL", tableName, body, CONSTANT.SQLInsertBodyEmptyError)
		return nil, CONSTANT.SQLInsertBodyEmptyError
	}

	SQLQuery, args := BuildInsertStatement(tableName, body)

	LOGGER.Log("InsertSQL", SQLQuery, args)
	result, err := db.Exec(SQLQuery, args...)
	if err != nil {
		LOGGER.Log("InsertSQL", tableName, body, err)
		return nil, err
	}
	return result, nil
}

// BuildInsertStatement - build insert statement with defined values
func BuildInsertStatement(tableName string, body map[string]string) (string, []interface{}) {
	args := []interface{}{}
	SQLQuery := "insert into " + tableName + " "
	keys := " ("
	values := " ("
	init := false
	i := 1
	for key, val := range body {
		if init {
			keys += ","
			values += ","
		}
		keys += ` "` + key + `" `
		values += " $" + strconv.Itoa(i)
		args = append(args, val)
		init = true
		i++
	}
	// add created_at, updated_at
	keys += `, "created_at", "updated_at" `
	values += ", $" + strconv.Itoa(i) + ", $" + strconv.Itoa(i+1)
	args = append(args, UTIL.GetCurrentTime(), UTIL.GetCurrentTime())

	keys += ")"
	values += ")"
	SQLQuery += keys + " values " + values
	return SQLQuery, args
}

// SelectSQL - query data with defined values
func SelectSQL(tableName string, columns []string, params map[string]string) ([]map[string]string, error) {
	args := []interface{}{}
	SQLQuery := "select " + strings.Join(columns, ",") + " from " + tableName + ""
	if len(params) > 0 {
		where := ""
		init := false
		i := 1
		for key, val := range params {
			if init {
				where += " and "
			}
			where += ` "` + key + `" = $` + strconv.Itoa(i)
			args = append(args, val)
			init = true
			i++
		}
		if strings.Compare(where, "") != 0 {
			SQLQuery += " where " + where
		}
	}
	return SelectProcess(SQLQuery, args...)
}

// SelectProcess - execute raw select statement, with all the input as params to prevent sql injection
func SelectProcess(SQLQuery string, params ...interface{}) ([]map[string]string, error) {
	LOGGER.Log("SelectProcess", SQLQuery, params)
	rows, err := db.Query(SQLQuery, params...)
	if err != nil {
		LOGGER.Log("SelectProcess", SQLQuery, params, err)
		return []map[string]string{}, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		LOGGER.Log("SelectProcess", SQLQuery, params, err)
		return []map[string]string{}, err
	}

	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols))
	data := []map[string]string{}
	rest := map[string]string{}
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		rest = map[string]string{}
		err = rows.Scan(dest...)
		if err != nil {
			LOGGER.Log("SelectProcess", SQLQuery, params, err)
			return []map[string]string{}, err
		}

		for i, raw := range rawResult {
			if raw == nil {
				rest[cols[i]] = ""
			} else {
				rest[cols[i]] = string(raw)
			}
		}

		data = append(data, rest)
	}

	return data, nil
}
