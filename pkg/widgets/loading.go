package widgets

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Loading struct {
	*tview.TextView
	loadChar string
	offset   int
}

func NewLoadingWidget() *Loading {
	text := tview.NewTextView().
		SetDynamicColors(true)
	text.SetBorder(true).
		SetBorderColor(tcell.ColorNavy)
	loadChar := "[yellow::]Loading...."
	return &Loading{
		TextView: text,
		loadChar: loadChar,
		offset:   len(loadChar) - 4,
	}
}

func (l *Loading) Load() bool {
	if !l.HasFocus() {
		return false
	}
	l.SetText(l.loadChar[:l.offset])
	if l.offset == len(l.loadChar) {
		l.offset = len(l.loadChar) - 4
	}
	time.Sleep(500 * time.Millisecond)
	l.offset += 1
	return true
}
