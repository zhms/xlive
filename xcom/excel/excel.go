package excel

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm/utils"
)

//用法示例
/*
	e := excel.NewExcelBuilder("_log", "abc")
	defer e.Close()
	e.SetTitle("Id", "订单号")
	e.SetValue("Id", "fuck", 2)
	filePath, _ := e.Save()
	return "/exports/"+path.Base(filePath)
*/

type ExcelBuilder struct {
	fileName       string
	sheetName      string
	excelfile      *excelize.File
	columnIndex    int
	titleColumnMap map[string]string
	titleOptions   map[string]map[string]interface{}
}

func NewExcelBuilder(fileName string) *ExcelBuilder {
	r := &ExcelBuilder{
		fileName:       fileName,
		sheetName:      "Sheet1",
		excelfile:      nil,
		columnIndex:    0,
		titleColumnMap: make(map[string]string),
		titleOptions:   make(map[string]map[string]interface{}),
	}
	r.Open()
	return r
}

func (s *ExcelBuilder) GetColumnName(columnNumber int) string {
	if columnNumber <= 26 {
		return fmt.Sprintf("%c", 'A'-1+columnNumber)
	} else {
		first := columnNumber / 26
		second := columnNumber % 26
		if first < 26 {
			return fmt.Sprintf("%c%c", 'A'-1+first, 'A'-1+second)
		} else {
			return s.GetColumnName(first) + fmt.Sprintf("%c", 'A'-1+second)
		}
	}
}

func (s *ExcelBuilder) SetTitleStyle() (err error) {
	if s.excelfile == nil {
		err = errors.New("excelize File did not open")
		return
	}
	//设置表头固定
	s.excelfile.SetPanes(s.sheetName, &excelize.Panes{
		Freeze:     true,
		Split:      false,
		XSplit:     0,
		YSplit:     1,
		ActivePane: "bottomLeft",
	})
	boldStyle, _ := s.excelfile.NewStyle(
		&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Font: &excelize.Font{
				Bold: true,
			},
		},
	)
	s.excelfile.SetCellStyle(s.sheetName, "A1", s.GetColumnName(s.columnIndex)+"1", boldStyle)
	return
}

func (s *ExcelBuilder) SetValueStyle(row int64) (err error) {
	if s.excelfile == nil {
		err = errors.New("excelize File did not open")
		return
	}
	centerStyle, _ := s.excelfile.NewStyle(
		&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		},
	)
	s.excelfile.SetCellStyle(s.sheetName, "A2", s.GetColumnName(s.columnIndex)+fmt.Sprintf("%d", row), centerStyle)
	return
}

func (s *ExcelBuilder) Open() {
	if s.excelfile == nil {
		s.excelfile = excelize.NewFile()
	}
}

func (s *ExcelBuilder) Close() {
	if s.excelfile != nil {
		s.excelfile.Close()
		s.excelfile = nil
		s.columnIndex = 0
		s.titleColumnMap = make(map[string]string)
	}
}

func (s *ExcelBuilder) SetTitle(key, title string) (err error) {
	if s.excelfile == nil {
		err = errors.New("excelize File did not open")
		return
	}
	if key == "" {
		err = errors.New("key can not be ''")
		return
	}
	s.columnIndex++
	indexName := s.GetColumnName(s.columnIndex)
	s.titleColumnMap[key] = indexName
	err = s.excelfile.SetCellValue(s.sheetName, indexName+"1", title)
	return
}
func (s *ExcelBuilder) SetTitleEx(key, title string, options map[string]interface{}) (err error) {
	if s.excelfile == nil {
		err = errors.New("excelize File did not open")
		return
	}
	if key == "" {
		err = errors.New("key can not be ''")
		return
	}
	s.columnIndex++
	indexName := s.GetColumnName(s.columnIndex)
	s.titleColumnMap[key] = indexName
	s.titleOptions[key] = options
	err = s.excelfile.SetCellValue(s.sheetName, indexName+"1", title)
	return
}

func (s *ExcelBuilder) GetTitleColumn(key string) (column string, ok bool) {
	column, ok = s.titleColumnMap[key]
	return
}

func (s *ExcelBuilder) SetColumnWidth(key string, width float64) (err error) {
	if s.excelfile == nil {
		err = errors.New("excelize File did not open")
		return
	}
	if column, ok := s.titleColumnMap[key]; ok {
		err = s.excelfile.SetColWidth(s.sheetName, column, column, width)
	} else {
		err = errors.New("key not exist")
	}
	return
}

func (s *ExcelBuilder) SetRowHeight(row int, height float64) (err error) {
	if s.excelfile == nil {
		err = errors.New("excelize File did not open")
		return
	}
	s.excelfile.SetRowHeight(s.sheetName, row, height)
	return
}

func (s *ExcelBuilder) SetValue(key string, value interface{}, row int64) (err error) {
	if s.excelfile == nil {
		err = errors.New("excelize File did not open")
		return
	}
	column, ok := s.titleColumnMap[key]
	if !ok {
		err = errors.New("column not exist")
		return
	}
	if options, ok := s.titleOptions[key]; ok {
		v := utils.ToString(value)
		vf, vok := options[v]
		if vok {
			err = s.excelfile.SetCellValue(s.sheetName, column+fmt.Sprintf("%d", row), vf)
		} else {
			err = s.excelfile.SetCellValue(s.sheetName, column+fmt.Sprintf("%d", row), value)
		}
	}
	err = s.excelfile.SetCellValue(s.sheetName, column+fmt.Sprintf("%d", row), value)
	return
}

func (s *ExcelBuilder) Write(ctx *gin.Context) {
	if s.excelfile == nil {
		return
	}
	n := url.PathEscape(s.fileName)
	ctx.Header("Content-Disposition", "attachment; filename="+n+".xlsx")
	ctx.Header("Content-Type", "application/octet-stream")
	s.excelfile.Write(ctx.Writer)
}
