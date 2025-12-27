package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SaveDialog represents a modal dialog for saving requests
type SaveDialog struct {
	Modal     *tview.Form
	Container *tview.Flex
	nameInput *tview.InputField

	onSave   func(name string)
	onCancel func()
}

// NewSaveDialog creates a new save dialog
func NewSaveDialog() *SaveDialog {
	sd := &SaveDialog{}
	sd.build()
	return sd
}

func (sd *SaveDialog) build() {
	sd.Modal = tview.NewForm()
	sd.Modal.SetBorder(true).
		SetTitle(" Save Request ").
		SetTitleAlign(tview.AlignCenter)

	sd.Modal.AddInputField("Name:", "", 40, nil, nil)
	sd.Modal.AddButton("Save", func() {
		if sd.onSave != nil {
			name := sd.Modal.GetFormItem(0).(*tview.InputField).GetText()
			sd.onSave(name)
		}
	})
	sd.Modal.AddButton("Cancel", func() {
		if sd.onCancel != nil {
			sd.onCancel()
		}
	})

	// Handle escape key
	sd.Modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if sd.onCancel != nil {
				sd.onCancel()
			}
			return nil
		}
		return event
	})

	// Center the modal
	sd.Container = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(sd.Modal, 10, 0, true).
			AddItem(nil, 0, 1, false), 50, 0, true).
		AddItem(nil, 0, 1, false)
}

// SetOnSave sets the save callback
func (sd *SaveDialog) SetOnSave(fn func(name string)) {
	sd.onSave = fn
}

// SetOnCancel sets the cancel callback
func (sd *SaveDialog) SetOnCancel(fn func()) {
	sd.onCancel = fn
}

// Reset clears the dialog
func (sd *SaveDialog) Reset() {
	sd.Modal.GetFormItem(0).(*tview.InputField).SetText("")
}

// SetName sets the name input value
func (sd *SaveDialog) SetName(name string) {
	sd.Modal.GetFormItem(0).(*tview.InputField).SetText(name)
}

// HelpBar represents the bottom help bar
type HelpBar struct {
	View *tview.TextView
}

// NewHelpBar creates a new help bar
func NewHelpBar() *HelpBar {
	hb := &HelpBar{}
	hb.View = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	hb.View.SetBorder(false)
	hb.SetDefaultHelp()
	return hb
}

// SetDefaultHelp sets the default help text
func (hb *HelpBar) SetDefaultHelp() {
	hb.View.SetText("[yellow]Ctrl+Enter[-]: Send | [yellow]Ctrl+S[-]: Save | [yellow]Ctrl+N[-]: New | [yellow]Tab[-]: Navigate | [yellow]Ctrl+Q[-]: Quit")
}

// SetText sets custom help text
func (hb *HelpBar) SetText(text string) {
	hb.View.SetText(text)
}

// ConfirmDialog represents a confirmation dialog
type ConfirmDialog struct {
	Modal     *tview.Modal
	Container *tview.Flex

	onConfirm func()
	onCancel  func()
}

// NewConfirmDialog creates a new confirmation dialog
func NewConfirmDialog(message string) *ConfirmDialog {
	cd := &ConfirmDialog{}
	cd.build(message)
	return cd
}

func (cd *ConfirmDialog) build(message string) {
	cd.Modal = tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" && cd.onConfirm != nil {
				cd.onConfirm()
			} else if cd.onCancel != nil {
				cd.onCancel()
			}
		})

	cd.Container = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(cd.Modal, 7, 0, true).
			AddItem(nil, 0, 1, false), 50, 0, true).
		AddItem(nil, 0, 1, false)
}

// SetOnConfirm sets the confirm callback
func (cd *ConfirmDialog) SetOnConfirm(fn func()) {
	cd.onConfirm = fn
}

// SetOnCancel sets the cancel callback
func (cd *ConfirmDialog) SetOnCancel(fn func()) {
	cd.onCancel = fn
}
