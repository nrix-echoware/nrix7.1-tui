package tui

import (
	"fmt"
	"strings"
	"terminal-echoware/pkg/types"
)

func (m *Model) View() string {
	if m.loading {
		return m.renderLoading()
	}

	var content string
	switch m.screen {
	case types.ScreenHome:
		content = m.renderHome()
	case types.ScreenSearch:
		content = m.renderSearch()
	case types.ScreenProduct:
		content = m.renderProduct()
	case types.ScreenCart:
		content = m.renderCart()
	case types.ScreenAddress:
		content = m.renderAddress()
	case types.ScreenCheckout:
		content = m.renderCheckout()
	case types.ScreenOrderSuccess:
		content = m.renderOrderSuccess()
	}

	// Add notification at the bottom if present
	if m.notification != nil {
		content += "\n" + RenderNotification(m.notification)
	}

	// Add error if present
	if m.err != nil {
		content += "\n" + ErrorStyle.Render(fmt.Sprintf("⚠ Error: %v", m.err))
	}

	return content
}

func (m *Model) renderLoading() string {
	var b strings.Builder
	
	b.WriteString("\n")
	b.WriteString(TitleStyle.Render(AsciiLogo))
	b.WriteString("\n\n")
	
	frame := string(LoadingFrames[m.loadingFrame%len(LoadingFrames)])
	b.WriteString(LoadingStyle.Render(fmt.Sprintf("  %s %s", frame, m.loadingMsg)))
	
	return b.String()
}

func (m *Model) renderHome() string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(TitleStyle.Render(AsciiLogo))
	b.WriteString("\n")
	
	// Cart badge
	if m.cart.Count() > 0 {
		b.WriteString(CartBadgeStyle.Render(fmt.Sprintf(" 🛒 Cart: %d items ", m.cart.Count())))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Products section
	b.WriteString(SubtitleStyle.Render("━━━ Products ━━━"))
	b.WriteString("\n\n")

	if len(m.homeProducts) == 0 {
		b.WriteString(NormalStyle.Render("  No products available."))
		b.WriteString("\n")
	} else {
		b.WriteString(RenderProductList(m.homeProducts, m.cursor))
	}

	// Footer with controls
	b.WriteString("\n")
	b.WriteString(RenderHelp("home"))

	return b.String()
}

func (m *Model) renderSearch() string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("🔍 Search Products", m.cart.Count()))
	b.WriteString("\n\n")

	// Search input
	b.WriteString(RenderInputField("Search", m.searchQuery, true))
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
		b.WriteString(RenderProductList(m.searchResults, m.cursor))
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderHelp("search"))

	return b.String()
}

