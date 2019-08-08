package handler

import (
	"bytes"
	"dg-server/core"
	"dg-server/store"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/labstack/echo"
)

func exportProject(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	id, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return err
	}

	projectStore := store.Stores().ProjectStore
	project, err := projectStore.FindID(id)
	if err != nil {
		return err
	}
	if project.ID == 0 {
		return &BusinessError{Message: "传入项目编号无效"}
	}

	name := project.Name + ".xlsx"
	name = base64.StdEncoding.EncodeToString([]byte(name))

	f := excelize.NewFile()

	err = buildExcelSheets(f, id)
	if err != nil {
		return err
	}

	//将数据存入buff中
	var buff bytes.Buffer
	if err := f.Write(&buff); err != nil {
		panic(err)
	}
	//设置请求头  使用浏览器下载
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+name)
	return c.Stream(http.StatusOK, echo.MIMEOctetStream, bytes.NewReader(buff.Bytes()))
}

func buildExcelSheets(f *excelize.File, pid int64) error {
	q := &core.DBQuery{
		PID: pid,
		Pager: core.Pager{
			Index: 0,
			Size:  9999999,
		},
	}

	dbStore := store.Stores().DataBaseStore
	list, _, err := dbStore.List(q)
	if err != nil {
		return err
	}

	/* 初始化通用样式 */
	styles := make(map[string]int)
	style, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#999999"],"pattern":1},"font":{"bold":true,"size":12,"family":"宋体"},"border":[{"type":"top","color":"#444444","style":1},{"type":"bottom","color":"#444444","style":1},{"type":"right","color":"#444444","style":1}]}`)
	styles["header-table-cell"] = style
	style, _ = f.NewStyle(`{"font":{"bold":true,"size":12,"family":"宋体"},"border":[{"type":"top","color":"#444444","style":1},{"type":"bottom","color":"#444444","style":1},{"type":"right","color":"#444444","style":1}]}`)
	styles["header-common-cell"] = style
	style, _ = f.NewStyle(`{"fill":{"type":"pattern","color":["#99CCFF"],"pattern":1},"font":{"bold":true,"size":12,"family":"宋体"},"border":[{"type":"top","color":"#444444","style":1},{"type":"bottom","color":"#444444","style":1},{"type":"right","color":"#444444","style":1}]}`)
	styles["header-field-cell"] = style

	for _, item := range list {
		// Create a new sheet.
		f.NewSheet(item.Name)
		if err = buildExcelTables(f, item, styles); err != nil {
			return err
		}
	}
	return nil
}

func buildExcelTables(f *excelize.File, _db *core.DataBase, styles map[string]int) error {
	q := &core.TableQuery{
		DID: _db.ID,
		Pager: core.Pager{
			Index: 0,
			Size:  9999999,
		},
	}

	cunrrentRowIndex := 1
	sheet := _db.Name

	tableStore := store.Stores().TableStore
	tables, _, err := tableStore.List(q)
	if err != nil {
		return err
	}

	columnStore := store.Stores().ColumnStore
	columns, _, _ := columnStore.List(&core.ColumnQuery{
		PID: _db.PID,
		DID: _db.ID,
		Pager: core.Pager{
			Index: 0,
			Size:  9999999,
		},
	})

	for _, table := range tables {
		table.Columns = make([]*core.Column, 0)
		for _, column := range columns {
			if column.TID == table.ID {
				table.Columns = append(table.Columns, column)
			}
		}
	}

	/* 设置列宽 */
	f.SetColWidth(sheet, "A", "A", 10)
	f.SetColWidth(sheet, "B", "B", 30)
	f.SetColWidth(sheet, "C", "D", 30)
	f.SetColWidth(sheet, "E", "F", 12)
	f.SetColWidth(sheet, "H", "H", 20)
	f.SetColWidth(sheet, "I", "I", 30)
	f.SetColWidth(sheet, "J", "J", 10)
	f.SetColWidth(sheet, "K", "K", 30)

	for _, table := range tables {
		/* 构建表头 */
		buildTableHeader(f, sheet, table, &cunrrentRowIndex, styles)

		/* 构建表体 */
		buildTableBody(f, sheet, table.Columns, &cunrrentRowIndex, styles)

		/* 每张表间隔一行 */
		f.SetRowHeight(sheet, cunrrentRowIndex, 18)
		cunrrentRowIndex++
	}

	return nil
}

