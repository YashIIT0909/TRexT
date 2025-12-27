package components

import (
	"fmt"

	"github.com/YashIIT0909/TRexT/internal/http"
	"github.com/YashIIT0909/TRexT/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// RequestPanel represents the request builder panel
type RequestPanel struct {
	Container    *tview.Flex
	MethodSelect *tview.DropDown
	URLInput     *tview.InputField
	HeadersInput *tview.TextArea
	BodyInput    *tview.TextArea
	SendButton   *tview.Button

	onSend func()
}

// NewRequestPanel creates a new request panel
func NewRequestPanel() *RequestPanel {
	rp := &RequestPanel{}
	rp.build()
	return rp
}

func (rp *RequestPanel) build() {
	// Method dropdown
	rp.MethodSelect = tview.NewDropDown().
		SetLabel("Method: ").
		SetOptions(http.SupportedMethods(), nil).
		SetCurrentOption(0).
		SetFieldWidth(10)
	rp.MethodSelect.SetBorder(false)

	// URL input
	rp.URLInput = tview.NewInputField().
		SetLabel("URL: ").
		SetPlaceholder("https://api.example.com/endpoint").
		SetFieldWidth(0)
	rp.URLInput.SetBorder(false)

	// Top row: method + URL
	topRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(rp.MethodSelect, 20, 0, false).
		AddItem(rp.URLInput, 0, 1, true)
	topRow.SetBorder(true).
		SetTitle(" Request ").
		SetTitleAlign(tview.AlignLeft)

	// Headers input
	rp.HeadersInput = tview.NewTextArea().
		SetPlaceholder("Content-Type: application/json\nAuthorization: Bearer token")
	rp.HeadersInput.SetBorder(true).
		SetTitle(" Headers ").
		SetTitleAlign(tview.AlignLeft)

	// Body input
	rp.BodyInput = tview.NewTextArea().
		SetPlaceholder(`{"key": "value"}`)
	rp.BodyInput.SetBorder(true).
		SetTitle(" Body ").
		SetTitleAlign(tview.AlignLeft)

	// Send button
	rp.SendButton = tview.NewButton("Send Request").
		SetSelectedFunc(func() {
			if rp.onSend != nil {
				rp.onSend()
			}
		})
	rp.SendButton.SetBackgroundColor(tcell.ColorDarkGreen)

	// Button container (centered)
	buttonContainer := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(rp.SendButton, 20, 0, false).
		AddItem(nil, 0, 1, false)

	// Main container
	rp.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topRow, 3, 0, true).
		AddItem(rp.HeadersInput, 0, 1, false).
		AddItem(rp.BodyInput, 0, 2, false).
		AddItem(buttonContainer, 1, 0, false)
}

// SetOnSend sets the callback for when send is triggered
func (rp *RequestPanel) SetOnSend(fn func()) {
	rp.onSend = fn
}

// GetRequest returns the current request from the panel
func (rp *RequestPanel) GetRequest() *http.Request {
	_, method := rp.MethodSelect.GetCurrentOption()
	if method == "" {
		method = "GET"
	}

	return &http.Request{
		Method:  method,
		URL:     rp.URLInput.GetText(),
		Headers: utils.ParseHeaders(rp.HeadersInput.GetText()),
		Body:    rp.BodyInput.GetText(),
	}
}

// SetRequest populates the panel with a request
func (rp *RequestPanel) SetRequest(req *http.Request) {
	// Set method
	for i, m := range http.SupportedMethods() {
		if m == req.Method {
			rp.MethodSelect.SetCurrentOption(i)
			break
		}
	}

	rp.URLInput.SetText(req.URL)
	rp.HeadersInput.SetText(utils.FormatHeaders(req.Headers), true)
	rp.BodyInput.SetText(req.Body, true)
}

// Clear resets the panel
func (rp *RequestPanel) Clear() {
	rp.MethodSelect.SetCurrentOption(0)
	rp.URLInput.SetText("")
	rp.HeadersInput.SetText("", true)
	rp.BodyInput.SetText("", true)
}

// GetFocusableItems returns the list of focusable items in order
func (rp *RequestPanel) GetFocusableItems() []tview.Primitive {
	return []tview.Primitive{
		rp.MethodSelect,
		rp.URLInput,
		rp.HeadersInput,
		rp.BodyInput,
		rp.SendButton,
	}
}

// StatusBadge returns a colored status badge string
func StatusBadge(statusCode int) string {
	color := "white"
	switch {
	case statusCode >= 200 && statusCode < 300:
		color = "green"
	case statusCode >= 300 && statusCode < 400:
		color = "yellow"
	case statusCode >= 400 && statusCode < 500:
		color = "orange"
	case statusCode >= 500:
		color = "red"
	}
	return fmt.Sprintf("[%s]%d[-]", color, statusCode)
}
