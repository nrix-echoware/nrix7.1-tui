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

func (m *Model) View() string {
	if m.loading {
		return m.renderLoading()
	}

	w := m.ContentWidth()
	
	var content string
	switch m.screen {
	case types.ScreenHome:
		content = m.renderHome(w)
	case types.ScreenSearch:
		content = m.renderSearch(w)
	case types.ScreenProduct:
		content = m.renderProduct(w)
	case types.ScreenCart:
		content = m.renderCart(w)
	case types.ScreenAddress:
		content = m.renderAddress(w)
	case types.ScreenCheckout:
		content = m.renderCheckout(w)
	case types.ScreenOrderSuccess:
		content = m.renderOrderSuccess(w)
	}

	// Add notification if present
	if m.notification != nil {
		content += "\n" + RenderNotification(m.notification)
	}

	// Add error if present
	if m.err != nil {
		content += "\n" + ErrorStyle.Render(fmt.Sprintf("⚠ Error: %v", m.err))
	}

	// Use viewport for scrollable content
	if m.viewportReady {
		m.viewport.SetContent(content)
		return m.viewport.View()
	}

	return content
}

func (m *Model) renderLoading() string {
	w := m.ContentWidth()
	cfg := config.GetConfig()
	
	var b strings.Builder
	b.WriteString("\n\n")
	b.WriteString(TitleStyle.Render(AsciiLogo))
	b.WriteString("\n\n")
	
	// About Us section
	aboutTitle := SubtitleStyle.Render("━━━ About Us ━━━")
	b.WriteString(aboutTitle)
	b.WriteString("\n\n")
	b.WriteString(BoxStyle.Width(w - 8).Render(wrapText(cfg.CompanyDescription, w-16)))
	b.WriteString("\n\n")
	
	frame := string(LoadingFrames[m.loadingFrame%len(LoadingFrames)])
	loadingText := LoadingStyle.Render(fmt.Sprintf("%s %s", frame, m.loadingMsg))
	b.WriteString(loadingText)
	
	// Center everything both horizontally and vertically
	content := b.String()
	
	centeredStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)
	
	return centeredStyle.Render(content)
}

func (m *Model) renderHome(w int) string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(TitleStyle.Render(AsciiLogo))
	b.WriteString("\n")
	
	// Cart badge on the right
	if m.cart.Count() > 0 {
		b.WriteString(CartBadgeStyle.Render(fmt.Sprintf(" 🛒 Cart: %d items ", m.cart.Count())))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Products section header
	b.WriteString(SubtitleStyle.Width(w).Render("━━━ Products ━━━"))
	b.WriteString("\n\n")

	// Product list
	if len(m.homeProducts) == 0 {
		b.WriteString(NormalStyle.Render("  No products available."))
		b.WriteString("\n")
	} else {
		b.WriteString(RenderProductList(m.homeProducts, m.cursor, w))
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderDivider(w))
	b.WriteString("\n")
	b.WriteString(RenderHelp("home", w))

	return b.String()
}

func (m *Model) renderSearch(w int) string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("🔍 Search Products", m.cart.Count(), w))
	b.WriteString("\n\n")

	// Search input
	b.WriteString(RenderInputField("Search", m.searchQuery, true, w))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("  Type to search • Tab: Execute • Enter: Select • ↑↓: Navigate"))
	b.WriteString("\n\n")

	// Results
	if len(m.searchResults) == 0 {
		if m.searchQuery == "" {
			b.WriteString(NormalStyle.Render("  Start typing to search..."))
		} else {
			b.WriteString(NormalStyle.Render("  No results found. Press Tab to search."))
		}
		b.WriteString("\n")
	} else {
		b.WriteString(SubtitleStyle.Render(fmt.Sprintf("━━━ Found %d results ━━━", len(m.searchResults))))
		b.WriteString("\n\n")
		b.WriteString(RenderProductList(m.searchResults, m.cursor, w))
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderDivider(w))
	b.WriteString("\n")
	b.WriteString(RenderHelp("search", w))

	return b.String()
}

