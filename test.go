package main

import (
	"github.com/SebastiaanKlippert/go-foxpro-dbf"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"

	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

type modelHandler struct {
	dbfresource *dbf.DBF
	tabletitle  []string
	lines       int
}

func newModelHandler() *modelHandler {
	dbf, _ := dbf.OpenFile("/Users/yangjian/Desktop/单考数据-8人.dbf", new(dbf.UTF8Decoder))

	m := new(modelHandler)
	m.dbfresource = dbf
	m.tabletitle = dbf.FieldNames()
	m.lines = int(dbf.NumRecords())
	return m
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{}
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {

}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return mh.lines
}

func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	mh.dbfresource.GoTo(uint32(row))
	deleted, _ := mh.dbfresource.Deleted()
	if deleted {
		return ui.TableString("Deleted")
	}
	field, _ := mh.dbfresource.Field(column)
	if field == nil {
		field = ""
	}

	I := strings.NewReader(field.(string))
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, _ := ioutil.ReadAll(O)
	return ui.TableString(d)
}

func setupUI() {
	mainwin := ui.NewWindow("libui Control Gallery", 640, 480, true)

	mh := newModelHandler()
	model := ui.NewTableModel(mh)

	table := ui.NewTable(&ui.TableParams{
		Model:                         model,
		RowBackgroundColorModelColumn: 3,
	})
	mainwin.SetChild(table)
	mainwin.SetMargined(true)

	for key, name := range mh.tabletitle {
		table.AppendTextColumn(name, key, ui.TableModelColumnAlwaysEditable, nil)
	}

	// tab := ui.NewTab()
	// mainwin.SetChild(tab)
	// mainwin.SetMargined(true)
	// hbox := ui.NewHorizontalBox()
	// hbox.SetPadded(true)
	// vbox := ui.NewVerticalBox()
	// vbox.SetPadded(true)
	// hbox.Append(vbox, false)
	// grid := ui.NewGrid()
	// grid.SetPadded(true)
	// vbox.Append(grid, false)
	// button := ui.NewButton("Open File")
	// entry := ui.NewEntry()
	// entry.SetReadOnly(true)
	// button.OnClicked(func(*ui.Button) {
	// 	filename := ui.OpenFile(mainwin)
	// 	if filename == "" {
	// 		filename = "(cancelled)"
	// 	}
	// 	entry.SetText(filename)
	// })
	// grid.Append(button,
	// 	0, 0, 1, 1,
	// 	false, ui.AlignFill, false, ui.AlignFill)
	// grid.Append(entry,
	// 	1, 0, 1, 1,
	// 	true, ui.AlignFill, false, ui.AlignFill)

	// tab.Append("Data Choosers", hbox)
	// tab.SetMargined(0, true)

	// tab.Append("Numbers and Lists", table)
	// tab.SetMargined(1, true)

	mainwin.Show()

	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		mh.dbfresource.Close()
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
