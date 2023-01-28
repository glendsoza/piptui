package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Dependency struct {
	*tview.TreeView
}

func NewDependencyWidget() *Dependency {
	root := tview.NewTreeNode(".")
	tree := tview.NewTreeView().SetRoot(root)
	tree.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)
	return &Dependency{TreeView: tree}
}
