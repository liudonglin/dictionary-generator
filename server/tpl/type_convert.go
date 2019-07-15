package tpl

import "strings"

func javaTypeConvert(_type string) string {
	result := ""
	upType := strings.ToUpper(_type)
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
