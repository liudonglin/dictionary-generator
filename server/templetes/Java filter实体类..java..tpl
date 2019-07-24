package {{ project.NameSpace }}.req;

import com.dk.foundation.engine.baseentity.Pager;
import lombok.Data;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;

import java.io.Serializable;

@Data
public class {{ fn.ToCamelString(table.Name) }}Filter extends Pager implements Serializable {

{ {% for column in columns %}
    @ApiModelProperty("{{ column.Title }}")	{% if !column.Null %}
    @NotNull(message = "{{ column.Title }}不能为空") {% endif %} {% if column.DataType=="varchar" %}
    @Length(max = {{column.Length}},message = "{{ column.Title }}长度不能超过{{column.Length}}") {% endif %} 
    private {% if column.IsEnum %}{{ fn.ToCamelString(column.Name) }}{% else %}{{ fn.SqlTypeConvertLanguageType(column,project.DataBase,project.Language) }}{% endif %} {{ fn.FirstToLower(fn.ToCamelString(column.Name)) }};  
{% endfor %}

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this, ToStringStyle.SHORT_PREFIX_STYLE);
    }
}
