package application

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v3/internal/capabilities"

	"github.com/wailsapp/wails/v3/pkg/icons"

	"github.com/samber/lo"

	"github.com/wailsapp/wails/v3/pkg/assetserver"
	"github.com/wailsapp/wails/v3/pkg/assetserver/webview"

	wailsruntime "github.com/wailsapp/wails/v3/internal/runtime"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var globalApplication *App

func init() {
	runtime.LockOSThread()
}

type EventListener struct {
	callback func(*Event)
}

func Get() *App {
	return globalApplication
}

func New(appOptions Options) *App {
	if globalApplication != nil {
		return globalApplication
	}

	mergeApplicationDefaults(&appOptions)

	result := newApplication(&appOptions)
	globalApplication = result

	if result.isDebugMode && result.Logger == nil {
		result.Logger = DefaultLogger()
	}

	result.Events = NewWailsEventProcessor(result.dispatchEventToWindows)

	opts := &assetserver.Options{
		Assets:     appOptions.Assets.FS,
		Handler:    appOptions.Assets.Handler,
		Middleware: assetserver.Middleware(appOptions.Assets.Middleware),
	}

	srv, err := assetserver.NewAssetServer(opts, false, result.Logger, wailsruntime.RuntimeAssetsBundle, result.isDebugMode)
	if err != nil {
		result.fatal(err.Error())
	}

	// Pass through the capabilities
	srv.GetCapabilities = func() []byte {
		return globalApplication.capabilities.AsBytes()
	}

	srv.GetFlags = func() []byte {
		updatedOptions := result.impl.GetFlags(appOptions)
		flags, err := json.Marshal(updatedOptions)
		if err != nil {
			log.Fatal("Invalid flags provided to application: ", err.Error())
		}
		return flags
	}

	srv.UseRuntimeHandler(NewMessageProcessor())
	result.assets = srv

	result.bindings, err = NewBindings(appOptions.Bind)
	if err != nil {
		println("Fatal error in application initialisation: ", err.Error())
		os.Exit(1)
	}

	result.plugins = NewPluginManager(appOptions.Plugins, srv)
	err = result.plugins.Init()
	if err != nil {
		result.Quit()
		os.Exit(1)
	}

	err = result.bindings.AddPlugins(appOptions.Plugins)
	if err != nil {
		println("Fatal error in application initialisation: ", err.Error())
		os.Exit(1)
	}

	return result
}

func mergeApplicationDefaults(o *Options) {
	if o.Name == "" {
		o.Name = "My Wails Application"
	}
	if o.Description == "" {
		o.Description = "An application written using Wails"
	}
	if o.Icon == nil {
		o.Icon = icons.ApplicationLightMode256
	}

}

type (
	platformApp interface {
		run() error
		destroy()
		setApplicationMenu(menu *Menu)
		name() string
		getCurrentWindowID() uint
		showAboutDialog(name string, description string, icon []byte)
		setIcon(icon []byte)
		on(id uint)
		dispatchOnMainThread(id uint)
		hide()
		show()
		getPrimaryScreen() (*Screen, error)
		getScreens() ([]*Screen, error)
		GetFlags(options Options) map[string]any
		isOnMainThread() bool
	}

	runnable interface {
		run()
	}
)

func processPanicHandlerRecover() {
	h := globalApplication.options.PanicHandler
	if h == nil {
		return
	}

	if err := recover(); err != nil {
		h(err)
	}
}

// Messages sent from javascript get routed here
type windowMessage struct {
	windowId uint
	message  string
}

var windowMessageBuffer = make(chan *windowMessage)

type dragAndDropMessage struct {
	windowId  uint
	filenames []string
}

var windowDragAndDropBuffer = make(chan *dragAndDropMessage)

var _ webview.Request = &webViewAssetRequest{}

const webViewRequestHeaderWindowId = "x-wails-window-id"
const webViewRequestHeaderWindowName = "x-wails-window-name"

type webViewAssetRequest struct {
	webview.Request
	windowId   uint
	windowName string
}

func (r *webViewAssetRequest) Header() (http.Header, error) {
	h, err := r.Request.Header()
	if err != nil {
		return nil, err
	}

	hh := h.Clone()
	hh.Set(webViewRequestHeaderWindowId, strconv.FormatUint(uint64(r.windowId), 10))
	return hh, nil
}

var webviewRequests = make(chan *webViewAssetRequest)

type eventHook struct {
	callback func(*Event)
}

