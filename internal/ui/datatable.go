package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/chayim/statii/internal/comms"
	"github.com/gdamore/tcell/v2"
	"github.com/lestrrat-go/strftime"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
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

	// anonymous function to handle the click
	anon := func() bool {
		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", m.URL).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", m.URL).Start()
		case "darwin":
			err = exec.Command("open", m.URL).Start()
		default:
			err = fmt.Errorf("unsupported platform")
		}
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		return true
	}

	t.InsertRow(1)
	t.SetCell(1, 0,
		&tview.TableCell{Text: m.Source, Align: tview.AlignLeft, Clicked: anon})
	t.SetCell(1, 1,
		&tview.TableCell{Text: m.ID, Align: tview.AlignCenter, Clicked: anon})
	t.SetCell(1, 2,
		&tview.TableCell{Text: m.Title, Align: tview.AlignCenter, Clicked: anon})
	t.SetCell(1, 3,
		&tview.TableCell{Text: m.Status, Align: tview.AlignCenter, Clicked: anon})
	t.SetCell(1, 4,
		&tview.TableCell{Text: now.FormatString(time.Now()), Align: tview.AlignRight, Clicked: anon})
}

// create a new table
func newTable(rows int64) *tview.Table {
	var messages = tview.NewTable()
	messages.SetTitle(" statii ~ messages ")
	messages.SetBorder(true)
	addHeaders(messages)

	return messages
}
