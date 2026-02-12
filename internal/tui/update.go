package tui

import (
	"fmt"
	"strings"
	"time"
	"terminal-echoware/pkg/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	loadingCmd := m.SetLoading(true, "Loading products...")
	return tea.Batch(tea.ClearScreen, loadingCmd, loadProductsCmd(m.apiClient, 0, 20))
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
	m.viewport.GotoTop()
	m.viewport.SetContent("")
	return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
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
	m.viewport.GotoTop()
	m.viewport.SetContent("")
	return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
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
	m.viewport.GotoTop()
	m.viewport.SetContent("")
	return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
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
	m.viewport.GotoTop()
	m.viewport.SetContent("")
	return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
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
		m.viewport.GotoTop()
		m.viewport.SetContent("")
		return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
	case "a":
		return m.addToCart()
	case "enter":
		return m.addToCart()
	case "tab":
		m.MoveFocusDown()
		return m, nil
	case "shift+tab":
		m.MoveFocusUp()
		return m, nil
	case "up", "k":
		// Scroll viewport up
		m.viewport.LineUp(1)
		return m, nil
	case "down", "j":
		// Scroll viewport down
		m.viewport.LineDown(1)
		return m, nil
	case "left", "h":
		m.CycleVariantLeft()
		return m, nil
	case "right", "l":
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
		err := m.cart.Add(*m.currentProduct, m.productQuantity, m.GetSelectedVariants())
		if err != nil {
			return m, m.SetNotification(err.Error(), "error")
		}
		return m, m.SetNotification(fmt.Sprintf("Added %d (%s) to cart!", m.productQuantity, variantStr), "success")
	}
	
	err := m.cart.Add(*m.currentProduct, m.productQuantity, nil)
	if err != nil {
		return m, m.SetNotification(err.Error(), "error")
	}
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
			const maxQuantity = 5
			if item.Quantity >= maxQuantity {
				return m, m.SetNotification(fmt.Sprintf("Maximum quantity of %d reached for this item", maxQuantity), "error")
			}
			item.Quantity++
			return m, m.SetNotification("Quantity increased", "info")
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
			m.cart.Remove(m.cart.Items[m.cursor].Product.ID, m.cart.Items[m.cursor].Variant)
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
		m.viewport.GotoTop()
		m.viewport.SetContent("")
		return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
	case "enter":
		if errMsg := m.validateShippingDetails(); errMsg != "" {
			return m, m.SetNotification(errMsg, "error")
		}
		m.screen = types.ScreenCheckout
		m.viewport.GotoTop()
		m.viewport.SetContent("")
		return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
	case "tab", "down":
		m.cursor = (m.cursor + 1) % 9
		return m, nil
	case "up":
		m.cursor = (m.cursor + 8) % 9
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
		if len(m.address.FullName) > 0 {
			m.address.FullName = m.address.FullName[:len(m.address.FullName)-1]
		}
	case 1:
		if len(m.address.Phone) > 0 {
			m.address.Phone = m.address.Phone[:len(m.address.Phone)-1]
		}
	case 2:
		if len(m.address.Email) > 0 {
			m.address.Email = m.address.Email[:len(m.address.Email)-1]
		}
	case 3:
		if len(m.address.AddressLine1) > 0 {
			m.address.AddressLine1 = m.address.AddressLine1[:len(m.address.AddressLine1)-1]
		}
	case 4:
		if len(m.address.AddressLine2) > 0 {
			m.address.AddressLine2 = m.address.AddressLine2[:len(m.address.AddressLine2)-1]
		}
	case 5:
		if len(m.address.City) > 0 {
			m.address.City = m.address.City[:len(m.address.City)-1]
		}
	case 6:
		if len(m.address.State) > 0 {
			m.address.State = m.address.State[:len(m.address.State)-1]
		}
	case 7:
		if len(m.address.PostalCode) > 0 {
			m.address.PostalCode = m.address.PostalCode[:len(m.address.PostalCode)-1]
		}
	case 8:
		if len(m.address.Country) > 0 {
			m.address.Country = m.address.Country[:len(m.address.Country)-1]
		}
	}
}

