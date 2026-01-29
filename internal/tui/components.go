package tui

import (
	"fmt"
	"strings"
	"terminal-echoware/pkg/config"
	"terminal-echoware/pkg/types"

	"github.com/charmbracelet/lipgloss"
)

func RenderProductLine(product types.Product, selected bool) string {
	style := ProductCardStyle
	if selected {
		style = ProductCardSelectedStyle
	}

	cursor := "  "
	if selected {
		cursor = "▸ "
	}

	price := PriceStyle.Render(fmt.Sprintf("₹%.0f", product.SellingPrice))
	brand := BrandStyle.Render(product.Brand)

	line := fmt.Sprintf("%s%-30s %s  %s", cursor, truncate(product.Name, 28), brand, price)
	return style.Render(line)
}

func RenderProductList(products []types.Product, cursor int) string {
	var b strings.Builder
	for i, product := range products {
		b.WriteString(RenderProductLine(product, i == cursor))
		b.WriteString("\n")
	}
	return b.String()
}

func RenderCartItem(item types.CartItem, selected bool) string {
	style := NormalStyle
	cursor := "  "
	if selected {
		style = SelectedStyle
		cursor = "▸ "
	}
	
	total := item.Product.SellingPrice * float64(item.Quantity)
	line := fmt.Sprintf("%s%-25s  x%d  %s",
		cursor,
		truncate(item.Product.Name, 23),
		item.Quantity,
		PriceStyle.Render(fmt.Sprintf("₹%.0f", total)),
	)
	return style.Render(line)
}

func RenderCartItemWithQty(item types.CartItem, selected bool) string {
	style := NormalStyle
	cursor := "  "
	if selected {
		style = SelectedStyle
		cursor = "▸ "
	}
	
	total := item.Product.SellingPrice * float64(item.Quantity)
	qtyStyle := HelpStyle
	if selected {
		qtyStyle = SuccessStyle
	}
	
	line := fmt.Sprintf("%s%-22s  %s  %s",
		cursor,
		truncate(item.Product.Name, 20),
		qtyStyle.Render(fmt.Sprintf("[-] %d [+]", item.Quantity)),
		PriceStyle.Render(fmt.Sprintf("₹%.0f", total)),
	)
	return style.Render(line)
}

func RenderCartList(items []types.CartItem, cursor int) string {
	var b strings.Builder
	for i, item := range items {
		b.WriteString(RenderCartItem(item, i == cursor))
		b.WriteString("\n")
	}
	return b.String()
}

func RenderCartListWithQty(items []types.CartItem, cursor int) string {
	var b strings.Builder
	for i, item := range items {
		b.WriteString(RenderCartItemWithQty(item, i == cursor))
		b.WriteString("\n")
	}
	return b.String()
}

func RenderHelp(screenName string) string {
	cfg := config.GetConfig()
	if !cfg.ShowControls {
		return ""
	}
	return FooterStyle.Render(HelpStyle.Render(cfg.GetHelpText(screenName)))
}

func RenderPrice(amount float64) string {
	return PriceStyle.Render(fmt.Sprintf("₹%.0f", amount))
}

func RenderInputField(label, value string, focused bool) string {
	style := InputStyle
	if focused {
		style = InputFocusedStyle
	}
	
	cursor := ""
	if focused {
		cursor = "▌"
	}
	
	content := fmt.Sprintf("%s: %s%s", label, value, cursor)
	return style.Render(content)
}

func RenderOrderItem(item types.OrderItem) string {
	total := item.Product.SellingPrice * float64(item.Quantity)
	return fmt.Sprintf("  %-22s  x%d  %s",
		truncate(item.Product.Name, 20),
		item.Quantity,
		RenderPrice(total),
	)
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

func RenderHeader(title string, cartCount int) string {
	var b strings.Builder
	
	b.WriteString(TitleStyle.Render(title))
	
	if cartCount > 0 {
		b.WriteString("  ")
		b.WriteString(CartBadgeStyle.Render(fmt.Sprintf(" 🛒 %d ", cartCount)))
	}
	
	return b.String()
}

func RenderDivider(width int) string {
	return HelpStyle.Render(strings.Repeat("─", width))
}

func RenderQuantitySelector(quantity int, focused bool) string {
	style := HelpStyle
	if focused {
		style = SuccessStyle
	}
	return style.Render(fmt.Sprintf("[ - ]  %d  [ + ]", quantity))
}

// RenderOptionRow renders a single option row (like quantity)
func RenderOptionRow(label string, value string, focused bool) string {
	rowStyle := OptionRowStyle
	if focused {
		rowStyle = OptionRowFocusedStyle
	}
	
	labelStr := OptionLabelStyle.Render(label + ":")
	
	var valueStr string
	if focused {
		valueStr = fmt.Sprintf("  ◀  %s  ▶", OptionValueSelectedStyle.Render(value))
	} else {
		valueStr = fmt.Sprintf("     %s   ", OptionValueStyle.Render(value))
	}
	
	return rowStyle.Render(fmt.Sprintf("  %s %s", labelStr, valueStr))
}

// RenderVariantRow renders a variant row with multiple options
func RenderVariantRow(variantName string, options []string, selectedIdx int, focused bool) string {
	rowStyle := OptionRowStyle
	if focused {
		rowStyle = OptionRowFocusedStyle
	}
	
	labelStr := OptionLabelStyle.Render(variantName + ":")
	
	var optionsStr strings.Builder
	if focused {
		optionsStr.WriteString("  ◀  ")
	} else {
		optionsStr.WriteString("     ")
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
	return s[:maxLen-3] + "..."
}
