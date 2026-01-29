package tui

import (
	"fmt"
	"terminal-echoware/pkg/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	loadingCmd := m.SetLoading(true, "Loading products...")
	return tea.Batch(loadingCmd, loadProductsCmd(m.apiClient, 0, 20))
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.viewportReady {
			m.InitViewport()
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
		return m, nil

	case tickMsg:
		if m.loading {
			m.loadingFrame = (m.loadingFrame + 1) % len(LoadingFrames)
			return m, tickCmd()
		}
		return m, nil

	case notificationClearMsg:
		m.ClearNotification()
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		
		// Handle scroll keys for viewport
		switch msg.String() {
		case "pgup":
			m.viewport.HalfViewUp()
		case "pgdown":
			m.viewport.HalfViewDown()
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		}
		
		return m.handleKeyPress(msg)

	case productsLoadedMsg:
		return m.handleProductsLoaded(msg)

	case productLoadedMsg:
		return m.handleProductLoaded(msg)

	case searchResultsMsg:
		return m.handleSearchResults(msg)

	case orderCreatedMsg:
		return m.handleOrderCreated(msg)
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) handleProductsLoaded(msg productsLoadedMsg) (tea.Model, tea.Cmd) {
	m.SetLoading(false, "")
	if msg.err != nil {
		m.SetError(msg.err)
		return m, nil
	}
	m.homeProducts = msg.products
	m.ClearError()
	return m, tea.ClearScreen
}

func (m *Model) handleProductLoaded(msg productLoadedMsg) (tea.Model, tea.Cmd) {
	m.SetLoading(false, "")
	if msg.err != nil {
		m.SetError(msg.err)
		return m, nil
	}
	m.currentProduct = msg.product
	m.ResetProductState()
	m.InitVariantSelections()
	m.screen = types.ScreenProduct
	m.ClearError()
	return m, tea.ClearScreen
}

func (m *Model) handleSearchResults(msg searchResultsMsg) (tea.Model, tea.Cmd) {
	m.SetLoading(false, "")
	if msg.err != nil {
		m.SetError(msg.err)
		return m, nil
	}
	m.searchResults = msg.products
	m.ResetCursor()
	m.ClearError()
	return m, tea.ClearScreen
}

func (m *Model) handleOrderCreated(msg orderCreatedMsg) (tea.Model, tea.Cmd) {
	m.SetLoading(false, "")
	if msg.err != nil {
		m.SetError(msg.err)
		return m, nil
	}
	m.order = msg.order
	m.screen = types.ScreenOrderSuccess
	m.ClearError()
	return m, tea.ClearScreen
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case types.ScreenHome:
		return m.handleHomeKeys(msg)
	case types.ScreenSearch:
		return m.handleSearchKeys(msg)
	case types.ScreenProduct:
		return m.handleProductKeys(msg)
	case types.ScreenCart:
		return m.handleCartKeys(msg)
	case types.ScreenAddress:
		return m.handleAddressKeys(msg)
	case types.ScreenCheckout:
		return m.handleCheckoutKeys(msg)
	case types.ScreenOrderSuccess:
		return m.handleOrderSuccessKeys(msg)
	}
	return m, nil
}

func (m *Model) handleHomeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "up", "k":
		m.NavigateUp()
		return m, nil
	case "down", "j":
		m.NavigateDown(len(m.homeProducts) - 1)
		return m, nil
	case "enter", " ":
		if len(m.homeProducts) > 0 && m.cursor < len(m.homeProducts) {
			loadingCmd := m.SetLoading(true, "Loading product...")
			return m, tea.Batch(loadingCmd, loadProductCmd(m.apiClient, m.homeProducts[m.cursor].ID))
		}
		return m, nil
	case "s", "/":
		cmd := m.GoToScreen(types.ScreenSearch)
		m.searchQuery = ""
		m.searchResults = nil
		return m, cmd
	case "c":
		return m, m.GoToScreen(types.ScreenCart)
	}
	return m, nil
}

func (m *Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	
	// Handle special keys first
	switch key {
	case "esc":
		cmd := m.GoToScreen(types.ScreenHome)
		m.searchResults = nil
		return m, cmd
	case "ctrl+c":
		return m, tea.Quit
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
		return m, nil
	case "enter":
		// If we have search results and cursor is on a product, open it
		if len(m.searchResults) > 0 && m.cursor < len(m.searchResults) {
			loadingCmd := m.SetLoading(true, "Loading product...")
			return m, tea.Batch(loadingCmd, loadProductCmd(m.apiClient, m.searchResults[m.cursor].ID))
		}
		// Otherwise, perform search
		if len(m.searchQuery) > 0 {
			loadingCmd := m.SetLoading(true, fmt.Sprintf("Searching for '%s'...", m.searchQuery))
			return m, tea.Batch(loadingCmd, searchProductsCmd(m.apiClient, m.searchQuery, 0, 20))
		}
		return m, nil
	case "tab":
		// Tab to search with current query
		if len(m.searchQuery) > 0 {
			loadingCmd := m.SetLoading(true, fmt.Sprintf("Searching for '%s'...", m.searchQuery))
			return m, tea.Batch(loadingCmd, searchProductsCmd(m.apiClient, m.searchQuery, 0, 20))
		}
		return m, nil
	case "up":
		m.NavigateUp()
		return m, nil
	case "down":
		m.NavigateDown(len(m.searchResults) - 1)
		return m, nil
	}
	
	// All other characters go to search query
	if len(msg.Runes) > 0 {
		m.searchQuery += string(msg.Runes)
	}
	return m, nil
}

