package widgets

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/glendsoza/piptui/pkg/pip"
	"github.com/rivo/tview"
)

type UI struct {
	app        *tview.Application
	packages   *Packages
	dependency *Dependency
	install    *Install
	loading    *Loading
	pipHelper  *pip.PipHelper
}

func NewUI() *UI {
	ui := &UI{
		app:        tview.NewApplication(),
		packages:   NewPackagesView(),
		dependency: NewDependencyWidget(),
		install:    NewInstallWidget(),
		loading:    NewLoadingWidget(),
		pipHelper:  pip.NewPipHelper(),
	}
	ui.preparePackagesWidget()
	ui.prepareInstallWidget()
	ui.prepareDependencyWidget()
	return ui
}

func (u *UI) preparePackagesWidget() {
	u.packages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlT:
			u.app.SetRoot(u.dependency, true)
			return nil
		case tcell.KeyCtrlD:
			if u.packages.table.HasFocus() {
				u.app.SetFocus(u.packages.text)
			} else {
				u.app.SetFocus(u.packages.table)
			}
			return nil
		case tcell.KeyCtrlI:
			u.app.SetRoot(u.install, true)
			return nil
		case tcell.KeyCtrlU:
			// column is always 0 as we are only interested in name and not version
			r, _ := u.packages.table.GetSelection()
			if r != 0 {
				selected := u.packages.table.GetCell(r, 0)
				modal := tview.NewModal()
				modal.SetText(fmt.Sprintf("Do you really want to uninstall %s[-:-:-]?\n (this is a force uninstall)", selected.Text))
				modal.AddButtons([]string{"Confirm", "Cancel"})
				modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					switch buttonLabel {
					case "Cancel":
						u.app.SetRoot(u.packages, true)
					case "Confirm":
						out := u.pipHelper.Uninstall(selected.Text)
						modal.SetText(string(out))
						u.loadPackagesWidget()
						u.app.SetRoot(modal, true)
					}
				})
				u.app.SetRoot(modal, true)
				return nil
			}
		}
		return event
	})
	u.packages.table.SetSelectionChangedFunc(func(row, column int) {
		if row != 0 {
			name := u.packages.table.GetCell(row, 0)
			version := u.packages.table.GetCell(row, 1)
			info := ""
			pkg, err := u.pipHelper.Show(name.Text, version.Text)
			if err != nil {
				info = err.Error()
			} else {
				info = pkg.ToString()
			}
			u.packages.text.SetText(info).ScrollToBeginning()
		} else {
			u.packages.text.SetText("You will see package details here").ScrollToBeginning()
		}

	})
}

func (u *UI) prepareInstallWidget() {
	Submitbutton := u.install.form.GetButton(0)
	Submitbutton.SetSelectedFunc(func() {
		nameAndVersion := u.install.form.GetFormItemByLabel("Package Name and Version")
		extraArguments := u.install.form.GetFormItemByLabel("Extra Arguments")
		nameAndVersionField, fieldOk := nameAndVersion.(*tview.InputField)
		extraArgumentsField, extraArgsOk := extraArguments.(*tview.InputField)
		if fieldOk && extraArgsOk {
			command := nameAndVersionField.GetText()
			if strings.TrimSpace(extraArgumentsField.GetText()) != "" {
				command = fmt.Sprintf("%s %s", strings.TrimSpace(extraArgumentsField.GetText()), command)
			}
			r, _ := u.pipHelper.InstallAndStream(command)
			// set the label as disabled
			Submitbutton.SetLabel("Disabled")
			go func() {
				reader := bufio.NewScanner(r)
				counter := 0
				for reader.Scan() {
					fmt.Fprintf(u.install.text, "[green]>>>[-:-:-] %s\n", reader.Text())
					counter += 1
					if counter >= 3 {
						counter = 0
						u.app.QueueUpdateDraw(
							func() {
								if u.install.HasFocus() {
									u.app.SetRoot(u.install, true)
								}
							},
						)
					}
				}
				u.loadPackagesWidget()
				Submitbutton.SetLabel("Submit")
				u.app.QueueUpdateDraw(func() {
					if u.packages.HasFocus() {
						u.app.SetRoot(u.packages, true)
					} else if u.install.HasFocus() {
						u.app.SetRoot(u.install, true)
					}
				})
			}()
		}
	})
	cancelButton := u.install.form.GetButton(1)
	cancelButton.SetSelectedFunc(func() {
		if u.pipHelper.InstallInProgress() {
			u.pipHelper.CancelInstall()
			fmt.Fprintln(u.install.text, "[red::]Cancelled[-:-:-]")
		} else {
			fmt.Fprintln(u.install.text, "[blue::]Nothing to cancel[-:-:-]")
		}
	})
	u.install.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			u.app.SetRoot(u.packages, true)
			return nil
		case tcell.KeyCtrlD:
			if u.install.form.HasFocus() {
				u.app.SetFocus(u.install.text)
			} else {
				u.app.SetFocus(u.install.form)
			}
			return nil
		}
		return event
	})
}

func (u *UI) prepareDependencyWidget() {
	u.dependency.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			u.app.SetRoot(u.packages, true)
			return nil
		}
		return event
	})
}

func (u *UI) loadDependencyWidget(packages pip.Packages) {
	packagesMap := map[string]*tview.TreeNode{}
	u.dependency.TreeView.GetRoot().ClearChildren()
	for _, pkg := range packages {
		packagesMap[pkg.Name] = tview.NewTreeNode(pkg.Name)
	}
	for _, pp := range packages {
		if len(pp.RequiredBy) == 0 {
			u.dependency.GetRoot().AddChild(packagesMap[pp.Name].SetColor(tcell.ColorGreen))
		} else {
			for _, dd := range pp.RequiredBy {
				packagesMap[dd].AddChild(packagesMap[pp.Name].SetColor(tcell.ColorYellow))
			}
		}
	}
}

func (u *UI) loadPackagesWidget() error {
	installedPackages, err := u.pipHelper.List()
	if err != nil {
		return err
	}
	c := 1
	u.packages.table.SetCell(0, 0, tview.NewTableCell("Name"))
	u.packages.table.SetCell(0, 1, tview.NewTableCell("Version"))
	u.packages.text.SetText("You will see package details here")
	for _, pkg := range installedPackages {
		u.packages.table.SetCell(c, 0, tview.NewTableCell(pkg.Name))
		u.packages.table.SetCell(c, 1, tview.NewTableCell(pkg.Version))
		c += 1
	}
	if len(installedPackages) > 0 {
		u.packages.text.SetTitle(installedPackages[0].Location)
	}
	u.loadDependencyWidget(installedPackages)
	return nil
}

func (u *UI) Run() error {
	fatalErrChan := make(chan error)
	go func() {
		for {
			if !u.loading.Load() {
				break
			}
			u.app.QueueUpdateDraw(func() {
				if u.loading.HasFocus() {
					u.app.SetRoot(u.loading, true)
				}
			})
		}
	}()
	go func() {
		err := u.loadPackagesWidget()
		if err != nil {
			fatalErrChan <- err
		}
		u.app.QueueUpdateDraw(func() {
			u.app.SetRoot(u.packages, true)
		})
	}()
	go func() {
		fatalErrChan <- u.app.SetRoot(u.loading, true).Run()
	}()
	err := <-fatalErrChan
	u.app.Stop()
	return err
}
