package tui

import (
	"fmt"
	"strings"
	"terminal-echoware/pkg/config"
	"terminal-echoware/pkg/types"

	"github.com/charmbracelet/lipgloss"
)

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	var currentLine strings.Builder
	currentLen := 0

	for _, word := range words {
		wordLen := len(word)
		if currentLen+wordLen+1 > width && currentLen > 0 {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLen = 0
		}
		if currentLen > 0 {
			currentLine.WriteString(" ")
			currentLen++
		}
		currentLine.WriteString(word)
		currentLen += wordLen
	}
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	return strings.Join(lines, "\n")
}

func (m *Model) divider(w int) string {
	return DividerStyle.Render(strings.Repeat("─", w))
}

func (m *Model) View() string {
	if m.loading {
		return m.renderLoading()
	}

	w := m.width
	if w < 40 {
		w = 40
	}

	var header, content, footer string

	switch m.screen {
	case types.ScreenHome:
		header, content, footer = m.renderHome(w)
	case types.ScreenSearch:
		header, content, footer = m.renderSearch(w)
	case types.ScreenProduct:
		header, content, footer = m.renderProduct(w)
	case types.ScreenCart:
		header, content, footer = m.renderCart(w)
	case types.ScreenAddress:
		header, content, footer = m.renderAddress(w)
	case types.ScreenCheckout:
		header, content, footer = m.renderCheckout(w)
	case types.ScreenOrderSuccess:
		header, content, footer = m.renderOrderSuccess(w)
	}

	// Add notification if present
	if m.notification != nil {
		footer = RenderNotification(m.notification) + "\n" + footer
	}

	// Add error if present
	if m.err != nil {
		footer = ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n" + footer
	}

	// Calculate viewport height
	headerHeight := strings.Count(header, "\n") + 1
	footerHeight := strings.Count(footer, "\n") + 1
	viewportHeight := m.height - headerHeight - footerHeight
	if viewportHeight < 5 {
		viewportHeight = 5
	}

	// Special handling for cart screen with sidebar
	if m.screen == types.ScreenCart {
		sidebarWidth := 28
		contentWidth := w - sidebarWidth - 3
		if contentWidth < 30 {
			contentWidth = w - 10
		}
		
		// Render sidebar
		sidebar := m.renderHotkeySidebar(viewportHeight)
		
		// Update viewport with main content only
		if m.viewportReady {
			m.viewport.Width = contentWidth
			m.viewport.Height = viewportHeight
			m.viewport.SetContent(content)
		}
		
		// Build final view with sidebar
		var b strings.Builder
		b.WriteString(header)
		
		// Combine sidebar and viewport side by side
		sidebarLines := strings.Split(strings.TrimRight(sidebar, "\n"), "\n")
		viewportContent := content
		if m.viewportReady {
			viewportContent = m.viewport.View()
		}
		viewportLines := strings.Split(strings.TrimRight(viewportContent, "\n"), "\n")
		
		maxLines := viewportHeight
		if len(viewportLines) > maxLines {
			maxLines = len(viewportLines)
		}
		
		for i := 0; i < maxLines; i++ {
			sidebarLine := ""
			if i < len(sidebarLines) {
				sidebarLine = sidebarLines[i]
			} else {
				// Pad sidebar line to maintain width
				sidebarLine = strings.Repeat(" ", sidebarWidth+2)
			}
			
			viewportLine := ""
			if i < len(viewportLines) {
				viewportLine = viewportLines[i]
			}
			
			b.WriteString(sidebarLine)
			b.WriteString(" │ ")
			b.WriteString(viewportLine)
			b.WriteString("\n")
		}
		
		b.WriteString(footer)
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Render(b.String())
	}

	// Special handling for product screen with sidebar
	if m.screen == types.ScreenProduct {
		sidebarWidth := 28
		separatorWidth := 3
		contentWidth := w - sidebarWidth - separatorWidth
		if contentWidth < 30 {
			contentWidth = 30
		}
		
		// Render sidebar
		sidebar := m.renderProductHotkeySidebar(viewportHeight)
		
		// Update viewport with main content only
		if m.viewportReady {
			m.viewport.Width = contentWidth
			m.viewport.Height = viewportHeight
			m.viewport.SetContent(content)
		}
		
		// Get viewport content
		viewportContent := content
		if m.viewportReady {
			viewportContent = m.viewport.View()
		}
		
		// Split into lines and trim
		sidebarLines := strings.Split(strings.TrimRight(sidebar, "\n"), "\n")
		viewportLines := strings.Split(strings.TrimRight(viewportContent, "\n"), "\n")
		
		// Get actual rendered sidebar width (accounting for ANSI codes)
		actualSidebarWidth := 0
		if len(sidebarLines) > 0 {
			actualSidebarWidth = lipgloss.Width(sidebarLines[0])
		}
		
		// Pad to same height
		maxLines := viewportHeight
		if len(viewportLines) > maxLines {
			maxLines = len(viewportLines)
		}
		if len(sidebarLines) < maxLines {
			emptySidebarLine := strings.Repeat(" ", actualSidebarWidth)
			for len(sidebarLines) < maxLines {
				sidebarLines = append(sidebarLines, emptySidebarLine)
			}
		}
		if len(viewportLines) < maxLines {
			for len(viewportLines) < maxLines {
				viewportLines = append(viewportLines, "")
			}
		}
		
		// Combine line by line
		var combinedLines []string
		for i := 0; i < maxLines && i < viewportHeight; i++ {
			sidebarLine := sidebarLines[i]
			viewportLine := ""
			if i < len(viewportLines) {
				viewportLine = viewportLines[i]
			}
			// Truncate viewport line if too long
			if lipgloss.Width(viewportLine) > contentWidth {
				viewportLine = viewportLine[:contentWidth]
			}
			combinedLine := lipgloss.JoinHorizontal(lipgloss.Left, sidebarLine, " │ ", viewportLine)
			combinedLines = append(combinedLines, combinedLine)
		}
		
		combinedContent := strings.Join(combinedLines, "\n")
		
		// Build final view
		var b strings.Builder
		b.WriteString(header)
		b.WriteString(combinedContent)
		if !strings.HasSuffix(combinedContent, "\n") {
			b.WriteString("\n")
		}
		b.WriteString(footer)
		
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Render(b.String())
	}

	// Update viewport
	if m.viewportReady {
		m.viewport.Width = w
		m.viewport.Height = viewportHeight
		m.viewport.SetContent(content)
	}

	// Build final view: header + viewport + footer
	var b strings.Builder
	b.WriteString(header)
	if m.viewportReady {
		b.WriteString(m.viewport.View())
	} else {
		b.WriteString(content)
	}
	b.WriteString("\n")
	b.WriteString(footer)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(b.String())
}

