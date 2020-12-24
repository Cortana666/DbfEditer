package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/tadvi/dbf"

	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

type modelHandler struct {
	dbfresource *dbf.DbfTable
	tabletitle  []dbf.DbfField
	lines       int
	filename    string
}

func newModelHandler(filename string) *modelHandler {
	dbf, _ := dbf.LoadFile(filename)

	m := new(modelHandler)
	m.dbfresource = dbf
	m.tabletitle = dbf.Fields()
	m.lines = int(dbf.NumRecords())
	m.filename = filename

	return m
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{}
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	mh.dbfresource.SetFieldValue(row, column, string(value.(ui.TableString)))
	mh.dbfresource.SaveFile(mh.filename)

}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	fmt.Println(m)
	return mh.lines
}

func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	field := mh.dbfresource.FieldValue(row, column)

	value, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(field), simplifiedchinese.GBK.NewDecoder()))

	return ui.TableString(value)
}

func setupUI() {

	box := ui.NewVerticalBox()
	button := ui.NewButton("打开文件")
	box.Append(button, false)

	window := ui.NewWindow("DbfEditer", 1024, 768, true)
	window.SetMargined(true)
	window.SetChild(box)

	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(window)
		if filename != "" {
			mh := newModelHandler(filename)
			model := ui.NewTableModel(mh)
			table := ui.NewTable(&ui.TableParams{
				Model:                         model,
				RowBackgroundColorModelColumn: -1,
			})

			for key, name := range mh.tabletitle {
				table.AppendTextColumn(name.Name, key, ui.TableModelColumnAlwaysEditable, nil)
			}

			window.SetChild(table)
		}
	})

	window.Show()

	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		window.Destroy()
		return true
	})
}

func main() {
	fmt.Println(1)
	ui.Main(setupUI)
}
