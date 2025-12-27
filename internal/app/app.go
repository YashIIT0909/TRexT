package app

import (
	"time"

	"github.com/YashIIT0909/TRexT/internal/components"
	"github.com/YashIIT0909/TRexT/internal/http"
	"github.com/YashIIT0909/TRexT/internal/storage"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App is the main application structure
type App struct {
	// tview application
	tviewApp *tview.Application

	// Layout containers
	rootFlex   *tview.Flex
	pages      *tview.Pages
	mainLayout *tview.Flex

	// Components
	collections  *components.CollectionsList
	requestPanel *components.RequestPanel
	responseView *components.ResponseView
	helpBar      *components.HelpBar
	saveDialog   *components.SaveDialog

	// Services
	httpClient *http.Client
	db         *storage.DB
	config     *storage.Config

	// State
	currentRequest   *http.Request
	currentRequestID int64
	focusIndex       int
	focusables       []tview.Primitive
}

// New creates a new App instance
func New() (*App, error) {
	// Load config
	config, err := storage.LoadConfig()
	if err != nil {
		config = storage.DefaultConfig()
	}

	// Apply theme
	ApplyTheme(GetTheme(config.Theme))

	// Initialize database
	db, err := storage.NewDB()
	if err != nil {
		return nil, err
	}

	app := &App{
		tviewApp:       tview.NewApplication(),
		httpClient:     http.NewClient(),
		db:             db,
		config:         config,
		currentRequest: http.NewRequest(),
	}

	app.buildUI()
	app.setupHandlers()
	app.loadSavedRequests()

	return app, nil
}

// buildUI constructs the user interface
func (a *App) buildUI() {
	// Create components
	a.collections = components.NewCollectionsList()
	a.requestPanel = components.NewRequestPanel()
	a.responseView = components.NewResponseView()
	a.helpBar = components.NewHelpBar()
	a.saveDialog = components.NewSaveDialog()

	// Set initial state
	a.responseView.Clear()

	// Create main layout: [Collections 1] | [Request 2] | [Response 2]
	a.mainLayout = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.collections.Container, 0, 1, true).
		AddItem(a.requestPanel.Container, 0, 2, false).
		AddItem(a.responseView.Container, 0, 2, false)

	// Root layout with help bar at bottom
	a.rootFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.mainLayout, 0, 1, true).
		AddItem(a.helpBar.View, 1, 0, false)

	// Pages for modal dialogs
	a.pages = tview.NewPages().
		AddPage("main", a.rootFlex, true, true).
		AddPage("save", a.saveDialog.Container, true, false)

	// Set focusable items for navigation
	a.focusables = []tview.Primitive{
		a.collections.List,
		a.requestPanel.URLInput,
		a.requestPanel.HeadersInput,
		a.requestPanel.BodyInput,
		a.requestPanel.SendButton,
		a.responseView.BodyView,
	}

	a.tviewApp.SetRoot(a.pages, true)
	a.tviewApp.SetFocus(a.collections.List)
}

// setupHandlers configures event handlers
func (a *App) setupHandlers() {
	// Global key bindings
	a.tviewApp.SetInputCapture(a.handleGlobalKeys)

	// Request panel handlers
	a.requestPanel.SetOnSend(a.executeRequest)

	// Collections list handlers
	a.collections.SetOnSelect(func(req *http.Request) {
		a.currentRequest = req
		a.currentRequestID = req.ID
		a.requestPanel.SetRequest(req)
		a.tviewApp.SetFocus(a.requestPanel.URLInput)
	})

	a.collections.SetOnNew(func() {
		a.newRequest()
	})

	a.collections.SetOnDelete(func(id int64) {
		a.deleteRequest(id)
	})

	// Save dialog handlers
	a.saveDialog.SetOnSave(func(name string) {
		a.saveRequest(name)
		a.pages.HidePage("save")
		a.tviewApp.SetFocus(a.collections.List)
	})

	a.saveDialog.SetOnCancel(func() {
		a.pages.HidePage("save")
		a.tviewApp.SetFocus(a.requestPanel.URLInput)
	})
}

