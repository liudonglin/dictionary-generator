package tpl

import (
	"dg-server/core"
	"dg-server/store"
	"strings"
	"time"

	"github.com/flosch/pongo2"
)

var tableTemplete = `
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
`

// TestGetTableScript ...
func TestGetTableScript(tid int64) string {
	tpl, err := pongo2.FromString(tableTemplete)

	tableStore := store.Stores().TableStore
	table, _ := tableStore.FindID(tid)

	dbStore := store.Stores().DataBaseStore
	database, _ := dbStore.FindID(table.DID)

	projectStore := store.Stores().ProjectStore
	project, _ := projectStore.FindID(table.PID)

	columnStore := store.Stores().ColumnStore
	columns, _, _ := columnStore.List(&core.ColumnQuery{
		TID: table.ID,
		Pager: core.Pager{
			Index: 0,
			Size:  9999999,
		},
	})

	indexs := make([]*core.Column, 0)
	for _, column := range columns {
		if column.Index {
			indexs = append(indexs, column)
		}
	}

	out, err := tpl.Execute(pongo2.Context{
		"project":  project,
		"database": database,
		"table":    table,
		"columns":  columns,
		"indexs":   indexs,
		"user":     "liudonglin",
		"fn":       fn,
	})

	if err != nil {
		print(err.Error())
	} else {
		print(out)
	}

	return out
}

// GetTableScript ...
func GetTableScript(req *core.TempleteLoadReq, userName string) (string, error) {
	tplStore := store.Stores().TempleteStore
	temp, _ := tplStore.FindID(req.TempleteID)
	tpl, err := pongo2.FromString(temp.Content)

	tableStore := store.Stores().TableStore
	table, _ := tableStore.FindID(req.TID)

	dbStore := store.Stores().DataBaseStore
	database, _ := dbStore.FindID(table.DID)

	projectStore := store.Stores().ProjectStore
	project, _ := projectStore.FindID(table.PID)

	columnStore := store.Stores().ColumnStore
	columns, _, _ := columnStore.List(&core.ColumnQuery{
		TID: table.ID,
		Pager: core.Pager{
			Index: 0,
			Size:  9999999,
		},
	})

	indexs := make([]*core.Column, 0)
	for _, column := range columns {
		if column.Index {
			indexs = append(indexs, column)
		}
	}

	out, err := tpl.Execute(pongo2.Context{
		"project":  project,
		"database": database,
		"table":    table,
		"columns":  columns,
		"indexs":   indexs,
		"user":     userName,
		"fn":       fn,
	})

	if err != nil {
		return "", err
	}
	return out, nil
}

// FnWrap ...
type FnWrap struct{}

var fn = &FnWrap{}

// LenColumn 获取列长度
func (*FnWrap) LenColumn(arr []*core.Column) int {
	return len(arr)
}

// IsLastColumn 判断是否是最后一列
func (*FnWrap) IsLastColumn(item *core.Column, list []*core.Column) bool {
	if item.ID == list[len(list)-1].ID {
		return true
	}
	return false
}

// ToSnakeString snake string, XxYy to xx_yy , XxYY to xx_yy
func (*FnWrap) ToSnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// ToCamelString camel string, xx_yy to XxYy
func (*FnWrap) ToCamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// FirstToLower 首字母转小写
func (*FnWrap) FirstToLower(s string) string {
	if len(s) < 1 {
		return ""
	}
	strArry := []rune(s)
	if strArry[0] >= 65 && strArry[0] <= 90 {
		strArry[0] += 32
	}
	return string(strArry)
}

// FirstToUpper 首字母转大写
func (*FnWrap) FirstToUpper(s string) string {
	if len(s) < 1 {
		return ""
	}
	strArry := []rune(s)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}

// GetNowDate 获取当前日期
func (*FnWrap) GetNowDate() string {
	return time.Now().Format("2006-01-02")
}
