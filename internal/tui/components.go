package tui

import (
	"fmt"
	"strings"
	"terminal-echoware/pkg/config"
	"terminal-echoware/pkg/types"

	"github.com/charmbracelet/lipgloss"
)

// Column widths for alignment
const (
	ColCursor  = 3
	ColName    = 40
	ColBrand   = 15
	ColPrice   = 12
	ColQty     = 12
)

func RenderProductLine(product types.Product, selected bool, width int) string {
	style := ProductCardStyle
	if selected {
		style = ProductCardSelectedStyle
	}

	cursor := "  "
	if selected {
		cursor = "▸ "
	}

	// Calculate dynamic name width based on available space
	nameWidth := width - ColCursor - ColBrand - ColPrice - 6
	if nameWidth < 20 {
		nameWidth = 20
	}
	if nameWidth > 50 {
		nameWidth = 50
	}

	name := padRight(truncate(product.Name, nameWidth), nameWidth)
	brand := padRight(truncate(product.Brand, ColBrand-2), ColBrand-2)
	price := fmt.Sprintf("₹%.0f", product.SellingPrice)

	line := fmt.Sprintf("%s%s  %s  %s", 
		cursor,
		NormalStyle.Render(name),
		BrandStyle.Render(brand),
		PriceStyle.Render(price),
	)
	return style.Width(width).Render(line)
}

func RenderProductList(products []types.Product, cursor int, width int) string {
	var b strings.Builder
	for i, product := range products {
		b.WriteString(RenderProductLine(product, i == cursor, width))
		b.WriteString("\n")
	}
	return b.String()
}

func RenderCartItem(item types.CartItem, selected bool, width int) string {
	style := NormalStyle
	cursor := "  "
	if selected {
		style = SelectedStyle
		cursor = "▸ "
	}
	
	nameWidth := width - ColCursor - ColQty - ColPrice - 8
	if nameWidth < 15 {
		nameWidth = 15
	}
	
	total := item.Product.SellingPrice * float64(item.Quantity)
	name := padRight(truncate(item.Product.Name, nameWidth), nameWidth)
	qty := fmt.Sprintf("x%d", item.Quantity)
	price := fmt.Sprintf("₹%.0f", total)
	
	line := fmt.Sprintf("%s%s  %s  %s", cursor, name, padLeft(qty, 4), padLeft(price, 10))
	return style.Width(width).Render(line)
}

func RenderCartItemWithQty(item types.CartItem, selected bool, width int) string {
	style := NormalStyle
	cursor := "  "
	if selected {
		style = SelectedStyle
		cursor = "▸ "
	}
	
	nameWidth := width - ColCursor - 20 - ColPrice - 6
	if nameWidth < 15 {
		nameWidth = 15
	}
	
	total := item.Product.SellingPrice * float64(item.Quantity)
	name := padRight(truncate(item.Product.Name, nameWidth), nameWidth)
	
	qtyStyle := HelpStyle
	if selected {
		qtyStyle = SuccessStyle
	}
	qtyStr := qtyStyle.Render(fmt.Sprintf("[ - ] %2d [ + ]", item.Quantity))
	price := PriceStyle.Render(fmt.Sprintf("₹%.0f", total))
	
	line := fmt.Sprintf("%s%s  %s  %s", cursor, name, qtyStr, price)
	return style.Width(width).Render(line)
}

func RenderCartList(items []types.CartItem, cursor int, width int) string {
	var b strings.Builder
	for i, item := range items {
		b.WriteString(RenderCartItem(item, i == cursor, width))
		b.WriteString("\n")
	}
	return b.String()
}

func RenderCartListWithQty(items []types.CartItem, cursor int, width int) string {
	var b strings.Builder
	for i, item := range items {
		b.WriteString(RenderCartItemWithQty(item, i == cursor, width))
		b.WriteString("\n")
	}
	return b.String()
}

func RenderHelp(screenName string, width int) string {
	cfg := config.GetConfig()
	if !cfg.ShowControls {
		return ""
	}
	return FooterStyle.Width(width).Render(HelpStyle.Render(cfg.GetHelpText(screenName)))
}