type App struct {
	options                       Options
	applicationEventListeners     map[uint][]*EventListener
	applicationEventListenersLock sync.RWMutex
	applicationEventHooks         map[uint][]*eventHook
	applicationEventHooksLock     sync.RWMutex

	// Windows
	windows     map[uint]*WebviewWindow
	windowsLock sync.Mutex

	// System Trays
	systemTrays      map[uint]*SystemTray
	systemTraysLock  sync.Mutex
	systemTrayID     uint
	systemTrayIDLock sync.RWMutex

	// MenuItems
	menuItems     map[uint]*MenuItem
	menuItemsLock sync.Mutex

	// Running
	running    bool
	runLock    sync.Mutex
	pendingRun []runnable

	bindings *Bindings
	plugins  *PluginManager

	// platform app
	impl platformApp

	// The main application menu
	ApplicationMenu *Menu

	clipboard *Clipboard
	Events    *EventProcessor
	Logger    *slog.Logger

	contextMenus     map[string]*Menu
	contextMenusLock sync.Mutex

	assets   *assetserver.AssetServer
	startURL string

	// Hooks
	windowCreatedCallbacks []func(window *WebviewWindow)
	pid                    int

	// Capabilities
	capabilities capabilities.Capabilities
	isDebugMode  bool
}

func (a *App) init() {
	a.applicationEventListeners = make(map[uint][]*EventListener)
	a.windows = make(map[uint]*WebviewWindow)
	a.systemTrays = make(map[uint]*SystemTray)
	a.contextMenus = make(map[string]*Menu)
	a.Logger = a.options.Logger
	a.pid = os.Getpid()
}

func (a *App) getSystemTrayID() uint {
	a.systemTrayIDLock.Lock()
	defer a.systemTrayIDLock.Unlock()
	a.systemTrayID++
	return a.systemTrayID
}

func (a *App) getWindowForID(id uint) *WebviewWindow {
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	return a.windows[id]
}

func (a *App) deleteWindowByID(id uint) {
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	delete(a.windows, id)
}

func (a *App) Capabilities() capabilities.Capabilities {
	return a.capabilities
}

func (a *App) On(eventType events.ApplicationEventType, callback func(event *Event)) func() {
	eventID := uint(eventType)
	a.applicationEventListenersLock.Lock()
	defer a.applicationEventListenersLock.Unlock()
	listener := &EventListener{
		callback: callback,
	}
	a.applicationEventListeners[eventID] = append(a.applicationEventListeners[eventID], listener)
	if a.impl != nil {
		go a.impl.on(eventID)
	}

	return func() {
		// lock the map
		a.applicationEventListenersLock.Lock()
		defer a.applicationEventListenersLock.Unlock()
		// Remove listener
		a.applicationEventListeners[eventID] = lo.Without(a.applicationEventListeners[eventID], listener)
	}
}

// RegisterHook registers a hook for the given event type. Hooks are called before the event listeners and can cancel the event.
// The returned function can be called to remove the hook.
func (a *App) RegisterHook(eventType events.ApplicationEventType, callback func(event *Event)) func() {
	eventID := uint(eventType)
	a.applicationEventHooksLock.Lock()
	defer a.applicationEventHooksLock.Unlock()
	thisHook := &eventHook{
		callback: callback,
	}
	a.applicationEventHooks[eventID] = append(a.applicationEventHooks[eventID], thisHook)

	return func() {
		a.applicationEventHooksLock.Lock()
		a.applicationEventHooks[eventID] = lo.Without(a.applicationEventHooks[eventID], thisHook)
		a.applicationEventHooksLock.Unlock()
	}
}

func (a *App) NewWebviewWindow() *WebviewWindow {
	return a.NewWebviewWindowWithOptions(WebviewWindowOptions{})
}

func (a *App) GetPID() int {
	return a.pid
}

func (a *App) info(message string, args ...any) {
	if a.Logger != nil {
		a.Logger.Info(message, args...)
	}
}

func (a *App) debug(message string, args ...any) {
	if a.Logger != nil {
		a.Logger.Debug(message, args...)
	}
}

func (a *App) fatal(message string, args ...any) {
	msg := "A FATAL ERROR HAS OCCURRED: " + message
	if a.Logger != nil {
		a.Logger.Error(msg, args...)
	} else {
		println(msg)
	}
	os.Exit(1)
}

