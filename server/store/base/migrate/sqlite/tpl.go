package sqlite

import (
	"database/sql"
	"dg-server/core"
	"time"
)

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
	{{ toCamelString(column.Name) }} string ` + "`" + `json:"{{ toSnakeString(column.Name) }}"` + "`" + `{% endfor %}
}
`

var migrationTempletes = []core.Templete{
	{
		Name:     "mysql数据表创建脚本",
		Content:  mysqlTableCreate,
		DataBase: "mysql",
		Type:     "table_create",
	},
	{
		Name:     "Golang实体",
		Content:  golangEntity,
		Type:     "code_entity",
		Language: "go",
	},
}

// MigrateTemplete 初始化Templete
func initTemplete(db *sql.DB) error {
	completed, _ := selectCompletedTpl(db)
	for _, mtpl := range migrationTempletes {

		if _, ok := completed[mtpl.Name]; ok {
			continue
		}

		if err := insertTemplete(db, &mtpl); err != nil {
			return err
		}
	}
	return nil
}

func insertTemplete(db *sql.DB, tpl *core.Templete) error {
	created := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(migrationTempleteInsert,
		tpl.Name,
		tpl.Content,
		tpl.Language,
		tpl.DataBase,
		tpl.Orm,
		tpl.Type,
		created,
		created)
	return err
}

func selectCompletedTpl(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelectTpl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

const migrationTempleteInsert = `
INSERT INTO templetes (
templete_name
,templete_content
,templete_language
,templete_data_base
,templete_orm
,templete_type
,templete_created
,templete_updated
) VALUES (
?
,?
,?
,?
,?
,?
,?
,?
)
`

var migrationSelectTpl = `
SELECT templete_name FROM templetes
`
