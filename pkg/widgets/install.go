package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Install struct {
	*tview.Flex
	form *tview.Form
	text *tview.TextView
}

func NewInstallWidget() *Install {
	form := tview.NewForm().
		AddInputField("Package Name and Version", "", 0, nil, nil).
		AddInputField("Extra Arguments", "", 0, nil, nil).
		AddButton("Submit", nil).
		AddButton("Cancel", nil)
	form.SetBorder(true).
		SetBorderColor(tcell.ColorNavy)
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetMaxLines(100)
	text.SetBorder(true).
		SetBorderColor(tcell.ColorYellowGreen)
	flex := tview.NewFlex().
		AddItem(form, 0, 1, true).
		AddItem(text, 0, 1, false)
	return &Install{
		Flex: flex,
		form: form,
		text: text,
	}
}
