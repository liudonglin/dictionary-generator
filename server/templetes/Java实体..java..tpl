package entity

import lombok.Data;

@Data
public class {{ fn.ToCamelString(table.Name) }} { {% for column in columns %}
	private {{ fn.SqlTypeConvertLanguageType(column,project.DataBase,project.Language) }} {{ fn.ToCamelString(column.Name) }};  
{% endfor %}
}
