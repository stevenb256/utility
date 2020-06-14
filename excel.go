package utl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	l "github.com/stevenb256/log"
	"github.com/tealeg/xlsx"
)

// Excel used to hold an excel file and columns
type Excel struct {
	i       int
	row     *xlsx.Row
	sheet   *xlsx.Sheet
	columns map[string]int
}

// OpenExcel opens an excel file
func OpenExcel(reader io.Reader, columns []string) (*Excel, error) {

	// locals
	var excel Excel

	// read the whole body
	data, err := ioutil.ReadAll(reader)
	if l.Check(err) {
		return nil, err
	}

	// open the binary
	file, err := xlsx.OpenBinary(data)
	if l.Check(err) {
		return nil, l.Fail(fmt.Errorf("can't OpenBinary on excel file: %s", err.Error()))
	}

	// if no sheets
	if 0 == len(file.Sheets) {
		return nil, l.Fail(errors.New("no sheets in excel file"))
	}

	// get first sheet
	excel.sheet = file.Sheets[0]

	// get row
	excel.row = excel.sheet.Row(0)

	// if no rows
	if nil == excel.row {
		return nil, l.Fail(errors.New("sheet has no rows or no header row"))
	}

	// make column map
	excel.columns = make(map[string]int)

	// find columns
	for _, name := range columns {
		for i := 0; i < len(excel.row.Cells); i++ {
			if name == strings.ToLower(excel.row.Cells[i].Value) {
				excel.columns[name] = i
				break
			}
		}
	}

	// not the same length?
	if len(excel.columns) != len(columns) {
		return nil, l.Fail(fmt.Errorf("spreadsheet must have column headers: %s", strings.Join(columns, ", ")))
	}

	// done
	return &excel, nil
}

// IsDone returns true if no more rows
func (e *Excel) IsDone() bool {
	if nil != e.row {
		e.i++
		e.row = e.sheet.Row(e.i)
	}
	if nil == e.row || 0 == len(e.row.Cells) {
		return true
	}
	return false
}

// getIndex returns index of column name or -1 if index is not valid
func (e *Excel) getIndex(name string) int {
	i, found := e.columns[name]
	if false == found {
		return -1
	}
	if nil == e.row || i >= len(e.row.Cells) {
		return -1
	}
	return i
}

// String gets current row/column as string
func (e *Excel) String(name string) string {
	i := e.getIndex(name)
	if -1 == i {
		return ""
	}
	return e.row.Cells[i].Value
}

// Int gets current row/column as string
func (e *Excel) Int(name string) int {
	i := e.getIndex(name)
	if -1 == i {
		return 0
	}
	return Atoi(e.row.Cells[i].Value)
}

// Date gets current row column as date
func (e *Excel) Date(name string) (time.Time, error) {
	i := e.getIndex(name)
	if -1 == i {
		var t time.Time
		return t, nil
	}
	return e.row.Cells[i].GetTime(false)
}
