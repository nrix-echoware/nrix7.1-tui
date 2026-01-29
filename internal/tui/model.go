package tui

import (
	"time"
	"terminal-echoware/internal/api"
	"terminal-echoware/pkg/config"
	"terminal-echoware/pkg/types"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time
type notificationClearMsg struct{}

type Notification struct {
	Message string
	Type    string // "success", "error", "info"
}

// VariantSelection holds the selected value index for each variant
type VariantSelection struct {
	VariantName   string
	SelectedIndex int
}

type Model struct {
	screen            types.Screen
	previousScreen    types.Screen
	apiClient         *api.Client
	cart              types.Cart
	homeProducts      []types.Product
	searchResults     []types.Product
	categories        []types.Category
	currentProduct    *types.Product
	productQuantity   int
	variantSelections []VariantSelection
	variantFocusIndex int // 0=quantity, 1+=variants
	searchQuery       string
	cursor            int
	err               error
	loading           bool
	loadingMsg        string
	loadingFrame      int
	address           types.ShippingDetails
	order             *types.Order
	width             int
	height            int
	notification      *Notification
	viewport          viewport.Model
	viewportReady     bool
}

func NewModel(apiClient *api.Client) *Model {
	config.InitConfig()
	return &Model{
		screen:            types.ScreenHome,
		apiClient:         apiClient,
		cart:              types.Cart{Items: []types.CartItem{}},
		cursor:            0,
		productQuantity:   1,
		variantSelections: []VariantSelection{},
		variantFocusIndex: 0,
		width:             80,
		height:            24,
		viewportReady:     false,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func clearNotificationCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return notificationClearMsg{}
	})
}

func (m *Model) SetLoading(loading bool, msg string) tea.Cmd {
	m.loading = loading
	m.loadingMsg = msg
	m.loadingFrame = 0
	if loading {
		return tickCmd()
	}
	return nil
}

func (m *Model) SetError(err error) {
	m.err = err
	m.loading = false
}

func (m *Model) ClearError() {
	m.err = nil
}

func (m *Model) ResetCursor() {
	m.cursor = 0
}

func (m *Model) ClearCart() {
	m.cart = types.Cart{Items: []types.CartItem{}}
}

func (m *Model) NavigateUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) NavigateDown(maxIndex int) {
	if m.cursor < maxIndex {
		m.cursor++
	}
}

func (m *Model) GoToScreen(screen types.Screen) tea.Cmd {
	m.previousScreen = m.screen
	m.screen = screen
	m.ResetCursor()
	m.viewport.GotoTop()
	return tea.ClearScreen
}

func (m *Model) SetNotification(message, notifType string) tea.Cmd {
	m.notification = &Notification{
		Message: message,
		Type:    notifType,
	}
	return clearNotificationCmd()
}

func (m *Model) ClearNotification() {
	m.notification = nil
}

func (m *Model) GetCurrentProducts() []types.Product {
	if m.screen == types.ScreenSearch {
		return m.searchResults
	}
	return m.homeProducts
}

func (m *Model) IncreaseQuantity() {
	if m.productQuantity < 99 {
		m.productQuantity++
	}
}

func (m *Model) DecreaseQuantity() {
	if m.productQuantity > 1 {
		m.productQuantity--
	}
}

func (m *Model) ResetProductState() {
	m.productQuantity = 1
	m.variantSelections = []VariantSelection{}
	m.variantFocusIndex = 0
}

// InitVariantSelections initializes variant selections for current product
func (m *Model) InitVariantSelections() {
	m.variantSelections = []VariantSelection{}
	if m.currentProduct == nil {
		return
	}
	for _, variant := range m.currentProduct.ProductVariants {
		m.variantSelections = append(m.variantSelections, VariantSelection{
			VariantName:   variant.VariantName,
			SelectedIndex: 0,
		})
	}
}

// GetTotalFocusItems returns count of focusable items (quantity + variants)
func (m *Model) GetTotalFocusItems() int {
	if m.currentProduct == nil {
		return 1
	}
	return 1 + len(m.currentProduct.ProductVariants)
}

// MoveFocusUp moves focus up in product page
func (m *Model) MoveFocusUp() {
	if m.variantFocusIndex > 0 {
		m.variantFocusIndex--
	}
}

// MoveFocusDown moves focus down in product page
func (m *Model) MoveFocusDown() {
	maxIndex := m.GetTotalFocusItems() - 1
	if m.variantFocusIndex < maxIndex {
		m.variantFocusIndex++
	}
}

// CycleVariantLeft moves selection left for focused variant
func (m *Model) CycleVariantLeft() {
	if m.variantFocusIndex == 0 {
		m.DecreaseQuantity()
		return
	}
	
	variantIdx := m.variantFocusIndex - 1
	if variantIdx >= 0 && variantIdx < len(m.variantSelections) && m.currentProduct != nil {
		variant := m.currentProduct.ProductVariants[variantIdx]
		sel := &m.variantSelections[variantIdx]
		if sel.SelectedIndex > 0 {
			sel.SelectedIndex--
		} else {
			sel.SelectedIndex = len(variant.VariantValues) - 1
		}
	}
}

// CycleVariantRight moves selection right for focused variant
func (m *Model) CycleVariantRight() {
	if m.variantFocusIndex == 0 {
		m.IncreaseQuantity()
		return
	}
	
	variantIdx := m.variantFocusIndex - 1
	if variantIdx >= 0 && variantIdx < len(m.variantSelections) && m.currentProduct != nil {
		variant := m.currentProduct.ProductVariants[variantIdx]
		sel := &m.variantSelections[variantIdx]
		if sel.SelectedIndex < len(variant.VariantValues)-1 {
			sel.SelectedIndex++
		} else {
			sel.SelectedIndex = 0
		}
	}
}

// GetSelectedVariants returns a map of variant name to selected value
func (m *Model) GetSelectedVariants() map[string]string {
	result := make(map[string]string)
	if m.currentProduct == nil {
		return result
	}
	
	for i, sel := range m.variantSelections {
		if i < len(m.currentProduct.ProductVariants) {
			variant := m.currentProduct.ProductVariants[i]
			if sel.SelectedIndex < len(variant.VariantValues) {
				result[variant.VariantName] = variant.VariantValues[sel.SelectedIndex].Label
			}
		}
	}
	return result
}

// GetSelectedVariantString returns a formatted string of selected variants
func (m *Model) GetSelectedVariantString() string {
	variants := m.GetSelectedVariants()
	if len(variants) == 0 {
		return ""
	}
	
	result := ""
	for name, value := range variants {
		if result != "" {
			result += ", "
		}
		result += name + ": " + value
	}
	return result
}

// ContentWidth returns usable content width
func (m *Model) ContentWidth() int {
	w := m.width - 4 // padding
	if w < 40 {
		return 40
	}
	if w > 120 {
		return 120
	}
	return w
}

// InitViewport initializes the viewport with current dimensions
func (m *Model) InitViewport() {
	m.viewport = viewport.New(m.width, m.height)
	m.viewport.YPosition = 0
	m.viewport.SetContent("")
	m.viewportReady = true
}
