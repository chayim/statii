package ui

import (
	"context"
	"time"

	"github.com/chayim/statii/internal/comms"
	"github.com/chayim/statii/internal/plugins"
	"github.com/gdamore/tcell/v2"
	"github.com/go-co-op/gocron"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()

// processMessages retrieves messages from a redis stream and stores
// the messages on the table
func processMessages(t *tview.Table, cfg *plugins.Config) {
	con := comms.NewConnection(cfg.Database, cfg.Size)

	ctx := context.TODO()
	// read off of stream and toss to table
	messages, err := con.ReadAll(ctx)
	if err != nil {
		return
	}
	for _, m := range messages {
		addRowAtTop(t, m, cfg.Size)
	}
}

// initialize the application container
func RunApp(cfg *plugins.Config) {
	app.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		if evt.Rune() == 113 { // q
			app.Stop()
		}
		return evt
	})

	// text menu
	text := tview.NewTextView()
	text.SetTextColor(tcell.ColorGreenYellow)
	text.SetTextAlign(tview.AlignCenter)
	text.SetText("(q)uit")

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	messages := newTable(cfg.Size)
	layout.SetBorderColor(tcell.ColorBlue)
	layout.AddItem(messages, 0, 1, true)
	layout.AddItem(text, 1, 1, false)

	// TODO handle this via goroutine
	s := gocron.NewScheduler(time.UTC)
	s.Every(30).Seconds().Do(processMessages, messages, cfg) // TODO something about the time for UI refresh
	s.StartAsync()

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
