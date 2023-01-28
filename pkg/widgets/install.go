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
		AddInputField("Name and Version", "", 0, nil, nil).
		AddInputField("File Path", "", 0, nil, nil).
		AddButton("Submit", nil)
	form.SetBorder(true).
		SetBorderColor(tcell.ColorNavy)
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetMaxLines(100)
	text.SetBorder(true).
		SetBorderColor(tcell.ColorYellowGreen)
	flex := tview.NewFlex().
		AddItem(form, 0, 1, true).
		AddItem(text, 0, 2, false)
	return &Install{
		Flex: flex,
		form: form,
		text: text,
	}
}
