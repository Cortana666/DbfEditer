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
}

func newModelHandler(filename string) *modelHandler {
	dbf, _ := dbf.LoadFile(filename)

	m := new(modelHandler)
	m.dbfresource = dbf
	m.tabletitle = dbf.Fields()
	m.lines = int(dbf.NumRecords())

	return m
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{}
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {

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
	mainwin := ui.NewWindow("libui Control Gallery", 640, 480, true)

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)
	button := ui.NewButton("打开文件")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		} else {
			mh := newModelHandler(filename)
			model := ui.NewTableModel(mh)

			table := ui.NewTable(&ui.TableParams{
				Model:                         model,
				RowBackgroundColorModelColumn: 3,
			})
			mainwin.SetChild(table)
			mainwin.SetMargined(true)
			for key, name := range mh.tabletitle {
				table.AppendTextColumn(name.Name, key, ui.TableModelColumnAlwaysEditable, nil)
			}

			// tab.Append("Numbers and Lists", table)
			// tab.SetMargined(1, true)

			mainwin.Show()
		}
		entry.SetText(filename)
	})
	grid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry,
		1, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	tab.Append("选择文件", hbox)
	tab.SetMargined(0, true)

	// mh := newModelHandler()
	// model := ui.NewTableModel(mh)

	// table := ui.NewTable(&ui.TableParams{
	// 	Model:                         model,
	// 	RowBackgroundColorModelColumn: 3,
	// })
	// mainwin.SetChild(table)
	// mainwin.SetMargined(true)

	// for key, name := range mh.tabletitle {
	// 	table.AppendTextColumn(name, key, ui.TableModelColumnAlwaysEditable, nil)
	// }

	// tab.Append("Numbers and Lists", table)
	// tab.SetMargined(1, true)

	mainwin.Show()

	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		// mh.dbfresource.Close()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})
}

func main() {
	fmt.Println(1)
	ui.Main(setupUI)
}
