USE `{{database.Name}}` ;
Drop Table If Exists `{{table.Name}}` ;
CREATE TABLE `{{table.Name}}` ( {% for column in columns %}
`{{column.Name}}` {{column.ColumnType}} {% if column.Null %}NULL{% else %}NOT NULL {% if column.DataType=="timestamp" %}DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP{% endif %}{% endif %} {% if column.AI %}AUTO_INCREMENT{% endif %} COMMENT '{{column.Title}}' {% if !fn.IsLastColumn(column,columns) %},{% endif %}{% endfor %}{% if fn.LenColumn(indexs)>0 %},{% endif %}
{% for column in indexs %}
{% if column.PK %}PRIMARY {% elif column.Unique %}UNIQUE {% endif %}KEY {% if !column.PK %}idx_{{table.Name}}_{{column.Name}}{% endif %} ({{column.Name}}) {% if !fn.IsLastColumn(column,indexs) %},{% endif %} {% endfor %}
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='{{table.Title}}';
