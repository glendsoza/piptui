package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Packages struct {
	*tview.Flex
	table *tview.Table
	text  *tview.TextView
}

func NewPackagesView() *Packages {
	text := tview.NewTextView()
	text.SetBorder(true)
	table := tview.NewTable().
		SetSelectable(true, false)
	table.SetBorder(true)
	table.SetBorderColor(tcell.ColorNavy)
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGreen))
	text.SetBorderColor(tcell.ColorYellowGreen)
	flex := tview.NewFlex().
		AddItem(table, 0, 1, true).
		AddItem(text, 0, 2, false)
	return &Packages{
		Flex:  flex,
		table: table,
		text:  text,
	}
}