func (m *Model) handleProductKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc", "b":
		// Go back to previous screen
		if m.previousScreen == types.ScreenSearch {
			m.screen = types.ScreenSearch
		} else {
			m.screen = types.ScreenHome
		}
		m.currentProduct = nil
		m.ResetProductState()
		m.ResetCursor()
		return m, tea.ClearScreen
	case "a":
		return m.addToCart()
	case "enter":
		return m.addToCart()
	case "tab", "down", "j":
		m.MoveFocusDown()
		return m, nil
	case "shift+tab", "up", "k":
		m.MoveFocusUp()
		return m, nil
	case "left", "h", "-", "_":
		m.CycleVariantLeft()
		return m, nil
	case "right", "l", "+", "=":
		m.CycleVariantRight()
		return m, nil
	case "c":
		return m, m.GoToScreen(types.ScreenCart)
	}
	return m, nil
}

func (m *Model) addToCart() (tea.Model, tea.Cmd) {
	if m.currentProduct == nil {
		return m, nil
	}
	
	// Check if all variants are selected (for products with variants)
	if len(m.currentProduct.ProductVariants) > 0 {
		variantStr := m.GetSelectedVariantString()
		m.cart.Add(*m.currentProduct, m.productQuantity)
		return m, m.SetNotification(fmt.Sprintf("Added %d (%s) to cart!", m.productQuantity, variantStr), "success")
	}
	
	m.cart.Add(*m.currentProduct, m.productQuantity)
	return m, m.SetNotification(fmt.Sprintf("Added %d to cart!", m.productQuantity), "success")
}

func (m *Model) handleCartKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc", "b":
		return m, m.GoToScreen(types.ScreenHome)
	case "up", "k":
		m.NavigateUp()
		return m, nil
	case "down", "j":
		m.NavigateDown(len(m.cart.Items) - 1)
		return m, nil
	case "+", "=":
		if m.cursor < len(m.cart.Items) {
			item := &m.cart.Items[m.cursor]
			if item.Quantity < 99 {
				item.Quantity++
				return m, m.SetNotification("Quantity increased", "info")
			}
		}
		return m, nil
	case "-", "_":
		if m.cursor < len(m.cart.Items) {
			item := &m.cart.Items[m.cursor]
			if item.Quantity > 1 {
				item.Quantity--
				return m, m.SetNotification("Quantity decreased", "info")
			}
		}
		return m, nil
	case "d", "x":
		if m.cursor < len(m.cart.Items) {
			name := m.cart.Items[m.cursor].Product.Name
			m.cart.Remove(m.cart.Items[m.cursor].Product.ID)
			if m.cursor >= len(m.cart.Items) && m.cursor > 0 {
				m.cursor--
			}
			return m, m.SetNotification(fmt.Sprintf("Removed %s", truncate(name, 20)), "info")
		}
		return m, nil
	case "enter", " ":
		if len(m.cart.Items) > 0 {
			return m, m.GoToScreen(types.ScreenAddress)
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) handleAddressKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = types.ScreenCart
		return m, tea.ClearScreen
	case "enter":
		if m.address.Phone != "" && m.address.Email != "" && m.address.Address != "" {
			m.screen = types.ScreenCheckout
			return m, tea.ClearScreen
		}
		return m, nil
	case "tab", "down":
		m.cursor = (m.cursor + 1) % 3
		return m, nil
	case "up":
		m.cursor = (m.cursor + 2) % 3
		return m, nil
	case "backspace":
		m.handleAddressBackspace()
		return m, nil
	default:
		if len(msg.Runes) > 0 {
			m.handleAddressInput(string(msg.Runes))
		}
		return m, nil
	}
}

func (m *Model) handleAddressBackspace() {
	switch m.cursor {
	case 0:
		if len(m.address.Phone) > 0 {
			m.address.Phone = m.address.Phone[:len(m.address.Phone)-1]
		}
	case 1:
		if len(m.address.Email) > 0 {
			m.address.Email = m.address.Email[:len(m.address.Email)-1]
		}
	case 2:
		if len(m.address.Address) > 0 {
			m.address.Address = m.address.Address[:len(m.address.Address)-1]
		}
	}
}

func (m *Model) handleAddressInput(input string) {
	switch m.cursor {
	case 0:
		m.address.Phone += input
	case 1:
		m.address.Email += input
	case 2:
		m.address.Address += input
	}
}

func (m *Model) handleCheckoutKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.screen = types.ScreenAddress
		return m, tea.ClearScreen
	case "enter", "y":
		return m.placeOrder()
	case "n":
		m.screen = types.ScreenAddress
		return m, tea.ClearScreen
	}
	return m, nil
}

func (m *Model) placeOrder() (tea.Model, tea.Cmd) {
	var orderItems []types.OrderItem
	for _, item := range m.cart.Items {
		orderItems = append(orderItems, types.OrderItem{
			Product:  item.Product,
			Quantity: item.Quantity,
		})
	}

	params := types.OrderCreateParams{
		TotalAmount:     m.cart.Total(),
		TotalDiscount:   0,
		OrderItems:      orderItems,
		ShippingDetails: m.address,
	}

	loadingCmd := m.SetLoading(true, "Placing order...")
	return m, tea.Batch(loadingCmd, createOrderCmd(m.apiClient, params))
}

func (m *Model) handleOrderSuccessKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc", " ":
		cmd := m.GoToScreen(types.ScreenHome)
		m.ClearCart()
		m.order = nil
		m.address = types.ShippingDetails{}
		return m, cmd
	}
	return m, nil
}