func RenderPrice(amount float64) string {
	return PriceStyle.Render(fmt.Sprintf("₹%.0f", amount))
}

func RenderInputField(label, value string, focused bool, width int) string {
	style := InputStyle.Width(width - 4)
	if focused {
		style = InputFocusedStyle.Width(width - 4)
	}
	
	cursor := ""
	if focused {
		cursor = "▌"
	}
	
	content := fmt.Sprintf("%s: %s%s", label, value, cursor)
	return style.Render(content)
}

func RenderOrderItem(item types.OrderItem, width int) string {
	nameWidth := width - 20
	if nameWidth < 15 {
		nameWidth = 15
	}
	
	total := item.Product.SellingPrice * float64(item.Quantity)
	name := padRight(truncate(item.Product.Name, nameWidth), nameWidth)
	return fmt.Sprintf("  %s  x%d  %s", name, item.Quantity, RenderPrice(total))
}

func RenderNotification(notif *Notification) string {
	if notif == nil {
		return ""
	}
	
	var style lipgloss.Style
	switch notif.Type {
	case "success":
		style = NotificationSuccessStyle
	case "error":
		style = NotificationErrorStyle
	default:
		style = NotificationInfoStyle
	}
	
	return style.Render(" " + notif.Message + " ")
}

func RenderHeader(title string, cartCount int, width int) string {
	titleStr := TitleStyle.Render(title)
	
	if cartCount > 0 {
		badge := CartBadgeStyle.Render(fmt.Sprintf(" 🛒 %d ", cartCount))
		// Calculate spacing
		titleLen := lipgloss.Width(titleStr)
		badgeLen := lipgloss.Width(badge)
		spacing := width - titleLen - badgeLen - 2
		if spacing < 2 {
			spacing = 2
		}
		return titleStr + strings.Repeat(" ", spacing) + badge
	}
	
	return titleStr
}

func RenderDivider(width int) string {
	return DividerStyle.Render(strings.Repeat("─", width))
}

func RenderQuantitySelector(quantity int, focused bool) string {
	style := HelpStyle
	if focused {
		style = SuccessStyle
	}
	return style.Render(fmt.Sprintf("[ - ]  %d  [ + ]", quantity))
}

// RenderOptionRow renders a single option row (like quantity)
func RenderOptionRow(label string, value string, focused bool, width int) string {
	rowStyle := OptionRowStyle.Width(width)
	if focused {
		rowStyle = OptionRowFocusedStyle.Width(width)
	}
	
	labelStr := OptionLabelStyle.Render(padRight(label+":", 12))
	
	var valueStr string
	if focused {
		valueStr = fmt.Sprintf("◀  %s  ▶", OptionValueSelectedStyle.Render(value))
	} else {
		valueStr = fmt.Sprintf("   %s   ", OptionValueStyle.Render(value))
	}
	
	return rowStyle.Render(fmt.Sprintf("  %s %s", labelStr, valueStr))
}

// RenderVariantRow renders a variant row with multiple options
func RenderVariantRow(variantName string, options []string, selectedIdx int, focused bool, width int) string {
	rowStyle := OptionRowStyle.Width(width)
	if focused {
		rowStyle = OptionRowFocusedStyle.Width(width)
	}
	
	labelStr := OptionLabelStyle.Render(padRight(variantName+":", 12))
	
	var optionsStr strings.Builder
	if focused {
		optionsStr.WriteString("◀  ")
	} else {
		optionsStr.WriteString("   ")
	}
	
	for i, opt := range options {
		if i == selectedIdx {
			optionsStr.WriteString(OptionValueSelectedStyle.Render(opt))
		} else {
			optionsStr.WriteString(OptionValueUnselectedStyle.Render(opt))
		}
		if i < len(options)-1 {
			optionsStr.WriteString(" ")
		}
	}
	
	if focused {
		optionsStr.WriteString("  ▶")
	}
	
	return rowStyle.Render(fmt.Sprintf("  %s %s", labelStr, optionsStr.String()))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

func padLeft(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(" ", length-len(s)) + s
}
