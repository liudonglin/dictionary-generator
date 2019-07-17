package {{ project.NameSpace }}.entity

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;
import org.hibernate.validator.constraints.Length;
import javax.validation.constraints.NotNull;

/**
 * {{ table.Description }}
 * @author {{ user }}
 * @date {{ fn.GetNowDate() }}
 * */

@Data
public class {{ fn.ToCamelString(table.Name) }} { {% for column in columns %}
	
	@ApiModelProperty("{{ column.Title }}")	{% if !column.Null %}
	@NotNull(message = "{{ column.Title }}不能为空") {% endif %} {% if column.DataType=="varchar" %}
	@Length(max = {{column.Length}},message = "{{ column.Title }}长度不能超过{{column.Length}}") {% endif %} 
	private {{ fn.SqlTypeConvertLanguageType(column,project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(column.Name)) }};  
{% endfor %}
}