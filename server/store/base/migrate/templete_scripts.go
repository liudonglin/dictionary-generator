package migrate

const mysqlTableCreate = `
USE {{database.Name}} ;
Drop Table If Exists {{table.Name}} ;
CREATE TABLE {{table.Name}} ( {% for column in columns %}
{{column.Name}} {{column.ColumnType}} {% if column.Null %}{% else %}NOT{% endif %} NULL {% if column.AI %}AUTO_INCREMENT{% endif %} COMMENT '{{column.Title}}' {% if !isLastColumn(column,columns) %},{% endif %}{% endfor %}{% if lenColumn(indexs)>0 %},{% endif %}
{% for column in indexs %}
{% if column.PK %}PRIMARY {% elif column.Unique %}UNIQUE {% endif %}KEY {% if !column.PK %}idx_{{table.Name}}_{{column.Name}}{% endif %} ({{column.Name}}) {% if !isLastColumn(column,indexs) %},{% endif %} {% endfor %}
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='{{database.Description}}';
`

const golangEntity = `
package entity

type {{ toCamelString(table.Name) }} struct { {% for column in columns %}
	{{ toCamelString(column.Name) }} {{ sqlTypeConvertLanguageType(column,project.DataBase,project.Language) }} ` + "`" + `json:"{{ toSnakeString(column.Name) }}"` + "`" + `{% endfor %}
}
`

const javaEntity = `
package entity

import lombok.Data;

@Data
public class {{ toCamelString(table.Name) }} { {% for column in columns %}
	private {{ sqlTypeConvertLanguageType(column,project.DataBase,project.Language) }} {{ toCamelString(column.Name) }};  
{% endfor %}
}
`

// MigrationTempletes ...
var MigrationTempletes = []core.Templete{
	{
		Name:     "mysql数据表创建脚本",
		Content:  mysqlTableCreate,
		DataBase: "mysql",
		Type:     "init",
	},
	{
		Name:     "Golang实体",
		Content:  golangEntity,
		Type:     "init",
		Language: "go",
	},
	{
		Name:     "Java实体",
		Content:  javaEntity,
		Type:     "init",
		Language: "java",
	},
}