func (a *App) error(message string, args ...any) {
	if a.Logger != nil {
		a.Logger.Error(message, args...)
	}
}

func (a *App) NewWebviewWindowWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow {
	newWindow := NewWindow(windowOptions)
	id := newWindow.id

	a.windowsLock.Lock()
	a.windows[id] = newWindow
	a.windowsLock.Unlock()

	// Call hooks
	for _, hook := range a.windowCreatedCallbacks {
		hook(newWindow)
	}

	a.runOrDeferToAppRun(newWindow)

	return newWindow
}

func (a *App) NewSystemTray() *SystemTray {
	id := a.getSystemTrayID()
	newSystemTray := NewSystemTray(id)

	a.systemTraysLock.Lock()
	a.systemTrays[id] = newSystemTray
	a.systemTraysLock.Unlock()

	a.runOrDeferToAppRun(newSystemTray)

	return newSystemTray
}

func (a *App) Run() error {
	a.logStartup()
	a.logPlatformInfo()

	// Setup panic handler
	defer processPanicHandlerRecover()

	a.impl = newPlatformApp(a)
	go func() {
		for {
			event := <-applicationEvents
			a.handleApplicationEvent(event)
		}
	}()
	go func() {
		for {
			event := <-windowEvents
			a.handleWindowEvent(event)
		}
	}()
	go func() {
		for {
			request := <-webviewRequests
			a.handleWebViewRequest(request)
		}
	}()
	go func() {
		for {
			event := <-windowMessageBuffer
			a.handleWindowMessage(event)
		}
	}()
	go func() {
		for {
			dragAndDropMessage := <-windowDragAndDropBuffer
			a.handleDragAndDropMessage(dragAndDropMessage)
		}
	}()

	go func() {
		for {
			menuItemID := <-menuItemClicked
			a.handleMenuItemClicked(menuItemID)
		}
	}()

	a.runLock.Lock()
	a.running = true

	for _, systray := range a.pendingRun {
		go systray.run()
	}
	a.pendingRun = nil

	a.runLock.Unlock()

	// set the application menu
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		a.impl.setApplicationMenu(a.ApplicationMenu)
	}
	a.impl.setIcon(a.options.Icon)

	err := a.impl.run()
	if err != nil {
		return err
	}

	a.plugins.Shutdown()

	return nil
}

func (a *App) handleApplicationEvent(event uint) {
	a.applicationEventListenersLock.RLock()
	listeners, ok := a.applicationEventListeners[event]
	a.applicationEventListenersLock.RUnlock()
	if !ok {
		return
	}

	thisEvent := &Event{}

	// Process Hooks
	a.applicationEventHooksLock.RLock()
	hooks, ok := a.applicationEventHooks[event]
	a.applicationEventHooksLock.RUnlock()
	if ok {
		for _, thisHook := range hooks {
			thisHook.callback(thisEvent)
			if thisEvent.Cancelled {
				return
			}
		}
	}

	for _, listener := range listeners {
		go listener.callback(thisEvent)
	}
}

func (a *App) handleDragAndDropMessage(event *dragAndDropMessage) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.handleDragAndDropMessage(event)
}

func (a *App) handleWindowMessage(event *windowMessage) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.handleMessage(event.message)
}

func (a *App) handleWebViewRequest(request *webViewAssetRequest) {
	// Get window from window map
	url, _ := request.URL()
	a.info("Window: '%s', Request: %s", request.windowName, url)
	a.assets.ServeWebViewRequest(request)
}

func (a *App) handleWindowEvent(event *windowEvent) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.WindowID]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.WindowID)
		return
	}
	window.handleWindowEvent(event.EventID)
}

func (a *App) handleMenuItemClicked(menuItemID uint) {
	menuItem := getMenuItemByID(menuItemID)
	if menuItem == nil {
		log.Printf("MenuItem #%d not found", menuItemID)
		return
	}
	menuItem.handleClick()
}

func (a *App) CurrentWindow() *WebviewWindow {
	if a.impl == nil {
		return nil
	}
	id := a.impl.getCurrentWindowID()
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	return a.windows[id]
}

func (a *App) Quit() {
	invokeSync(func() {
		a.windowsLock.Lock()
		for _, window := range a.windows {
			window.Destroy()
		}
		a.windowsLock.Unlock()
		a.systemTraysLock.Lock()
		for _, systray := range a.systemTrays {
			systray.Destroy()
		}
		a.systemTraysLock.Unlock()
		if a.impl != nil {
			a.impl.destroy()
		}
	})
}