func (m *Model) renderProduct(w int) string {
	if m.currentProduct == nil {
		return ErrorStyle.Render("Product not found.")
	}

	var b strings.Builder
	p := m.currentProduct

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("📦 Product Details", m.cart.Count(), w))
	b.WriteString("\n\n")

	// ============ ROW 1: Product Name + Description (side by side) ============
	leftColWidth := w / 2
	rightColWidth := w - leftColWidth - 2

	// Left: Product Name Box
	leftContent := BoxStyle.Width(leftColWidth - 2).Render(
		TitleStyle.Render(wrapText(p.Name, leftColWidth-8)),
	)

	// Right: Description Box (500 chars max)
	desc := p.ProductDescription
	if len(desc) > 500 {
		desc = desc[:500] + "..."
	}
	rightContent := BoxStyle.Width(rightColWidth - 2).Render(
		SubtitleStyle.Render("📝 Description") + "\n\n" +
			NormalStyle.Render(wrapText(desc, rightColWidth-8)),
	)

	// Join horizontally
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, "  ", rightContent)
	b.WriteString(row1)
	b.WriteString("\n\n")

	// ============ ROW 2: Features ============
	if len(p.Features) > 0 {
		var featuresContent strings.Builder
		featuresContent.WriteString(SubtitleStyle.Render("✨ Features"))
		featuresContent.WriteString("\n\n")
		for _, feature := range p.Features {
			featuresContent.WriteString(fmt.Sprintf("  • %s\n", feature))
		}
		b.WriteString(BoxStyle.Width(w - 2).Render(featuresContent.String()))
		b.WriteString("\n\n")
	}

	// ============ ROW 3: Brand + Price Info ============
	var priceInfo strings.Builder
	priceInfo.WriteString(fmt.Sprintf("🏷️  Brand: %s\n\n", BrandStyle.Render(p.Brand)))
	priceInfo.WriteString(fmt.Sprintf("💰 Selling Price: %s\n", PriceStyle.Render(fmt.Sprintf("₹%.0f", p.SellingPrice))))
	priceInfo.WriteString(fmt.Sprintf("📋 MRP: ₹%.0f", p.MRPPrice))
	if p.MRPPrice > p.SellingPrice {
		discount := ((p.MRPPrice - p.SellingPrice) / p.MRPPrice) * 100
		priceInfo.WriteString(fmt.Sprintf("\n\n🎉 %s", SuccessStyle.Render(fmt.Sprintf("%.0f%% OFF!", discount))))
	}
	b.WriteString(BoxStyle.Width(w - 2).Render(priceInfo.String()))
	b.WriteString("\n\n")

	// ============ ROW 4: Tags + Categories ============
	var tagsContent strings.Builder
	if len(p.Tags) > 0 {
		tagsContent.WriteString("🏷️  Tags: ")
		for i, tag := range p.Tags {
			if i > 0 {
				tagsContent.WriteString(" • ")
			}
			tagsContent.WriteString(BadgeStyle.Render(" " + tag + " "))
		}
	}
	if len(p.CategoryDetails) > 0 {
		if tagsContent.Len() > 0 {
			tagsContent.WriteString("\n\n")
		}
		tagsContent.WriteString("📂 Categories: ")
		var catNames []string
		for _, cat := range p.CategoryDetails {
			catNames = append(catNames, cat.Name)
		}
		tagsContent.WriteString(strings.Join(catNames, ", "))
	}
	if tagsContent.Len() > 0 {
		b.WriteString(BoxStyle.Width(w - 2).Render(tagsContent.String()))
		b.WriteString("\n\n")
	}

	// ============ ROW 5: Select Options ============
	var optionsContent strings.Builder
	optionsContent.WriteString(SubtitleStyle.Render("⚙️  Select Options"))
	optionsContent.WriteString("\n\n")

	// Quantity selector
	qtyFocused := m.variantFocusIndex == 0
	optionsContent.WriteString(RenderOptionRow("Quantity", fmt.Sprintf("%d", m.productQuantity), qtyFocused, w-8))
	optionsContent.WriteString("\n")

	// Variant selectors
	for i, variant := range p.ProductVariants {
		focused := m.variantFocusIndex == i+1

		options := []string{}
		for _, val := range variant.VariantValues {
			options = append(options, val.Label)
		}

		selectedIdx := 0
		if i < len(m.variantSelections) {
			selectedIdx = m.variantSelections[i].SelectedIndex
		}

		optionsContent.WriteString(RenderVariantRow(variant.VariantName, options, selectedIdx, focused, w-8))
		optionsContent.WriteString("\n")
	}

	b.WriteString(BoxStyle.Width(w - 2).Render(optionsContent.String()))

	// Footer
	b.WriteString("\n\n")
	b.WriteString(RenderDivider(w))
	b.WriteString("\n")
	helpText := "Tab/↑↓: Navigate • ←→: Change • A/Enter: Add • C: Cart • Esc: Back"
	b.WriteString(FooterStyle.Width(w).Render(HelpStyle.Render(helpText)))

	return b.String()
}

