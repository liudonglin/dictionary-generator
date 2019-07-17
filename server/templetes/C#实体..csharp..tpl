using System;

namespace entity
{
	public class {{ fn.ToCamelString(table.Name) }} 
	{
		{% for column in columns %}
		public {{ fn.SqlTypeConvertLanguageType(column,project.DataBase,project.Language) }} {{ fn.ToCamelString(column.Name) }} { get; set; }
		{% endfor %}
	}
}
