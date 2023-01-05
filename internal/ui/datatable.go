package ui

import (
	"time"

	"github.com/chayim/statii/internal/comms"
	"github.com/gdamore/tcell/v2"
	"github.com/lestrrat-go/strftime"
	"github.com/rivo/tview"
)

var headers = []string{"Source", "Item", "Title", "Status", "Delivered"}

// addHeaders does a one time header addition to a table
func addHeaders(t *tview.Table) {
	if t.GetRowCount() > 0 {
		return
	}

	t.InsertRow(0)
	colour := tcell.ColorYellow
	for c := 0; c < len(headers); c++ {
		t.SetCell(0, c,
			tview.NewTableCell(headers[c]).
				SetTextColor(colour).
				SetExpansion(5).
				SetAlign(tview.AlignCenter))
	}

}

// handle data insertion, ensuring there is a number of rows met
func addRowAtTop(t *tview.Table, m *comms.Message, rowCount int64) {
	rows := t.GetRowCount()
	if int64(rows) == rowCount {
		t.RemoveRow(int(rowCount - 1))
	}

	// after the header
	now, _ := strftime.New("%Y-%m-%d %H:%M")

	t.InsertRow(1)
	t.SetCell(1, 0,
		&tview.TableCell{Text: m.Source, Align: tview.AlignLeft})
	t.SetCell(1, 1,
		&tview.TableCell{Text: m.ID, Align: tview.AlignCenter})
	t.SetCell(1, 2,
		&tview.TableCell{Text: m.Title, Align: tview.AlignCenter})
	t.SetCell(1, 3,
		&tview.TableCell{Text: m.Status, Align: tview.AlignCenter})
	t.SetCell(1, 4,
		&tview.TableCell{Text: now.FormatString(time.Now()), Align: tview.AlignRight})
}

// create a new table
func newTable(rows int64) *tview.Table {
	var messages = tview.NewTable()
	messages.SetTitle(" statii ~ messages ")
	messages.SetBorder(true)
	addHeaders(messages)

	return messages
}
