package tpl

import (
	"dg-server/core"
	"strings"
)

// SqlTypeConvertLanguageType 数据库类型转编程语言类型
func (*FnWrap) SqlTypeConvertLanguageType(col *core.Column, dataBase, language string) string {
	result := ""

	if dataBase == "mysql" && language == "java" {
		result = mysqlConvertJava(col)
	}
	if dataBase == "mysql" && language == "go" {
		result = mysqlConvertGo(col)
	}
	if dataBase == "mysql" && language == "csharp" {
		result = mysqlConvertCsharp(col)
	}
	return result
}

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
	case "DATE", "DATETIME", "TIMESTAMP":
		result = "Date"
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

func mysqlConvertCsharp(col *core.Column) string {
	result := ""
	upType := strings.ToUpper(col.DataType)
	switch upType {
	case "TINYINT":
		result = "bool?"
	case "INT", "INTEGER", "SMALLINT":
		result = "int?"
	case "BIGINT":
		result = "long?"
	case "FLOAT":
		result = "float?"
	case "DOUBLE":
		result = "double?"
	case "DECIMAL":
		result = "decimal?"
	case "DATE", "DATETIME", "TIMESTAMP":
		result = "DateTime?"
	case "CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT", "ENUM", "SET":
		result = "string"
	}
	return result
}
