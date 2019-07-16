package tpl

import (
	"dg-server/core"
	"strings"
)

func mysqlConvertJava(col *core.Column) string {
	result := ""
	upType := strings.ToUpper(col.DataType)
	switch upType {
	case "BIT", "BINARY", "VARBINARY", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
		result = "byte[]"
	case "TINYINT":
		result = "Boolean"
	case "SMALLINT", "MEDIUMINT", "INT", "INTEGER":
		result = "Integer"
	case "BIGINT":
		result = "Long"
	case "FLOAT":
		result = "Float"
	case "DOUBLE":
		result = "Double"
	case "DECIMAL":
		result = "BigDecimal"
	case "DATE":
		result = "Date"
	case "DATETIME", "TIMESTAMP":
		result = "Timestamp"
	case "CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT", "ENUM", "SET":
		result = "String"
	}
	return result
}

func mysqlConvertGo(col *core.Column) string {
	result := ""
	upType := strings.ToUpper(col.DataType)
	switch upType {
	case "TINYINT", "INT", "INTEGER", "SMALLINT":
		result = "int32"
	case "BIGINT":
		result = "int64"
	case "FLOAT":
		result = "float32"
	case "DOUBLE":
		result = "float64"
	case "DATE", "DATETIME", "TIMESTAMP":
		result = "time.Time"
	case "CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT", "ENUM", "SET", "DECIMAL":
		result = "string"
	}
	return result
}
