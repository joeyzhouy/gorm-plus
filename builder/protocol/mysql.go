package protocol

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm-plus/extension"
	"gorm-plus/utils"
	"strings"
)

const mysql = "mysql"
const mysqlConnectUrl = "%s:%s@tcp(%s:%d)/%s"
const showTables = "show tables"
const yes = "YES"
const PRI = "PRI"
const columnQuery = "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE, COLUMN_COMMENT,TABLE_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' AND table_name IN ('%s') ORDER BY TABLE_NAME"

func init() {
	extension.RegisterProtocol(mysql, func(param utils.Param) (extension.Protocol, error) {
		return &mysqlProtocol{
			Param: param,
		}, nil
	})
}

type mysqlProtocol struct {
	Param utils.Param
	db    *sql.DB
}

func (m *mysqlProtocol) GetColumnWithTableNames(tableNames []string) ([]utils.TableInfo, error) {
	if len(tableNames) == 0 {
		return nil, nil
	}
	result := make([]utils.TableInfo, 0)
	row, err := m.db.Query(fmt.Sprintf(columnQuery, m.Param.DBName, strings.Join(tableNames, "','")))
	if err != nil {
		return nil, err
	}
	var name, key, dataType, nullable, comment, tableName, preName string
	var tempColumnInfo utils.ColumnInfo
	columns := make([]utils.ColumnInfo, 0)
	for row.Next() {
		row.Scan(&name, &key, &dataType, &nullable, &comment, &tableName)
		if preName == "" {
			preName = tableName
		} else if preName != tableName {
			result = append(result, utils.TableInfo{Name: preName, Columns: columns})
			columns = make([]utils.ColumnInfo, 0)
			preName = tableName
		}
		tempColumnInfo = utils.ColumnInfo{
			Name:     name,
			Comment:  comment,
			DataType: dataType,
			Key:      key}
		if nullable == yes {
			tempColumnInfo.Nullable = true
		} else {
			tempColumnInfo.Nullable = false
		}
		tempColumnInfo.GoType = mapping(tempColumnInfo.DataType, tempColumnInfo.Nullable)
		if tempColumnInfo.Key == PRI {
			tempColumnInfo.Key = utils.GormPrimary
		} else {
			tempColumnInfo.Key = ""
		}
		columns = append(columns, tempColumnInfo)
	}
	if len(columns) > 0 {
		result = append(result, utils.TableInfo{Name: preName, Columns: columns})
	}
	return result, nil
}

func mapping(mysqlType string, nullable bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			return utils.SqlNullInt32
		}
		return utils.GoInt
	case "bigint":
		if nullable {
			return utils.SqlNullInt64
		}
		return utils.GoInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			return utils.SqlNullString
		}
		return "string"
	case "date", "datetime", "time", "timestamp":
		if nullable {
			return utils.SqlNullTime
		}
		return utils.GoTime
	case "decimal", "double":
		if nullable {
			return utils.SqlNullFloat64
		}
		return utils.GoFloat64
	case "float":
		if nullable {
			return utils.SqlNullFloat64
		}
		return utils.GoFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return utils.GoByteArray
	}
	return ""
}

func (m *mysqlProtocol) GetTableNames() ([]string, error) {
	rows, err := m.db.Query(showTables)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	var tableName string
	for rows.Next() {
		rows.Scan(&tableName)
		result = append(result, tableName)
	}
	return result, nil
}

func (m *mysqlProtocol) GetConnection() (*sql.DB, error) {
	var err error
	connectUrl := fmt.Sprintf(mysqlConnectUrl, m.Param.UserName, m.Param.Password,
		m.Param.Host, m.Param.Port, m.Param.DBName)
	m.db, err = sql.Open(mysql, connectUrl)
	return m.db, err
}