func (m *Model) renderLoading() string {
	cfg := config.GetConfig()

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(TitleStyle.Render(cfg.ShopName))
	b.WriteString("\n\n")

	// About Us
	b.WriteString(m.divider(50))
	b.WriteString("\n")
	b.WriteString(TitleStyle.Render("ABOUT US"))
	b.WriteString("\n")
	b.WriteString(m.divider(50))
	b.WriteString("\n\n")
	b.WriteString(wrapText(cfg.CompanyDescription, 50))
	b.WriteString("\n\n")
	b.WriteString(m.divider(50))
	b.WriteString("\n\n")

	frame := string(LoadingFrames[m.loadingFrame%len(LoadingFrames)])
	b.WriteString(LoadingStyle.Render(fmt.Sprintf("%s %s", frame, m.loadingMsg)))

	// Center everything
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(b.String())
}

// ==================== HOME ====================

func (m *Model) renderHome(w int) (header, content, footer string) {
	// HEADER
	cfg := config.GetConfig()
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	leftPart := TitleStyle.Render(cfg.ShopName)
	rightPart := ""
	if m.cart.Count() > 0 {
		rightPart = CartBadgeStyle.Render(fmt.Sprintf(" Cart(%d) ", m.cart.Count()))
	}
	h.WriteString(m.headerRow(leftPart, rightPart, w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n\n")
	header = h.String()

	// CONTENT
	var c strings.Builder
	if len(m.homeProducts) == 0 {
		c.WriteString("No products available.\n")
	} else {
		for i, p := range m.homeProducts {
			c.WriteString(m.renderProductLine(p, i == m.cursor, w))
			c.WriteString("\n")
		}
	}
	content = c.String()

	// FOOTER
	footer = m.renderFooter("↑/↓ Navigate   Enter View   S Search   C Cart   Q Quit", w)
	return
}

// ==================== SEARCH ====================

func (m *Model) renderSearch(w int) (header, content, footer string) {
	// HEADER
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	h.WriteString(m.headerRow("← Back", "SEARCH", w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n\n")
	h.WriteString(fmt.Sprintf("Search: %s▌\n", m.searchQuery))
	h.WriteString(HelpStyle.Render("Type to search • Tab to execute • Enter to select"))
	h.WriteString("\n\n")
	header = h.String()

	// CONTENT
	var c strings.Builder
	if len(m.searchResults) == 0 {
		if m.searchQuery == "" {
			c.WriteString("Start typing to search...\n")
		} else {
			c.WriteString("No results found.\n")
		}
	} else {
		c.WriteString(fmt.Sprintf("Found %d results:\n\n", len(m.searchResults)))
		for i, p := range m.searchResults {
			c.WriteString(m.renderProductLine(p, i == m.cursor, w))
			c.WriteString("\n")
		}
	}
	content = c.String()

	// FOOTER
	footer = m.renderFooter("↑/↓ Navigate   Enter Select   Esc Back", w)
	return
}

func (m *Model) renderProductHotkeySidebar(height int) string {
	sidebarWidth := 28
	var sb strings.Builder
	
	sb.WriteString(SubtitleStyle.Render("KEYBOARD SHORTCUTS"))
	sb.WriteString("\n")
	sb.WriteString(m.divider(sidebarWidth - 4))
	sb.WriteString("\n\n")
	
	hotkeys := []struct {
		key string
		desc string
	}{
		{"Tab", "Next Option"},
		{"Shift+Tab", "Prev Option"},
		{"← / →", "Change Value"},
		{"↑ / ↓", "Scroll"},
		{"A / Enter", "Add to Cart"},
		{"C", "View Cart"},
		{"Esc / B", "Back"},
		{"Q", "Quit"},
	}
	
	for _, hk := range hotkeys {
		keyPart := HelpStyle.Render(fmt.Sprintf("%-14s", hk.key))
		descPart := NormalStyle.Render(hk.desc)
		sb.WriteString(fmt.Sprintf("%s %s\n", keyPart, descPart))
	}
	
	lines := strings.Count(sb.String(), "\n")
	remaining := height - lines - 1
	if remaining > 0 {
		sb.WriteString(strings.Repeat("\n", remaining))
	}
	
	return lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorMuted).
		Padding(1, 1).
		Render(sb.String())
}

// ==================== PRODUCT ====================

func (m *Model) renderProduct(w int) (header, content, footer string) {
	if m.currentProduct == nil {
		return "", "Product not found.", ""
	}
	p := m.currentProduct

	// HEADER (compact)
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	cartStr := ""
	if m.cart.Count() > 0 {
		cartStr = fmt.Sprintf("Cart(%d)", m.cart.Count())
	}
	h.WriteString(m.headerRow("← Back", cartStr, w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	
	// Compact product title and price in header
	titleLine := p.Name
	if p.Brand != "" {
		titleLine = fmt.Sprintf("%s - %s", p.Brand, p.Name)
	}
	h.WriteString(TitleStyle.Render(titleLine))
	h.WriteString("  ")
	
	priceStr := PriceStyle.Render(fmt.Sprintf("₹%.0f", p.SellingPrice))
	if p.MRPPrice > p.SellingPrice {
		discount := ((p.MRPPrice - p.SellingPrice) / p.MRPPrice) * 100
		priceStr += " " + HelpStyle.Render(fmt.Sprintf("₹%.0f", p.MRPPrice))
		priceStr += " " + SuccessStyle.Render(fmt.Sprintf("%.0f%% OFF", discount))
	}
	h.WriteString(priceStr)
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	header = h.String()

	// CONTENT (compact, no boxes, minimal spacing)
	contentWidth := w - 28 - 3 // sidebar width + separator
	if contentWidth < 30 {
		contentWidth = 30
	}
	
	var c strings.Builder
	
	// Description (compact)
	c.WriteString(SubtitleStyle.Render("DESCRIPTION"))
	c.WriteString("\n")
	if p.ProductDescription != "" {
		c.WriteString(wrapText(p.ProductDescription, contentWidth-2))
	} else {
		c.WriteString(HelpStyle.Render("No description available."))
	}
	c.WriteString("\n\n")

	// Features (compact)
	if len(p.Features) > 0 {
		c.WriteString(SubtitleStyle.Render("FEATURES"))
		c.WriteString("\n")
		for i, f := range p.Features {
			c.WriteString(fmt.Sprintf("  • %s", f))
			if i < len(p.Features)-1 {
				c.WriteString("\n")
			}
		}
		c.WriteString("\n\n")
	}

	// Options (compact, no box)
	c.WriteString(SubtitleStyle.Render("OPTIONS"))
	c.WriteString("\n")

	// Quantity
	qtyFocused := m.variantFocusIndex == 0
	qtyLine := m.renderOptionLine("Quantity", fmt.Sprintf("%d", m.productQuantity), qtyFocused)
	c.WriteString(qtyLine)
	c.WriteString("\n")

	// Variants
	for i, variant := range p.ProductVariants {
		focused := m.variantFocusIndex == i+1
		selectedIdx := 0
		if i < len(m.variantSelections) {
			selectedIdx = m.variantSelections[i].SelectedIndex
		}

		var opts []string
		for j, val := range variant.VariantValues {
			if j == selectedIdx {
				opts = append(opts, OptionValueSelectedStyle.Render(val.Label))
			} else {
				opts = append(opts, OptionValueUnselectedStyle.Render(val.Label))
			}
		}
		variantLine := m.renderOptionLine(variant.VariantName, strings.Join(opts, " "), focused)
		c.WriteString(variantLine)
		c.WriteString("\n")
	}
	c.WriteString("\n")

	// Tags (compact)
	if len(p.Tags) > 0 {
		c.WriteString(SubtitleStyle.Render("TAGS"))
		c.WriteString("\n")
		for _, tag := range p.Tags {
			c.WriteString(BadgeStyle.Render(" #" + tag + " "))
		}
		c.WriteString("\n")
	}

	content = c.String()

	// FOOTER (minimal)
	footer = m.renderFooter("", w)
	return
}

// ==================== CART ====================

func (m *Model) renderHotkeySidebar(height int) string {
	sidebarWidth := 28
	var sb strings.Builder
	
	sb.WriteString(SubtitleStyle.Render("KEYBOARD SHORTCUTS"))
	sb.WriteString("\n")
	sb.WriteString(m.divider(sidebarWidth - 4))
	sb.WriteString("\n\n")
	
	hotkeys := []struct {
		key string
		desc string
	}{
		{"↑ / k", "Navigate Up"},
		{"↓ / j", "Navigate Down"},
		{"+ / =", "Increase Qty"},
		{"- / _", "Decrease Qty"},
		{"D / x", "Remove Item"},
		{"Enter", "Checkout"},
		{"Esc / b", "Back"},
		{"Q", "Quit"},
	}
	
	for _, hk := range hotkeys {
		keyPart := HelpStyle.Render(fmt.Sprintf("%-12s", hk.key))
		descPart := NormalStyle.Render(hk.desc)
		sb.WriteString(fmt.Sprintf("%s %s\n", keyPart, descPart))
	}
	
	// Fill remaining height
	lines := strings.Count(sb.String(), "\n")
	remaining := height - lines - 1
	if remaining > 0 {
		sb.WriteString(strings.Repeat("\n", remaining))
	}
	
	return lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorMuted).
		Padding(1, 1).
		Render(sb.String())
}

func (m *Model) renderCart(w int) (header, content, footer string) {
	// HEADER
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	h.WriteString(m.headerRow("← Back", "SHOPPING CART", w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n\n")
	header = h.String()

	// CONTENT (main content only, sidebar handled in View)
	contentWidth := w - 28 - 3 // sidebar width + spacing
	if contentWidth < 30 {
		contentWidth = w - 10
	}
	
	var c strings.Builder
	if len(m.cart.Items) == 0 {
		c.WriteString("Your cart is empty.\n\n")
		c.WriteString(HelpStyle.Render("Press Esc to browse products"))
		c.WriteString("\n")
	} else {
		for i, item := range m.cart.Items {
			selected := i == m.cursor
			c.WriteString(m.renderCartLine(item, selected, contentWidth))
			c.WriteString("\n")
		}
		c.WriteString("\n")
		c.WriteString(m.divider(contentWidth))
		c.WriteString("\n")
		c.WriteString(TitleStyle.Render(fmt.Sprintf("Total: ₹%.0f", m.cart.Total())))
		c.WriteString("\n")
	}
	content = c.String()

	// FOOTER
	footer = m.renderFooter("", w)
	return
}

// ==================== ADDRESS ====================

func (m *Model) renderAddress(w int) (header, content, footer string) {
	// HEADER
	cfg := config.GetConfig()
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	h.WriteString(m.headerRow("← Back", cfg.ShopName, w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n\n")
	header = h.String()

	// CONTENT
	var c strings.Builder
	c.WriteString(TitleStyle.Render("SHIPPING DETAILS"))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Full Name", m.address.FullName, m.cursor == 0, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Phone", m.address.Phone, m.cursor == 1, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Email", m.address.Email, m.cursor == 2, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Address 1", m.address.AddressLine1, m.cursor == 3, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Address 2", m.address.AddressLine2, m.cursor == 4, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("City", m.address.City, m.cursor == 5, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("State", m.address.State, m.cursor == 6, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Postal", m.address.PostalCode, m.cursor == 7, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Country", m.address.Country, m.cursor == 8, w))
	c.WriteString("\n")
	content = c.String()

	// FOOTER
	footer = m.renderFooter("Tab/↓ Next   ↑ Previous   Enter Continue   Esc Back", w)
	return
}

// ==================== CHECKOUT ====================

func (m *Model) renderCheckout(w int) (header, content, footer string) {
	// HEADER
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	h.WriteString(m.headerRow("← Back", "CHECKOUT", w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	header = h.String()

	// Two-partition layout: Order Summary (left) | Shipping Address (right)
	leftWidth := (w - 3) / 2
	rightWidth := w - leftWidth - 3
	
	// LEFT: Order Summary (boxed)
	var leftBox strings.Builder
	leftBox.WriteString(SubtitleStyle.Render("ORDER SUMMARY"))
	leftBox.WriteString("\n")
	
	// Items list
	maxItems := 10
	itemsToShow := m.cart.Items
	if len(itemsToShow) > maxItems {
		itemsToShow = itemsToShow[:maxItems]
	}
	
	for i, item := range itemsToShow {
		total := item.Product.SellingPrice * float64(item.Quantity)
		name := truncate(item.Product.Name, leftWidth-20)
		leftBox.WriteString(fmt.Sprintf("  %d. %s\n", i+1, name))
		leftBox.WriteString(fmt.Sprintf("     Qty: %d  %s\n", item.Quantity, PriceStyle.Render(fmt.Sprintf("₹%.0f", total))))
		if i < len(itemsToShow)-1 {
			leftBox.WriteString("\n")
		}
	}
	
	if len(m.cart.Items) > maxItems {
		leftBox.WriteString(fmt.Sprintf("\n  ... and %d more item(s)\n", len(m.cart.Items)-maxItems))
	}
	
	leftBox.WriteString("\n")
	leftBox.WriteString(m.divider(leftWidth - 4))
	leftBox.WriteString("\n")
	totalLine := fmt.Sprintf("  Total: %s", PriceStyle.Render(fmt.Sprintf("₹%.0f", m.cart.Total())))
	leftBox.WriteString(TitleStyle.Render(totalLine))
	
	leftBoxRendered := BoxStyle.
		Width(leftWidth - 2).
		BorderForeground(ColorPrimary).
		Render(leftBox.String())
	
	// RIGHT: Shipping Address (boxed)
	var rightBox strings.Builder
	rightBox.WriteString(SubtitleStyle.Render("SHIPPING ADDRESS"))
	rightBox.WriteString("\n\n")
	
	rightBox.WriteString(fmt.Sprintf("  %s: %s\n", HelpStyle.Render("Name"), NormalStyle.Render(m.address.FullName)))
	rightBox.WriteString(fmt.Sprintf("  %s: %s\n", HelpStyle.Render("Phone"), NormalStyle.Render(m.address.Phone)))
	rightBox.WriteString(fmt.Sprintf("  %s: %s\n", HelpStyle.Render("Email"), NormalStyle.Render(m.address.Email)))
	rightBox.WriteString("\n")
	rightBox.WriteString(fmt.Sprintf("  %s:\n", HelpStyle.Render("Address")))
	rightBox.WriteString(fmt.Sprintf("    %s\n", NormalStyle.Render(m.address.AddressLine1)))
	if strings.TrimSpace(m.address.AddressLine2) != "" {
		rightBox.WriteString(fmt.Sprintf("    %s\n", NormalStyle.Render(m.address.AddressLine2)))
	}
	rightBox.WriteString(fmt.Sprintf("    %s, %s %s\n", 
		NormalStyle.Render(m.address.City),
		NormalStyle.Render(m.address.State),
		NormalStyle.Render(m.address.PostalCode)))
	rightBox.WriteString(fmt.Sprintf("    %s\n", NormalStyle.Render(m.address.Country)))
	rightBox.WriteString("\n")
	rightBox.WriteString(SuccessStyle.Render("  Enter/Y to confirm"))
	
	rightBoxRendered := BoxStyle.
		Width(rightWidth - 2).
		BorderForeground(ColorSecondary).
		Render(rightBox.String())
	
	// Combine left and right using lipgloss
	combinedContent := lipgloss.JoinHorizontal(lipgloss.Top, leftBoxRendered, " │ ", rightBoxRendered)
	
	content = combinedContent + "\n"

	// FOOTER
	footer = m.renderFooter("Enter/Y Confirm   N/Esc Back", w)
	return
}

// ==================== ORDER SUCCESS ====================

func (m *Model) renderOrderSuccess(w int) (header, content, footer string) {
	// HEADER
	cfg := config.GetConfig()
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	h.WriteString(m.headerRow("", cfg.ShopName, w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	header = h.String()

	// CONTENT
	var c strings.Builder
	
	// Success message
	c.WriteString("\n")
	successMsg := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center).
		Foreground(ColorAccent).
		Bold(true).
		Render("✓ ORDER PLACED SUCCESSFULLY!")
	c.WriteString(successMsg)
	c.WriteString("\n\n")
	
	// Order details box
	if m.order != nil {
		var orderBox strings.Builder
		orderBox.WriteString(SubtitleStyle.Render("ORDER DETAILS"))
		orderBox.WriteString("\n")
		orderBox.WriteString(m.divider(w - 4))
		orderBox.WriteString("\n")
		orderBox.WriteString(fmt.Sprintf("  %s: %s\n", HelpStyle.Render("Order ID"), TitleStyle.Render(m.order.ID)))
		orderBox.WriteString(fmt.Sprintf("  %s: %s\n", HelpStyle.Render("Total"), PriceStyle.Render(fmt.Sprintf("₹%.0f", m.order.TotalAmount))))
		orderBox.WriteString(fmt.Sprintf("  %s: %s\n", HelpStyle.Render("Status"), SuccessStyle.Render(string(m.order.Status.Type))))
		orderBox.WriteString(m.divider(w - 4))
		
		boxContent := orderBox.String()
		orderBoxRendered := BoxStyle.
			Width(w - 4).
			BorderForeground(ColorAccent).
			Render(boxContent)
		c.WriteString(orderBoxRendered)
		c.WriteString("\n\n")
	}
	
	// Thank you message
	thankYouMsg := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center).
		Foreground(ColorSecondary).
		Render("Thank you for your order!")
	c.WriteString(thankYouMsg)
	c.WriteString("\n")
	content = c.String()

	// FOOTER
	footer = m.renderFooter("Press any key to continue", w)
	return
}

// ==================== HELPER RENDERERS ====================

func (m *Model) headerRow(left, right string, w int) string {
	leftLen := lipgloss.Width(left)
	rightLen := lipgloss.Width(right)
	space := w - leftLen - rightLen
	if space < 1 {
		space = 1
	}
	return left + strings.Repeat(" ", space) + right
}

func (m *Model) renderFooter(helpText string, w int) string {
	var f strings.Builder
	f.WriteString(m.divider(w))
	f.WriteString("\n")
	f.WriteString(HelpStyle.Render(helpText))
	f.WriteString("\n")
	f.WriteString(m.divider(w))
	return f.String()
}

func (m *Model) renderProductLine(p types.Product, selected bool, w int) string {
	cursor := "  "
	style := NormalStyle
	if selected {
		cursor = "▸ "
		style = SelectedStyle
	}

	nameW := w - 30
	if nameW < 20 {
		nameW = 20
	}
	name := truncate(p.Name, nameW)
	price := fmt.Sprintf("₹%.0f", p.SellingPrice)

	line := fmt.Sprintf("%s%-*s  %s", cursor, nameW, name, PriceStyle.Render(price))
	return style.Render(line)
}

func (m *Model) renderCartLine(item types.CartItem, selected bool, w int) string {
	cursor := "  "
	style := NormalStyle
	if selected {
		cursor = "▸ "
		style = SelectedStyle
	}

	nameW := w - 35
	if nameW < 15 {
		nameW = 15
	}
	name := truncate(item.Product.Name, nameW)
	total := item.Product.SellingPrice * float64(item.Quantity)

	qtyStr := fmt.Sprintf("[-] %2d [+]", item.Quantity)
	if selected {
		qtyStr = SuccessStyle.Render(qtyStr)
	} else {
		qtyStr = HelpStyle.Render(qtyStr)
	}

	line := fmt.Sprintf("%s%-*s  %s  %s", cursor, nameW, name, qtyStr, PriceStyle.Render(fmt.Sprintf("₹%.0f", total)))
	return style.Render(line)
}

func (m *Model) renderOptionLine(label, value string, focused bool) string {
	style := NormalStyle
	prefix := "  "
	if focused {
		style = SelectedStyle
		prefix = "▸ "
	}
	return style.Render(fmt.Sprintf("%s%-12s: %s", prefix, label, value))
}

func (m *Model) renderInputLine(label, value string, focused bool, w int) string {
	cursor := ""
	style := NormalStyle
	if focused {
		cursor = "▌"
		style = SelectedStyle
	}
	return style.Width(w - 4).Render(fmt.Sprintf("%-10s: %s%s", label, value, cursor))
}