func buildTableHeader(f *excelize.File, sheet string, table *core.Table, cunrrentRowIndex *int, styles map[string]int) {

	f.SetRowHeight(sheet, *cunrrentRowIndex, 18)
	f.SetCellValue(sheet, fmt.Sprintf("A%d", *cunrrentRowIndex), "Table")
	style, _ := styles["header-table-cell"]
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *cunrrentRowIndex), fmt.Sprintf("A%d", *cunrrentRowIndex), style)

	style, _ = styles["header-common-cell"]
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", *cunrrentRowIndex), fmt.Sprintf("K%d", *cunrrentRowIndex), style)

	f.SetCellValue(sheet, fmt.Sprintf("B%d", *cunrrentRowIndex), table.Name)
	f.SetCellValue(sheet, fmt.Sprintf("C%d", *cunrrentRowIndex), table.Title)
	f.SetCellValue(sheet, fmt.Sprintf("D%d", *cunrrentRowIndex), table.Description)
	// err := f.MergeCell(sheet, fmt.Sprintf("D%d", *cunrrentRowIndex), fmt.Sprintf("K%d", cunrrentRowIndex))
	// if err != nil {

	// }

	*cunrrentRowIndex++
	f.SetRowHeight(sheet, *cunrrentRowIndex, 18)

	f.SetCellValue(sheet, fmt.Sprintf("A%d", *cunrrentRowIndex), "主键")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", *cunrrentRowIndex), "自增")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", *cunrrentRowIndex), "列名")
	f.SetCellValue(sheet, fmt.Sprintf("D%d", *cunrrentRowIndex), "数据类型")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", *cunrrentRowIndex), "长度")
	f.SetCellValue(sheet, fmt.Sprintf("F%d", *cunrrentRowIndex), "可空")
	f.SetCellValue(sheet, fmt.Sprintf("G%d", *cunrrentRowIndex), "索引列")
	f.SetCellValue(sheet, fmt.Sprintf("H%d", *cunrrentRowIndex), "枚举")
	f.SetCellValue(sheet, fmt.Sprintf("I%d", *cunrrentRowIndex), "标题")
	f.SetCellValue(sheet, fmt.Sprintf("J%d", *cunrrentRowIndex), "唯一列")
	f.SetCellValue(sheet, fmt.Sprintf("K%d", *cunrrentRowIndex), "描述")
	style, _ = styles["header-field-cell"]
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *cunrrentRowIndex), fmt.Sprintf("K%d", *cunrrentRowIndex), style)

	*cunrrentRowIndex++
}

func buildTableBody(f *excelize.File, sheet string, columns []*core.Column, cunrrentRowIndex *int, styles map[string]int) {
	for _, col := range columns {
		f.SetRowHeight(sheet, *cunrrentRowIndex, 18)
		style, _ := styles["header-common-cell"]
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", *cunrrentRowIndex), fmt.Sprintf("K%d", *cunrrentRowIndex), style)

		if col.PK {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", *cunrrentRowIndex), "PK")
		}
		if col.AI {
			f.SetCellValue(sheet, fmt.Sprintf("B%d", *cunrrentRowIndex), "AI")
		}

		f.SetCellValue(sheet, fmt.Sprintf("C%d", *cunrrentRowIndex), col.Name)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", *cunrrentRowIndex), col.DataType)
		f.SetCellStr(sheet, fmt.Sprintf("E%d", *cunrrentRowIndex), col.Length)

		if col.Null {
			f.SetCellValue(sheet, fmt.Sprintf("F%d", *cunrrentRowIndex), "Null")
		} else {
			f.SetCellValue(sheet, fmt.Sprintf("F%d", *cunrrentRowIndex), "Not Null")
		}

		if col.Index {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", *cunrrentRowIndex), "INDEX")
		}

		f.SetCellValue(sheet, fmt.Sprintf("H%d", *cunrrentRowIndex), col.Enum)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", *cunrrentRowIndex), col.Title)

		if col.Index {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", *cunrrentRowIndex), "UNI")
		}

		f.SetCellValue(sheet, fmt.Sprintf("K%d", *cunrrentRowIndex), col.Description)

		*cunrrentRowIndex++
	}
}