func (m *Model) renderProduct() string {
	if m.currentProduct == nil {
		return ErrorStyle.Render("Product not found.")
	}

	var b strings.Builder
	p := m.currentProduct

	// Header with full product name
	b.WriteString("\n")
	b.WriteString(TitleStyle.Render(p.Name))
	b.WriteString("\n")
	
	// Cart badge
	if m.cart.Count() > 0 {
		b.WriteString(CartBadgeStyle.Render(fmt.Sprintf(" 🛒 %d ", m.cart.Count())))
	}
	b.WriteString("\n\n")

	// Price and brand section
	priceBox := BoxStyle.Render(fmt.Sprintf(
		"💰 Price: %s\n🏷️  Brand: %s",
		RenderPrice(p.SellingPrice),
		BrandStyle.Render(p.Brand),
	))
	b.WriteString(priceBox)
	b.WriteString("\n")

	if p.MRPPrice > p.SellingPrice {
		discount := ((p.MRPPrice - p.SellingPrice) / p.MRPPrice) * 100
		b.WriteString(SuccessStyle.Render(fmt.Sprintf("  🎉 %.0f%% OFF (MRP: ₹%.0f)", discount, p.MRPPrice)))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Quantity selector (focus index 0)
	b.WriteString(SubtitleStyle.Render("📦 Options"))
	b.WriteString("\n\n")
	
	qtyFocused := m.variantFocusIndex == 0
	b.WriteString(RenderOptionRow("Quantity", fmt.Sprintf("%d", m.productQuantity), qtyFocused))
	b.WriteString("\n")

	// Variant selectors
	for i, variant := range p.ProductVariants {
		focused := m.variantFocusIndex == i+1
		selectedValue := ""
		if i < len(m.variantSelections) {
			selIdx := m.variantSelections[i].SelectedIndex
			if selIdx < len(variant.VariantValues) {
				selectedValue = variant.VariantValues[selIdx].Label
			}
		}
		
		// Show all options with selected highlighted
		options := []string{}
		for j, val := range variant.VariantValues {
			if j == m.variantSelections[i].SelectedIndex {
				options = append(options, fmt.Sprintf("[%s]", val.Label))
			} else {
				options = append(options, val.Label)
			}
		}
		
		b.WriteString(RenderVariantRow(variant.VariantName, options, m.variantSelections[i].SelectedIndex, focused))
		b.WriteString("\n")
		_ = selectedValue // not used directly, but keeping for future use
	}
	b.WriteString("\n")

	// Description
	if p.ProductDescription != "" {
		b.WriteString(SubtitleStyle.Render("📝 Description"))
		b.WriteString("\n")
		desc := p.ProductDescription
		if len(desc) > 100 {
			desc = desc[:100] + "..."
		}
		b.WriteString(NormalStyle.Render("  " + desc))
		b.WriteString("\n\n")
	}

	// Features
	if len(p.Features) > 0 {
		b.WriteString(SubtitleStyle.Render("✨ Features"))
		b.WriteString("\n")
		for i, feature := range p.Features {
			if i >= 4 {
				b.WriteString(HelpStyle.Render(fmt.Sprintf("  ... and %d more", len(p.Features)-4)))
				b.WriteString("\n")
				break
			}
			b.WriteString(NormalStyle.Render(fmt.Sprintf("  • %s", feature)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Tags
	if len(p.Tags) > 0 {
		b.WriteString(HelpStyle.Render("🏷️  Tags: " + strings.Join(p.Tags, ", ")))
		b.WriteString("\n")
	}

	// Footer with controls
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render(strings.Repeat("─", 60)))
	b.WriteString("\n")
	b.WriteString(FooterStyle.Render(HelpStyle.Render("Tab/↑↓: Navigate • ←→: Change • A/Enter: Add • C: Cart • Esc: Back")))

	return b.String()
}

func (m *Model) renderCart() string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("🛒 Shopping Cart", 0))
	b.WriteString("\n\n")

	if len(m.cart.Items) == 0 {
		b.WriteString(NormalStyle.Render("  Your cart is empty."))
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("  Press 'b' to browse products"))
		b.WriteString("\n")
	} else {
		// Items header
		b.WriteString(SubtitleStyle.Render("━━━ Items ━━━"))
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("  Use +/- to change quantity"))
		b.WriteString("\n\n")
		
		// Items list
		b.WriteString(RenderCartListWithQty(m.cart.Items, m.cursor))
		
		// Total
		b.WriteString("\n")
		totalBox := BoxStyle.Render(fmt.Sprintf("💰 Total: %s", RenderPrice(m.cart.Total())))
		b.WriteString(totalBox)
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(FooterStyle.Render(HelpStyle.Render("↑↓: Navigate • +/-: Qty • D: Remove • Enter: Checkout • Esc: Back")))

	return b.String()
}

func (m *Model) renderAddress() string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("📦 Shipping Details", m.cart.Count()))
	b.WriteString("\n\n")

	// Form fields
	b.WriteString(RenderInputField("📱 Phone", m.address.Phone, m.cursor == 0))
	b.WriteString("\n\n")
	b.WriteString(RenderInputField("📧 Email", m.address.Email, m.cursor == 1))
	b.WriteString("\n\n")
	b.WriteString(RenderInputField("🏠 Address", m.address.Address, m.cursor == 2))
	b.WriteString("\n\n")

	// Instructions
	b.WriteString(HelpStyle.Render("  Tab/↓: Next field • ↑: Previous • Enter: Continue"))
	b.WriteString("\n")

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderHelp("address"))

	return b.String()
}

func (m *Model) renderCheckout() string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderHeader("✓ Confirm Order", 0))
	b.WriteString("\n\n")

	// Order summary
	b.WriteString(SubtitleStyle.Render("━━━ Order Summary ━━━"))
	b.WriteString("\n\n")

	for _, item := range m.cart.Items {
		b.WriteString(RenderOrderItem(types.OrderItem{
			Product:  item.Product,
			Quantity: item.Quantity,
		}))
		b.WriteString("\n")
	}

	// Total
	b.WriteString("\n")
	totalBox := BoxStyle.Render(fmt.Sprintf("💰 Total: %s", RenderPrice(m.cart.Total())))
	b.WriteString(totalBox)
	b.WriteString("\n\n")

	// Shipping details
	b.WriteString(SubtitleStyle.Render("━━━ Shipping To ━━━"))
	b.WriteString("\n\n")
	b.WriteString(NormalStyle.Render(fmt.Sprintf("  📱 %s", m.address.Phone)))
	b.WriteString("\n")
	b.WriteString(NormalStyle.Render(fmt.Sprintf("  📧 %s", m.address.Email)))
	b.WriteString("\n")
	b.WriteString(NormalStyle.Render(fmt.Sprintf("  🏠 %s", m.address.Address)))
	b.WriteString("\n\n")

	// Confirmation prompt
	b.WriteString(SuccessStyle.Render("  Press Enter or Y to place order"))
	b.WriteString("\n")

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderHelp("checkout"))

	return b.String()
}

func (m *Model) renderOrderSuccess() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(SuccessStyle.Render(SuccessArt))
	b.WriteString("\n\n")

	if m.order != nil {
		infoBox := BoxStyle.Render(fmt.Sprintf(
			"Order ID: %s\nTotal: %s\nStatus: %s",
			m.order.ID,
			RenderPrice(m.order.TotalAmount),
			string(m.order.Status.Type),
		))
		b.WriteString(infoBox)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("  Thank you for your order! 🎉"))
	b.WriteString("\n")

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderHelp("order_success"))

	return b.String()
}