func (m *Model) renderCart(w int) string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("🛒 Shopping Cart", 0, w))
	b.WriteString("\n\n")

	if len(m.cart.Items) == 0 {
		b.WriteString(NormalStyle.Render("  Your cart is empty."))
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("  Press Esc to browse products"))
		b.WriteString("\n")
	} else {
		// Items header
		b.WriteString(SubtitleStyle.Width(w).Render("━━━ Items ━━━"))
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("  Use +/- to change quantity"))
		b.WriteString("\n\n")
		
		// Items list
		b.WriteString(RenderCartListWithQty(m.cart.Items, m.cursor, w))
		
		// Total
		b.WriteString("\n")
		totalBox := BoxStyle.Width(w - 4).Render(fmt.Sprintf("💰 Total: %s", RenderPrice(m.cart.Total())))
		b.WriteString(totalBox)
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderDivider(w))
	b.WriteString("\n")
	helpText := "↑↓/jk: Navigate • +/-: Qty • D: Remove • Enter: Checkout • Esc: Back"
	b.WriteString(FooterStyle.Width(w).Render(HelpStyle.Render(helpText)))

	return b.String()
}

func (m *Model) renderAddress(w int) string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("📦 Shipping Details", m.cart.Count(), w))
	b.WriteString("\n\n")

	// Form fields
	b.WriteString(RenderInputField("📱 Phone", m.address.Phone, m.cursor == 0, w))
	b.WriteString("\n\n")
	b.WriteString(RenderInputField("📧 Email", m.address.Email, m.cursor == 1, w))
	b.WriteString("\n\n")
	b.WriteString(RenderInputField("🏠 Address", m.address.Address, m.cursor == 2, w))
	b.WriteString("\n\n")

	// Instructions
	b.WriteString(HelpStyle.Render("  Tab/↓: Next field • ↑: Previous • Enter: Continue"))
	b.WriteString("\n")

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderDivider(w))
	b.WriteString("\n")
	b.WriteString(RenderHelp("address", w))

	return b.String()
}

func (m *Model) renderCheckout(w int) string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("✓ Confirm Order", 0, w))
	b.WriteString("\n\n")

	// Order summary
	b.WriteString(SubtitleStyle.Width(w).Render("━━━ Order Summary ━━━"))
	b.WriteString("\n\n")

	for _, item := range m.cart.Items {
		b.WriteString(RenderOrderItem(types.OrderItem{
			Product:  item.Product,
			Quantity: item.Quantity,
		}, w))
		b.WriteString("\n")
	}

	// Total
	b.WriteString("\n")
	totalBox := BoxStyle.Width(w - 4).Render(fmt.Sprintf("💰 Total: %s", RenderPrice(m.cart.Total())))
	b.WriteString(totalBox)
	b.WriteString("\n\n")

	// Shipping details
	b.WriteString(SubtitleStyle.Width(w).Render("━━━ Shipping To ━━━"))
	b.WriteString("\n\n")
	b.WriteString(NormalStyle.Width(w).Render(fmt.Sprintf("  📱 %s", m.address.Phone)))
	b.WriteString("\n")
	b.WriteString(NormalStyle.Width(w).Render(fmt.Sprintf("  📧 %s", m.address.Email)))
	b.WriteString("\n")
	b.WriteString(NormalStyle.Width(w).Render(fmt.Sprintf("  🏠 %s", m.address.Address)))
	b.WriteString("\n\n")

	// Confirmation prompt
	b.WriteString(SuccessStyle.Render("  Press Enter or Y to place order"))
	b.WriteString("\n")

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderDivider(w))
	b.WriteString("\n")
	b.WriteString(RenderHelp("checkout", w))

	return b.String()
}

func (m *Model) renderOrderSuccess(w int) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(SuccessStyle.Render(SuccessArt))
	b.WriteString("\n\n")

	if m.order != nil {
		infoBox := BoxStyle.Width(w - 4).Render(fmt.Sprintf(
			"Order ID: %s\nTotal: %s\nStatus: %s",
			m.order.ID,
			RenderPrice(m.order.TotalAmount),
			string(m.order.Status.Type),
		))
		b.WriteString(infoBox)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Width(w).Render("  Thank you for your order! 🎉"))
	b.WriteString("\n")

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderDivider(w))
	b.WriteString("\n")
	b.WriteString(RenderHelp("order_success", w))

	return b.String()
}
