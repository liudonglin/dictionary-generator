package entity

type {{ fn.ToCamelString(table.Name) }} struct { {% for column in columns %}
	{{ fn.ToCamelString(column.Name) }} {{ fn.SqlTypeConvertLanguageType(column,project.DataBase,project.Language) }} `json:"{{ fn.ToSnakeString(column.Name) }}"`{% endfor %}
}
