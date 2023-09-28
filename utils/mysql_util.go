package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/canteen_management/logger"
)

const (
	logTagMysqlUtil = "mysql_util"
)

const (
	maxDealPerTimes = 512
)

type SqlHandle interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func GetSqlPlaceholder(len int) string {
	if len <= 0 {
		return ""
	}
	return strings.Repeat(",?", len)[1:]
}

func SqlQueryRowWithLock(mysql SqlHandle, tableName string, resultDAO interface{}, condition string, condValues ...interface{}) error {
	allField := GetFieldsTagByKey(resultDAO, "json")
	allAddr := GetFieldsAddr(resultDAO)

	sqlStr := fmt.Sprintf("SELECT `%v` FROM `%v` %v FOR UPDATE", strings.Join(allField, "`,`"), tableName, condition)

	row := mysql.QueryRow(sqlStr, condValues...)
	return row.Scan(allAddr...)
}

func SqlQueryRow(mysql SqlHandle, tableName string, resultDAO interface{}, condition string, condValues ...interface{}) error {
	allField := GetFieldsTagByKey(resultDAO, "json")
	allAddr := GetFieldsAddr(resultDAO)

	sqlStr := fmt.Sprintf("SELECT `%v` FROM `%v` %v", strings.Join(allField, "`,`"), tableName, condition)

	row := mysql.QueryRow(sqlStr, condValues...)
	return row.Scan(allAddr...)
}

