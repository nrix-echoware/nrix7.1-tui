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

// ==================== PRODUCT ====================

func (m *Model) renderProduct(w int) (header, content, footer string) {
	if m.currentProduct == nil {
		return "", "Product not found.", ""
	}
	p := m.currentProduct

	// HEADER (fixed)
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
	h.WriteString("\n\n")

	// Product title + price (in header)
	h.WriteString(TitleStyle.Render(p.Name))
	if p.Brand != "" {
		h.WriteString(HelpStyle.Render(fmt.Sprintf(" (%s)", p.Brand)))
	}
	h.WriteString("\n\n")

	// Price line
	priceStr := PriceStyle.Render(fmt.Sprintf("₹%.0f", p.SellingPrice))
	if p.MRPPrice > p.SellingPrice {
		discount := ((p.MRPPrice - p.SellingPrice) / p.MRPPrice) * 100
		priceStr += HelpStyle.Render(fmt.Sprintf("  ₹%.0f", p.MRPPrice))
		priceStr += SuccessStyle.Render(fmt.Sprintf("  (%.0f%% OFF)", discount))
	}
	h.WriteString(priceStr)
	h.WriteString("\n\n")
	header = h.String()

	// CONTENT (scrollable)
	var c strings.Builder

	// Description
	c.WriteString(m.divider(w))
	c.WriteString("\n")
	c.WriteString(SubtitleStyle.Render("DESCRIPTION"))
	c.WriteString("\n")
	c.WriteString(m.divider(w))
	c.WriteString("\n\n")
	if p.ProductDescription != "" {
		c.WriteString(wrapText(p.ProductDescription, w-2))
	} else {
		c.WriteString("No description available.")
	}
	c.WriteString("\n\n")

	// Features
	if len(p.Features) > 0 {
		c.WriteString(m.divider(w))
		c.WriteString("\n")
		c.WriteString(SubtitleStyle.Render("FEATURES"))
		c.WriteString("\n")
		c.WriteString(m.divider(w))
		c.WriteString("\n\n")
		for _, f := range p.Features {
			c.WriteString(fmt.Sprintf("• %s\n", f))
		}
		c.WriteString("\n")
	}

	// Options
	c.WriteString(m.divider(w))
	c.WriteString("\n")
	c.WriteString(SubtitleStyle.Render("OPTIONS"))
	c.WriteString("\n")
	c.WriteString(m.divider(w))
	c.WriteString("\n\n")

	// Quantity
	qtyFocused := m.variantFocusIndex == 0
	c.WriteString(m.renderOptionLine("Quantity", fmt.Sprintf("%d", m.productQuantity), qtyFocused))
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
		c.WriteString(m.renderOptionLine(variant.VariantName, strings.Join(opts, " "), focused))
		c.WriteString("\n")
	}
	c.WriteString("\n")

	// Tags
	if len(p.Tags) > 0 {
		c.WriteString(m.divider(w))
		c.WriteString("\n")
		c.WriteString(SubtitleStyle.Render("TAGS"))
		c.WriteString("\n")
		c.WriteString(m.divider(w))
		c.WriteString("\n\n")
		for _, tag := range p.Tags {
			c.WriteString(HelpStyle.Render(fmt.Sprintf("#%s  ", tag)))
		}
		c.WriteString("\n\n")
	}

	content = c.String()

	// FOOTER (fixed)
	footer = m.renderFooter("↑/↓ Scroll   ←/→ Change   A Add to Cart   C Cart   Esc Back", w)
	return
}

// ==================== CART ====================

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

	// CONTENT
	var c strings.Builder
	if len(m.cart.Items) == 0 {
		c.WriteString("Your cart is empty.\n\n")
		c.WriteString(HelpStyle.Render("Press Esc to browse products"))
		c.WriteString("\n")
	} else {
		for i, item := range m.cart.Items {
			selected := i == m.cursor
			c.WriteString(m.renderCartLine(item, selected, w))
			c.WriteString("\n")
		}
		c.WriteString("\n")
		c.WriteString(m.divider(w))
		c.WriteString("\n")
		c.WriteString(TitleStyle.Render(fmt.Sprintf("Total: ₹%.0f", m.cart.Total())))
		c.WriteString("\n")
	}
	content = c.String()

	// FOOTER
	footer = m.renderFooter("↑/↓ Navigate   +/- Qty   D Remove   Enter Checkout   Esc Back", w)
	return
}

// ==================== ADDRESS ====================

func (m *Model) renderAddress(w int) (header, content, footer string) {
	// HEADER
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	h.WriteString(m.headerRow("← Back", "SHIPPING DETAILS", w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n\n")
	header = h.String()

	// CONTENT
	var c strings.Builder
	c.WriteString(m.renderInputLine("Phone", m.address.Phone, m.cursor == 0, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Email", m.address.Email, m.cursor == 1, w))
	c.WriteString("\n\n")
	c.WriteString(m.renderInputLine("Address", m.address.Address, m.cursor == 2, w))
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
	h.WriteString(m.headerRow("← Back", "CONFIRM ORDER", w))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n\n")
	header = h.String()

	// CONTENT
	var c strings.Builder
	c.WriteString(SubtitleStyle.Render("ORDER SUMMARY"))
	c.WriteString("\n\n")
	for _, item := range m.cart.Items {
		total := item.Product.SellingPrice * float64(item.Quantity)
		c.WriteString(fmt.Sprintf("%-40s x%d    ₹%.0f\n", truncate(item.Product.Name, 40), item.Quantity, total))
	}
	c.WriteString("\n")
	c.WriteString(m.divider(w))
	c.WriteString("\n")
	c.WriteString(TitleStyle.Render(fmt.Sprintf("Total: ₹%.0f", m.cart.Total())))
	c.WriteString("\n\n")

	c.WriteString(SubtitleStyle.Render("SHIPPING TO"))
	c.WriteString("\n\n")
	c.WriteString(fmt.Sprintf("Phone:   %s\n", m.address.Phone))
	c.WriteString(fmt.Sprintf("Email:   %s\n", m.address.Email))
	c.WriteString(fmt.Sprintf("Address: %s\n", m.address.Address))
	c.WriteString("\n")
	c.WriteString(SuccessStyle.Render("Press Enter or Y to place order"))
	c.WriteString("\n")
	content = c.String()

	// FOOTER
	footer = m.renderFooter("Enter/Y Confirm   N/Esc Back", w)
	return
}

// ==================== ORDER SUCCESS ====================

func (m *Model) renderOrderSuccess(w int) (header, content, footer string) {
	// HEADER
	var h strings.Builder
	h.WriteString(m.divider(w))
	h.WriteString("\n")
	h.WriteString(lipgloss.NewStyle().Width(w).Align(lipgloss.Center).Render(SuccessStyle.Render("ORDER PLACED!")))
	h.WriteString("\n")
	h.WriteString(m.divider(w))
	h.WriteString("\n\n")
	header = h.String()

	// CONTENT
	var c strings.Builder
	if m.order != nil {
		c.WriteString(fmt.Sprintf("Order ID: %s\n", m.order.ID))
		c.WriteString(fmt.Sprintf("Total:    ₹%.0f\n", m.order.TotalAmount))
		c.WriteString(fmt.Sprintf("Status:   %s\n", m.order.Status.Type))
	}
	c.WriteString("\n")
	c.WriteString(SuccessStyle.Render("Thank you for your order!"))
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
