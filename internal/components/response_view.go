package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/YashIIT0909/TRexT/internal/http"
	"github.com/YashIIT0909/TRexT/internal/utils"
	"github.com/rivo/tview"
)

// ResponseView represents the response display panel
type ResponseView struct {
	Container   *tview.Flex
	StatusBar   *tview.TextView
	HeadersView *tview.TextView
	BodyView    *tview.TextView
	tabs        *tview.TextView
	currentTab  string
	response    *http.Response
}

// NewResponseView creates a new response view
func NewResponseView() *ResponseView {
	rv := &ResponseView{
		currentTab: "body",
	}
	rv.build()
	return rv
}

func (rv *ResponseView) build() {
	// Status bar
	rv.StatusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	rv.StatusBar.SetBorder(false)

	// Tabs
	rv.tabs = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[\"body\"][darkcyan]Body[white][\"\"] | [\"headers\"]Headers[\"\"]")
	rv.tabs.SetBorder(false)

	// Header row
	headerRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(rv.StatusBar, 0, 1, false).
		AddItem(rv.tabs, 30, 0, false)

	// Headers view
	rv.HeadersView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	rv.HeadersView.SetBorder(false)

	// Body view
	rv.BodyView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	rv.BodyView.SetBorder(false)

	// Content container (shows either body or headers)
	contentBox := tview.NewFlex().
		AddItem(rv.BodyView, 0, 1, false)
	contentBox.SetBorder(true).
		SetTitle(" Response ").
		SetTitleAlign(tview.AlignLeft)

	// Main container
	rv.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(headerRow, 1, 0, false).
		AddItem(contentBox, 0, 1, false)
}

// SetResponse displays the response
func (rv *ResponseView) SetResponse(resp *http.Response) {
	rv.response = resp

	if resp.Error != nil {
		rv.StatusBar.SetText(fmt.Sprintf("[red]Error:[-] %s", resp.Error.Error()))
		rv.BodyView.SetText("")
		rv.HeadersView.SetText("")
		return
	}

	// Status bar
	statusColor := "green"
	switch {
	case resp.StatusCode >= 300 && resp.StatusCode < 400:
		statusColor = "yellow"
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		statusColor = "orange"
	case resp.StatusCode >= 500:
		statusColor = "red"
	}

	statusText := fmt.Sprintf("[%s]%s[-] | %dms | %s",
		statusColor,
		resp.Status,
		resp.Duration.Milliseconds(),
		formatSize(resp.Size),
	)
	rv.StatusBar.SetText(statusText)

	// Format body
	body := resp.BodyString()
	if utils.IsValidJSON(body) {
		if formatted, err := utils.FormatJSON(body); err == nil {
			body = formatted
		}
	}
	rv.BodyView.SetText(body)

	// Format headers
	var headerLines []string
	// Sort headers for consistent display
	keys := make([]string, 0, len(resp.Headers))
	for k := range resp.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		values := resp.Headers[key]
		for _, v := range values {
			headerLines = append(headerLines, fmt.Sprintf("[darkcyan]%s:[-] %s", key, v))
		}
	}
	rv.HeadersView.SetText(strings.Join(headerLines, "\n"))
}

// Clear resets the response view
func (rv *ResponseView) Clear() {
	rv.response = nil
	rv.StatusBar.SetText("[gray]No response yet[-]")
	rv.BodyView.SetText("")
	rv.HeadersView.SetText("")
}

// ShowTab switches between body and headers view
func (rv *ResponseView) ShowTab(tab string) {
	rv.currentTab = tab
	// Update tab highlighting
	switch tab {
	case "body":
		rv.tabs.SetText("[darkcyan]Body[-] | [gray]Headers[-]")
	case "headers":
		rv.tabs.SetText("[gray]Body[-] | [darkcyan]Headers[-]")
	}
}

// ToggleTab switches between tabs
func (rv *ResponseView) ToggleTab() {
	if rv.currentTab == "body" {
		rv.ShowTab("headers")
	} else {
		rv.ShowTab("body")
	}
}

// GetCurrentTab returns the current tab
func (rv *ResponseView) GetCurrentTab() string {
	return rv.currentTab
}

// formatSize formats bytes to human readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
