package config

import "strings"

type KeyBinding struct {
	Key         string
	Description string
}

type ControlsConfig struct {
	ShowHelp      bool
	HelpPosition  string
	KeyBindings   map[string]KeyBinding
	CustomBindings map[string]string
}

type AppConfig struct {
	APIBaseURL         string
	SSHPort            string
	ShowControls       bool
	ShopName           string
	CompanyName        string
	CompanyDescription string
	Controls           ControlsConfig
	Theme              ThemeConfig
}

type ThemeConfig struct {
	PrimaryColor   string
	SecondaryColor string
	AccentColor    string
	ErrorColor     string
	SuccessColor   string
}

var GlobalConfig *AppConfig

func InitConfig() {
	GlobalConfig = &AppConfig{
		APIBaseURL:         "https://lowkey-backend-omega.vercel.app",
		SSHPort:            "2222",
		ShowControls:       true,
		ShopName:           "Nrix7 Shop",
		CompanyName:        "Nrix7 E-Commerce",
		CompanyDescription: "Welcome to Nrix7 - Your trusted destination for quality products at unbeatable prices. We bring you the latest trends in fashion, electronics, home essentials, and more. Shop with confidence through our secure terminal interface. Fast shipping, easy returns, and 24/7 customer support.",
		Controls: ControlsConfig{
			ShowHelp:     true,
			HelpPosition: "bottom",
			KeyBindings: map[string]KeyBinding{
				"navigate_up":    {Key: "↑/k", Description: "Navigate up"},
				"navigate_down":  {Key: "↓/j", Description: "Navigate down"},
				"select":         {Key: "Enter", Description: "Select/Confirm"},
				"search":         {Key: "S", Description: "Search"},
				"cart":           {Key: "C", Description: "View Cart"},
				"add_to_cart":    {Key: "A", Description: "Add to Cart"},
				"delete":         {Key: "D", Description: "Delete"},
				"back":            {Key: "Esc/B", Description: "Back"},
				"quit":            {Key: "Q/Ctrl+C", Description: "Quit"},
				"tab":             {Key: "Tab", Description: "Switch field"},
			},
			CustomBindings: make(map[string]string),
		},
		Theme: ThemeConfig{
			PrimaryColor:   "205",
			SecondaryColor: "117",
			AccentColor:    "82",
			ErrorColor:     "196",
			SuccessColor:   "46",
		},
	}
}

func GetConfig() *AppConfig {
	if GlobalConfig == nil {
		InitConfig()
	}
	return GlobalConfig
}

func (c *AppConfig) GetKeyBinding(action string) KeyBinding {
	if binding, ok := c.Controls.KeyBindings[action]; ok {
		return binding
	}
	return KeyBinding{Key: "", Description: ""}
}

func (c *AppConfig) GetHelpText(screen string) string {
	if !c.ShowControls || !c.Controls.ShowHelp {
		return ""
	}

	var bindings []string
	switch screen {
	case "home":
		bindings = []string{
			c.Controls.KeyBindings["navigate_up"].Key + ": " + c.Controls.KeyBindings["navigate_up"].Description,
			c.Controls.KeyBindings["navigate_down"].Key + ": " + c.Controls.KeyBindings["navigate_down"].Description,
			c.Controls.KeyBindings["select"].Key + ": View",
			c.Controls.KeyBindings["search"].Key + ": " + c.Controls.KeyBindings["search"].Description,
			c.Controls.KeyBindings["cart"].Key + ": " + c.Controls.KeyBindings["cart"].Description,
			c.Controls.KeyBindings["quit"].Key + ": " + c.Controls.KeyBindings["quit"].Description,
		}
	case "search":
		bindings = []string{
			"Type: Search",
			c.Controls.KeyBindings["navigate_up"].Key + ": " + c.Controls.KeyBindings["navigate_up"].Description,
			c.Controls.KeyBindings["navigate_down"].Key + ": " + c.Controls.KeyBindings["navigate_down"].Description,
			c.Controls.KeyBindings["select"].Key + ": Search",
			c.Controls.KeyBindings["back"].Key + ": " + c.Controls.KeyBindings["back"].Description,
		}
	case "product":
		bindings = []string{
			c.Controls.KeyBindings["add_to_cart"].Key + ": " + c.Controls.KeyBindings["add_to_cart"].Description,
			c.Controls.KeyBindings["cart"].Key + ": " + c.Controls.KeyBindings["cart"].Description,
			c.Controls.KeyBindings["back"].Key + ": " + c.Controls.KeyBindings["back"].Description,
			c.Controls.KeyBindings["quit"].Key + ": " + c.Controls.KeyBindings["quit"].Description,
		}
	case "cart":
		bindings = []string{
			c.Controls.KeyBindings["navigate_up"].Key + ": " + c.Controls.KeyBindings["navigate_up"].Description,
			c.Controls.KeyBindings["navigate_down"].Key + ": " + c.Controls.KeyBindings["navigate_down"].Description,
			c.Controls.KeyBindings["delete"].Key + ": " + c.Controls.KeyBindings["delete"].Description,
			c.Controls.KeyBindings["select"].Key + ": Checkout",
			c.Controls.KeyBindings["back"].Key + ": " + c.Controls.KeyBindings["back"].Description,
		}
	case "address":
		bindings = []string{
			c.Controls.KeyBindings["tab"].Key + ": " + c.Controls.KeyBindings["tab"].Description,
			c.Controls.KeyBindings["select"].Key + ": Continue",
			c.Controls.KeyBindings["back"].Key + ": " + c.Controls.KeyBindings["back"].Description,
		}
	case "checkout":
		bindings = []string{
			c.Controls.KeyBindings["select"].Key + "/Y: Place Order",
			"N: Back",
			c.Controls.KeyBindings["back"].Key + ": " + c.Controls.KeyBindings["back"].Description,
		}
	case "order_success":
		bindings = []string{
			c.Controls.KeyBindings["select"].Key + ": Continue Shopping",
		}
	default:
		bindings = []string{
			c.Controls.KeyBindings["quit"].Key + ": " + c.Controls.KeyBindings["quit"].Description,
		}
	}

	return strings.Join(bindings, " | ")
}