// handleGlobalKeys handles global keyboard shortcuts
func (a *App) handleGlobalKeys(event *tcell.EventKey) *tcell.EventKey {
	// Check if we're in a modal
	if name, _ := a.pages.GetFrontPage(); name != "main" {
		return event
	}

	switch {
	case event.Key() == tcell.KeyCtrlQ:
		// Quit
		a.Stop()
		return nil

	case event.Key() == tcell.KeyCtrlN:
		// New request
		a.newRequest()
		return nil

	case event.Key() == tcell.KeyCtrlS:
		// Save request
		a.showSaveDialog()
		return nil

	case event.Key() == tcell.KeyEnter && event.Modifiers()&tcell.ModCtrl != 0:
		// Send request (Ctrl+Enter)
		a.executeRequest()
		return nil

	case event.Key() == tcell.KeyTab:
		// Navigate to next focusable
		a.focusNext()
		return nil

	case event.Key() == tcell.KeyBacktab:
		// Navigate to previous focusable
		a.focusPrev()
		return nil

	case event.Key() == tcell.KeyCtrlH:
		// Focus collections (left)
		a.tviewApp.SetFocus(a.collections.List)
		a.focusIndex = 0
		return nil

	case event.Key() == tcell.KeyCtrlL:
		// Focus response (right)
		a.tviewApp.SetFocus(a.responseView.BodyView)
		a.focusIndex = len(a.focusables) - 1
		return nil

	case event.Key() == tcell.KeyCtrlU:
		// Focus URL input
		a.tviewApp.SetFocus(a.requestPanel.URLInput)
		a.focusIndex = 1
		return nil
	}

	return event
}

// focusNext focuses the next widget
func (a *App) focusNext() {
	a.focusIndex = (a.focusIndex + 1) % len(a.focusables)
	a.tviewApp.SetFocus(a.focusables[a.focusIndex])
}

// focusPrev focuses the previous widget
func (a *App) focusPrev() {
	a.focusIndex--
	if a.focusIndex < 0 {
		a.focusIndex = len(a.focusables) - 1
	}
	a.tviewApp.SetFocus(a.focusables[a.focusIndex])
}

// executeRequest sends the HTTP request
func (a *App) executeRequest() {
	req := a.requestPanel.GetRequest()
	if req.URL == "" {
		return
	}

	// Update status
	a.responseView.StatusBar.SetText("[yellow]Sending request...[-]")
	a.tviewApp.ForceDraw()

	// Execute in goroutine to not block UI
	go func() {
		resp := a.httpClient.Execute(req)

		// Update UI in main thread
		a.tviewApp.QueueUpdateDraw(func() {
			a.responseView.SetResponse(resp)

			// Add to history
			if a.config.History.Enabled && resp.Error == nil {
				entry := &storage.HistoryEntry{
					URL:        req.URL,
					Method:     req.Method,
					StatusCode: resp.StatusCode,
					Duration:   resp.Duration.Milliseconds(),
					Timestamp:  time.Now().Unix(),
				}
				_ = a.db.AddToHistory(entry)
			}
		})
	}()
}

// newRequest creates a new empty request
func (a *App) newRequest() {
	a.currentRequest = http.NewRequest()
	a.currentRequestID = 0
	a.requestPanel.Clear()
	a.responseView.Clear()
	a.tviewApp.SetFocus(a.requestPanel.URLInput)
}

// showSaveDialog shows the save request dialog
func (a *App) showSaveDialog() {
	req := a.requestPanel.GetRequest()
	if req.Name != "" {
		a.saveDialog.SetName(req.Name)
	} else {
		a.saveDialog.Reset()
	}
	a.pages.ShowPage("save")
	a.tviewApp.SetFocus(a.saveDialog.Modal)
}

// saveRequest saves the current request
func (a *App) saveRequest(name string) {
	req := a.requestPanel.GetRequest()
	req.Name = name
	req.ID = a.currentRequestID

	savedReq := storage.FromHTTPRequest(req, 1) // Default collection

	if err := a.db.SaveRequest(savedReq); err != nil {
		// TODO: Show error
		return
	}

	a.currentRequestID = savedReq.ID
	a.loadSavedRequests()
}

// deleteRequest deletes a saved request
func (a *App) deleteRequest(id int64) {
	if err := a.db.DeleteRequest(id); err != nil {
		return
	}

	if a.currentRequestID == id {
		a.newRequest()
	}

	a.loadSavedRequests()
}

// loadSavedRequests loads all saved requests into the collections list
func (a *App) loadSavedRequests() {
	requests, err := a.db.GetAllRequests()
	if err != nil {
		return
	}
	a.collections.SetRequests(requests)
}

// Run starts the application
func (a *App) Run() error {
	return a.tviewApp.Run()
}

// Stop stops the application
func (a *App) Stop() {
	if a.db != nil {
		a.db.Close()
	}
	a.tviewApp.Stop()
}
