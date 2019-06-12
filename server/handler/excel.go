package handler

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func testExcel() {
    f := excelize.NewFile()
    // Create a new sheet.
    index := f.NewSheet("Sheet2")
    // Set value of a cell.
    f.SetCellValue("Sheet2", "A2", "Hello world.")
    f.SetCellValue("Sheet1", "B2", 100)
    // Set active sheet of the workbook.
    f.SetActiveSheet(index)
    // Save xlsx file by the given path.
    err := f.SaveAs("./Book1.xlsx")
    if err != nil {
        fmt.Println(err)
    }
}