package components

import (
	"fmt"

	"github.com/YashIIT0909/TRexT/internal/http"
	"github.com/YashIIT0909/TRexT/internal/storage"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CollectionsList represents the sidebar with saved requests
type CollectionsList struct {
	Container *tview.Flex
	List      *tview.List
	requests  []*storage.SavedRequest

	onSelect func(req *http.Request)
	onNew    func()
	onDelete func(id int64)
}

// NewCollectionsList creates a new collections list
func NewCollectionsList() *CollectionsList {
	cl := &CollectionsList{
		requests: make([]*storage.SavedRequest, 0),
	}
	cl.build()
	return cl
}

func (cl *CollectionsList) build() {
	cl.List = tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkCyan)

	cl.List.SetBorder(true).
		SetTitle(" Collections ").
		SetTitleAlign(tview.AlignLeft)

	// Add "New Request" item at the top
	cl.List.AddItem("+ New Request", "Create a new request", 'n', func() {
		if cl.onNew != nil {
			cl.onNew()
		}
	})

	// Handle selection
	cl.List.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if index == 0 {
			// "New Request" item
			if cl.onNew != nil {
				cl.onNew()
			}
			return
		}

		// Adjust index for saved requests (subtract 1 for "New Request" item)
		reqIndex := index - 1
		if reqIndex >= 0 && reqIndex < len(cl.requests) {
			if cl.onSelect != nil {
				cl.onSelect(cl.requests[reqIndex].ToHTTPRequest())
			}
		}
	})

	// Handle delete with 'd' key
	cl.List.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'd' || event.Rune() == 'D' {
			index := cl.List.GetCurrentItem()
			if index > 0 { // Don't delete "New Request" item
				reqIndex := index - 1
				if reqIndex >= 0 && reqIndex < len(cl.requests) {
					if cl.onDelete != nil {
						cl.onDelete(cl.requests[reqIndex].ID)
					}
				}
			}
			return nil
		}
		return event
	})

	// Help text at bottom
	helpText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[gray]n: new | d: delete[-]").
		SetTextAlign(tview.AlignCenter)
	helpText.SetBorder(false)

	cl.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(cl.List, 0, 1, true).
		AddItem(helpText, 1, 0, false)
}

// SetRequests populates the list with saved requests
func (cl *CollectionsList) SetRequests(requests []*storage.SavedRequest) {
	cl.requests = requests
	cl.refresh()
}

// refresh updates the list display
func (cl *CollectionsList) refresh() {
	// Remember current selection
	currentIndex := cl.List.GetCurrentItem()

	// Clear and rebuild
	cl.List.Clear()

	// Add "New Request" item
	cl.List.AddItem("+ New Request", "Create a new request", 'n', nil)

	// Add saved requests
	for _, req := range cl.requests {
		methodColor := getMethodColor(req.Method)
		mainText := fmt.Sprintf("[%s]%s[-] %s", methodColor, req.Method, req.Name)
		secondaryText := truncateURL(req.URL, 30)
		cl.List.AddItem(mainText, secondaryText, 0, nil)
	}

	// Restore selection if valid
	if currentIndex < cl.List.GetItemCount() {
		cl.List.SetCurrentItem(currentIndex)
	}
}

// AddRequest adds a request to the list
func (cl *CollectionsList) AddRequest(req *storage.SavedRequest) {
	cl.requests = append(cl.requests, req)
	cl.refresh()
}

// RemoveRequest removes a request from the list by ID
func (cl *CollectionsList) RemoveRequest(id int64) {
	for i, req := range cl.requests {
		if req.ID == id {
			cl.requests = append(cl.requests[:i], cl.requests[i+1:]...)
			break
		}
	}
	cl.refresh()
}

// SetOnSelect sets the callback for when a request is selected
func (cl *CollectionsList) SetOnSelect(fn func(req *http.Request)) {
	cl.onSelect = fn
}

// SetOnNew sets the callback for when new request is triggered
func (cl *CollectionsList) SetOnNew(fn func()) {
	cl.onNew = fn
}

// SetOnDelete sets the callback for when delete is triggered
func (cl *CollectionsList) SetOnDelete(fn func(id int64)) {
	cl.onDelete = fn
}

// getMethodColor returns a color for HTTP method
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "green"
	case "POST":
		return "yellow"
	case "PUT":
		return "blue"
	case "PATCH":
		return "cyan"
	case "DELETE":
		return "red"
	case "HEAD":
		return "purple"
	case "OPTIONS":
		return "gray"
	default:
		return "white"
	}
}

// truncateURL shortens a URL for display
func truncateURL(url string, maxLen int) string {
	if len(url) <= maxLen {
		return url
	}
	return url[:maxLen-3] + "..."
}
