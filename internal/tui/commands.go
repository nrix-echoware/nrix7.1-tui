package tui

import (
	"terminal-echoware/internal/api"
	"terminal-echoware/pkg/types"

	tea "github.com/charmbracelet/bubbletea"
)

type productsLoadedMsg struct {
	products []types.Product
	count    int
	err      error
}

type productLoadedMsg struct {
	product *types.Product
	err     error
}

type searchResultsMsg struct {
	products []types.Product
	count    int
	err      error
}

type orderCreatedMsg struct {
	order *types.Order
	err   error
}

func loadProductsCmd(client *api.Client, skip, take int) tea.Cmd {
	return func() tea.Msg {
		active := true
		products, count, err := client.ListProducts(types.ProductListParams{
			Skip:              skip,
			Take:              take,
			Active:            &active,
			IncludeCategories: true,
		})
		return productsLoadedMsg{products: products, count: count, err: err}
	}
}

func loadProductCmd(client *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		product, err := client.GetProduct(id)
		return productLoadedMsg{product: product, err: err}
	}
}

func searchProductsCmd(client *api.Client, query string, skip, take int) tea.Cmd {
	return func() tea.Msg {
		products, count, err := client.SearchProducts(types.ProductSearchParams{
			SearchTerm:        query,
			Skip:              skip,
			Take:              take,
			IncludeCategories: true,
		})
		return searchResultsMsg{products: products, count: count, err: err}
	}
}

func createOrderCmd(client *api.Client, params types.OrderCreateParams) tea.Cmd {
	return func() tea.Msg {
		order, err := client.CreateOrder(params)
		return orderCreatedMsg{order: order, err: err}
	}
}