func (a *App) SetMenu(menu *Menu) {
	a.ApplicationMenu = menu
	if a.impl != nil {
		a.impl.setApplicationMenu(menu)
	}
}
func (a *App) ShowAboutDialog() {
	if a.impl != nil {
		a.impl.showAboutDialog(a.options.Name, a.options.Description, a.options.Icon)
	}
}

func InfoDialog() *MessageDialog {
	return newMessageDialog(InfoDialogType)
}

func QuestionDialog() *MessageDialog {
	return newMessageDialog(QuestionDialogType)
}

func WarningDialog() *MessageDialog {
	return newMessageDialog(WarningDialogType)
}

func ErrorDialog() *MessageDialog {
	return newMessageDialog(ErrorDialogType)
}

func OpenDirectoryDialog() *MessageDialog {
	return newMessageDialog(OpenDirectoryDialogType)
}

func OpenFileDialog() *OpenFileDialogStruct {
	return newOpenFileDialog()
}

func SaveFileDialog() *SaveFileDialogStruct {
	return newSaveFileDialog()
}

func (a *App) GetPrimaryScreen() (*Screen, error) {
	return a.impl.getPrimaryScreen()
}

func (a *App) GetScreens() ([]*Screen, error) {
	return a.impl.getScreens()
}

func (a *App) Clipboard() *Clipboard {
	if a.clipboard == nil {
		a.clipboard = newClipboard()
	}
	return a.clipboard
}

func (a *App) dispatchOnMainThread(fn func()) {
	// If we are on the main thread, just call the function
	if a.impl.isOnMainThread() {
		fn()
		return
	}

	mainThreadFunctionStoreLock.Lock()
	id := generateFunctionStoreID()
	mainThreadFunctionStore[id] = fn
	mainThreadFunctionStoreLock.Unlock()
	// Call platform specific dispatch function
	a.impl.dispatchOnMainThread(id)
}

func OpenFileDialogWithOptions(options *OpenFileDialogOptions) *OpenFileDialogStruct {
	result := OpenFileDialog()
	result.SetOptions(options)
	return result
}

func SaveFileDialogWithOptions(s *SaveFileDialogOptions) *SaveFileDialogStruct {
	result := SaveFileDialog()
	result.SetOptions(s)
	return result
}

func (a *App) dispatchEventToWindows(event *WailsEvent) {
	for _, window := range a.windows {
		window.dispatchWailsEvent(event)
	}
}

func (a *App) Hide() {
	if a.impl != nil {
		a.impl.hide()
	}
}

func (a *App) Show() {
	if a.impl != nil {
		a.impl.show()
	}
}

func (a *App) RegisterContextMenu(name string, menu *Menu) {
	a.contextMenusLock.Lock()
	defer a.contextMenusLock.Unlock()
	a.contextMenus[name] = menu
}

func (a *App) getContextMenu(name string) (*Menu, bool) {
	a.contextMenusLock.Lock()
	defer a.contextMenusLock.Unlock()
	menu, ok := a.contextMenus[name]
	return menu, ok

}

func (a *App) OnWindowCreation(callback func(window *WebviewWindow)) {
	a.windowCreatedCallbacks = append(a.windowCreatedCallbacks, callback)
}

func (a *App) GetWindowByName(name string) *WebviewWindow {
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	for _, window := range a.windows {
		if window.Name() == name {
			return window
		}
	}
	return nil
}

func (a *App) runOrDeferToAppRun(r runnable) {
	a.runLock.Lock()
	running := a.running
	if !running {
		a.pendingRun = append(a.pendingRun, r)
	}
	a.runLock.Unlock()

	if running {
		r.run()
	}
}

func invokeSync(fn func()) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer processPanicHandlerRecover()
		fn()
		wg.Done()
	})
	wg.Wait()
}

func invokeSyncWithResult[T any](fn func() T) (res T) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer processPanicHandlerRecover()
		res = fn()
		wg.Done()
	})
	wg.Wait()
	return res
}

func invokeSyncWithError(fn func() error) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer processPanicHandlerRecover()
		err = fn()
		wg.Done()
	})
	wg.Wait()
	return
}

func invokeSyncWithResultAndError[T any](fn func() (T, error)) (res T, err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer processPanicHandlerRecover()
		res, err = fn()
		wg.Done()
	})
	wg.Wait()
	return res, err
}