func SqlQueryWithLock(mysql SqlHandle, tableName string, dao interface{}, condition string, condValues ...interface{}) (resultSlice interface{}, err error) {
	allField := GetFieldsTagByKey(dao, "json")

	sqlStr := fmt.Sprintf("SELECT `%v` FROM `%v` %v FOR UPDATE", strings.Join(allField, "`,`"), tableName, condition)

	rows, err := mysql.Query(sqlStr, condValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	daoType := reflect.TypeOf(dao)
	result := reflect.MakeSlice(reflect.SliceOf(daoType), 0, 0)
	for rows.Next() {
		t := reflect.New(daoType.Elem())
		err = rows.Scan(GetFieldsAddr(t.Interface())...)
		if err != nil {
			return result, err
		}
		result = reflect.Append(result, t)
	}
	return result.Interface(), nil
}

func SqlQuery(mysql SqlHandle, tableName string, dao interface{}, condition string, condValues ...interface{}) (resultSlice interface{}, err error) {
	allField := GetFieldsTagByKey(dao, "json")

	sqlStr := fmt.Sprintf("SELECT `%v` FROM `%v` %v", strings.Join(allField, "`,`"), tableName, condition)

	rows, err := mysql.Query(sqlStr, condValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	daoType := reflect.TypeOf(dao)
	result := reflect.MakeSlice(reflect.SliceOf(daoType), 0, 0)
	for rows.Next() {
		t := reflect.New(daoType.Elem())
		err := rows.Scan(GetFieldsAddr(t.Interface())...)
		if err != nil {
			return result, err
		}
		result = reflect.Append(result, t)
	}
	return result.Interface(), nil
}

func SqlInsert(mysql SqlHandle, tableName string, dao interface{}, skipTags ...string) (lastInsertID int64, err error) {
	allField := GetFieldsTagByKey(dao, "json", skipTags...)
	allValue := GetFieldsValue(dao, skipTags...)

	sqlStr := fmt.Sprintf("INSERT INTO `%v` SET `%v`=? ", tableName,
		strings.Join(allField, "`=?,`"))

	result, err := mysql.Exec(sqlStr, allValue...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		return 0, fmt.Errorf("SqlInsert failed, table:%v affect rows:%v err:%v", tableName, rows, err)
	}
	return result.LastInsertId()
}

func SqlUpsert(mysql SqlHandle, tableName string, dao interface{}, skipTags ...string) (lastInsertID int64, err error) {
	allField := GetFieldsTagByKey(dao, "json", skipTags...)
	allValue := GetFieldsValue(dao, skipTags...)

	sqlStr := fmt.Sprintf("INSERT INTO `%v`(`%v`) VALUES(%v) ON DUPLICATE KEY UPDATE", tableName,
		strings.Join(allField, "`,`"), GetSqlPlaceholder(len(allField)))

	result, err := mysql.Exec(sqlStr, allValue...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		return 0, fmt.Errorf("SqlInsert failed, table:%v affect rows:%v err:%v", tableName, rows, err)
	}
	return result.LastInsertId()
}

func SqlBatchUpdateTag(mysql SqlHandle, tableName string, daoList []interface{}, conditionTag, updateTag string) error {
	if len(daoList) == 0 {
		return nil
	}
	condition := fmt.Sprintf(" Case `%v` ", conditionTag)
	params, keys := make([]interface{}, 0), make([]interface{}, 0)
	for _, dao := range daoList {
		conditionVal, allValue := GetSpecifiedFieldsValueWithSpecialField(dao, conditionTag, updateTag)
		condition += fmt.Sprintf(" WHEN ? THEN ? ")
		params = append(params, conditionVal, allValue[0])
		keys = append(keys, conditionVal)
	}
	params = append(params, keys...)
	condition += fmt.Sprintf(" END WHERE `%v` in (%v)", conditionTag, GetSqlPlaceholder(len(daoList)))
	sqlStr := fmt.Sprintf("UPDATE `%v` SET `%v`= %v ", tableName, updateTag, condition)

	result, err := mysql.Exec(sqlStr, params...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SqlUpdate failed, table:%v affect rows:%v err:%v", tableName, rows, err)
	}
	return nil
}

func SqlUpdateWithUpdateTags(mysql SqlHandle, tableName string, dao interface{}, conditionTag string, updateTags ...string) error {
	allField := GetSpecifiedFieldsTag(dao, "json", updateTags...)
	conditionVal, allValue := GetSpecifiedFieldsValueWithSpecialField(dao, conditionTag, updateTags...)
	allValue = append(allValue, conditionVal)

	//sqlStr := fmt.Sprintf("INSERT INTO %v(`%v`) VALUES(%v)", tableName,
	//	strings.Join(allField, "`,`"), GetSqlPlaceholder(len(allField)))
	sqlStr := fmt.Sprintf("UPDATE `%v` SET `%v`=?  WHERE `%v`=?", tableName,
		strings.Join(allField, "`=?,`"), conditionTag)

	result, err := mysql.Exec(sqlStr, allValue...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SqlUpdate failed, table:%v affect rows:%v err:%v", tableName, rows, err)
	}
	return nil
}

func SqlUpdate(mysql SqlHandle, tableName string, dao interface{}, skipTags ...string) (lastInsertID int64, err error) {
	allField := GetFieldsTagByKey(dao, "json", skipTags...)
	allValue := GetFieldsValue(dao, skipTags...)

	//sqlStr := fmt.Sprintf("INSERT INTO %v(`%v`) VALUES(%v)", tableName,
	//	strings.Join(allField, "`,`"), GetSqlPlaceholder(len(allField)))
	sqlStr := fmt.Sprintf("UPDATE  `%v`(`%v`) SET VALUES(%v)", tableName,
		strings.Join(allField, "`,`"), GetSqlPlaceholder(len(allField)))

	result, err := mysql.Exec(sqlStr, allValue...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		return 0, fmt.Errorf("SqlInsert failed, table:%v affect rows:%v err:%v", tableName, rows, err)
	}
	return result.LastInsertId()
}

func SqlInsertBatch(mysql SqlHandle, tableName string, daoSlice interface{}, skipTags ...string) error {
	refType := reflect.ValueOf(daoSlice)
	sliceLen := refType.Len()
	if sliceLen <= 0 {
		return nil
	}

	exampleDAO := refType.Index(0).Interface()
	idFieldIndex, i := -1, -1
	for _, tag := range GetFieldsTagByKey(exampleDAO, "json") {
		i++
		if tag == "id" {
			idFieldIndex = i
			break
		}
	}

	allField := GetFieldsTagByKey(exampleDAO, "json", skipTags...)
	placeholder := GetSqlPlaceholder(len(allField))

	left, right := 0, 0
	for ; left < sliceLen; left = right {
		right = left + maxDealPerTimes
		if right > sliceLen {
			right = sliceLen
		}
		dealCnt := right - left

		sqlStr := fmt.Sprintf("INSERT INTO %v(`%v`) VALUES %v", tableName,
			strings.Join(allField, "`,`"),
			strings.Repeat(fmt.Sprintf(",(%v)", placeholder), dealCnt)[1:],
		)
		values := make([]interface{}, 0)
		for i := left; i < right; i++ {
			dao := refType.Index(i).Interface()
			values = append(values, GetFieldsValue(dao, skipTags...)...)
		}
		result, err := mysql.Exec(sqlStr, values...)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil || rows != int64(dealCnt) {
			return fmt.Errorf("SqlInsertBatch failed, table:%v err:%v rows should be:%v but real:%v", tableName, err, dealCnt, rows)
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		// 更新ID
		if idFieldIndex >= 0 {
			for i := left; i < right; i++ {
				refType.Index(i).Elem().Field(idFieldIndex).SetUint(uint64(lastID))
				lastID++
			}
		}

	}

	return nil
}

type SqlTxFunc func(tx *sql.Tx) error

func SqlTransaction(mysql *sql.DB, runFunc ...SqlTxFunc) error {
	tx, err := mysql.Begin()
	if err != nil {
		return fmt.Errorf("mysql begin err:%v", err)
	}
	defer func() {
		p := recover()
		if p != nil {
			logger.Error(logTagMysqlUtil, "run sql recover:%v", p)
			err := tx.Rollback()
			if err != nil {
				logger.Error(logTagMysqlUtil, "mysql rollback err:%v", err)
			}
			time.Sleep(time.Second)
			panic(p)
		}
	}()

	var runErr error = nil
	for _, run := range runFunc {
		runErr = run(tx)
		if runErr != nil {
			break
		}
	}

	if runErr == nil {
		err := tx.Commit()
		if err != nil {
			logger.Error(logTagMysqlUtil, "mysql commit err:%v", err)
			return err
		}
	} else {
		logger.Warn(logTagMysqlUtil, "mysql rollback, transaction err:%v", runErr)
		err := tx.Rollback()
		if err != nil {
			logger.Error(logTagMysqlUtil, "mysql rollback err:%v", err)
		}
	}
	return runErr
}

func GetIDLowerBound(mysql SqlHandle, tableName string, cmpTime time.Time) (int64, error) {
	type D struct {
		ID    int64     `json:"id"`
		Ctime time.Time `json:"ctime"`
	}
	d := &D{}
	var L, R int64
	err := SqlQueryRow(mysql, tableName, d, "ORDER BY `id` ASC LIMIT 1")
	if err != nil {
		return 0, err
	}
	if d.Ctime.Before(cmpTime) == false {
		return d.ID, nil
	}
	L = d.ID

	err = SqlQueryRow(mysql, tableName, d, "ORDER BY `id` DESC LIMIT 1")
	if err != nil {
		return 0, err
	}
	if d.Ctime.Before(cmpTime) == true {
		return d.ID + 1, nil
	}
	R = d.ID

	for round := 0; L < R; round++ {
		if round >= 50 {
			return 0, fmt.Errorf("GetIDLowerBound error, round >= 50, args:%v,%v", tableName, cmpTime)
		}

		err = SqlQueryRow(mysql, tableName, d, "WHERE `id`<=? ORDER BY `id` DESC LIMIT 1", L+(R-L)/2)
		if err != nil {
			return 0, err
		}
		if d.Ctime.Before(cmpTime) == false {
			R = d.ID
		} else {
			if d.ID >= L {
				L = d.ID + 1
			} else {

				// 由于id不连续，多加的判断代码
				err = SqlQueryRow(mysql, tableName, d, "WHERE `id`>=? ORDER BY `id` ASC LIMIT 1", L+(R-L+1)/2)
				if err != nil {
					return 0, err
				}
				if d.ID >= R {
					return R, nil
				}

				if d.Ctime.Before(cmpTime) == false {
					R = d.ID
				} else {
					L = d.ID + 1
				}
			}
		}
	}
	return L, nil
}