func (m *Model) handleAddressInput(input string) {
	switch m.cursor {
	case 0:
		m.address.FullName += input
	case 1:
		m.address.Phone += input
	case 2:
		m.address.Email += input
	case 3:
		m.address.AddressLine1 += input
	case 4:
		m.address.AddressLine2 += input
	case 5:
		m.address.City += input
	case 6:
		m.address.State += input
	case 7:
		m.address.PostalCode += input
	case 8:
		m.address.Country += input
	}
}

func (m *Model) handleCheckoutKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.screen = types.ScreenAddress
		m.viewport.GotoTop()
		m.viewport.SetContent("")
		return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
	case "enter", "y":
		if errMsg := m.validateShippingDetails(); errMsg != "" {
			return m, m.SetNotification(errMsg, "error")
		}
		return m.placeOrder()
	case "n":
		m.screen = types.ScreenAddress
		m.viewport.GotoTop()
		m.viewport.SetContent("")
		return m, tea.Sequence(tea.ClearScreen, tea.WindowSize())
	}
	return m, nil
}

func (m *Model) placeOrder() (tea.Model, tea.Cmd) {
	var orderItems []types.OrderItemInput
	for _, item := range m.cart.Items {
		total := item.Product.SellingPrice * float64(item.Quantity)
		orderItems = append(orderItems, types.OrderItemInput{
			ProductID:   item.Product.ID,
			ProductName: item.Product.Name,
			Variant:     item.Variant,
			Quantity:    item.Quantity,
			Price:       item.Product.SellingPrice,
			Total:       total,
		})
	}

	subtotal := m.cart.Total()
	discount := 0.0
	shipping := 0.0
	total := subtotal - discount + shipping

	m.address.Address = ""
	m.address.IsDefault = false

	params := types.OrderCreateParams{
		ShippingAddress: m.address,
		Items:           orderItems,
		SpecialMessage:  "",
		Pricing: types.OrderPricingInput{
			Subtotal: subtotal,
			Discount: discount,
			Shipping: shipping,
			Total:    total,
		},
		UserEmail:     m.address.Email,
		Timestamp:     time.Now().Format(time.RFC3339),
		PaymentMethod: "cod",
	}

	loadingCmd := m.SetLoading(true, "Placing order...")
	return m, tea.Batch(loadingCmd, createOrderCmd(m.apiClient, params))
}

func (m *Model) validateShippingDetails() string {
	fullName := strings.TrimSpace(m.address.FullName)
	phone := strings.TrimSpace(m.address.Phone)
	email := strings.TrimSpace(m.address.Email)
	address1 := strings.TrimSpace(m.address.AddressLine1)
	city := strings.TrimSpace(m.address.City)
	state := strings.TrimSpace(m.address.State)
	postal := strings.TrimSpace(m.address.PostalCode)
	country := strings.TrimSpace(m.address.Country)

	if fullName == "" {
		return "Full name is required"
	}
	if phone == "" || !isValidPhone(phone) {
		return "Enter a valid phone number"
	}
	if email == "" || !isValidEmail(email) {
		return "Enter a valid email address"
	}
	if address1 == "" {
		return "Address line 1 is required"
	}
	if city == "" {
		return "City is required"
	}
	if state == "" {
		return "State is required"
	}
	if postal == "" || len(postal) < 4 {
		return "Enter a valid postal code"
	}
	if country == "" {
		return "Country is required"
	}
	return ""
}

func isValidEmail(email string) bool {
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1
}

func isValidPhone(phone string) bool {
	digits := 0
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits++
		} else if r != '+' && r != '-' && r != ' ' {
			return false
		}
	}
	return digits >= 8
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
