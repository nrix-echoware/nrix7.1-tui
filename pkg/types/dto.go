package types

type Media struct {
	URL      string `json:"url"`
	MimeType string `json:"mimetype"`
	Size     int64  `json:"size"`
}

type Discount struct {
	Rate float64      `json:"rate"`
	Type DiscountType `json:"type"`
}

type ProductVariantValue struct {
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

type ProductVariant struct {
	VariantName  string                `json:"variant_name"`
	VariantValues []ProductVariantValue `json:"variant_values"`
}

type CategoryDetail struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Discount    Discount `json:"discount"`
	Medias      []Media  `json:"medias"`
}

type Product struct {
	ID                string           `json:"_id"`
	Name              string           `json:"name"`
	Brand             string           `json:"brand"`
	Categories        []string         `json:"categories"`
	ProductDescription string          `json:"product_description"`
	MRPPrice          float64          `json:"mrp_price"`
	SellingPrice      float64          `json:"selling_price"`
	Tags              []string         `json:"tags"`
	Medias            []Media          `json:"medias"`
	Features          []string         `json:"features"`
	Active            bool             `json:"active"`
	ProductVariants   []ProductVariant `json:"product_variants"`
	CategoryDetails   []CategoryDetail  `json:"category_details"`
}

type Category struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Medias      []Media  `json:"medias"`
	Discount    Discount `json:"discount"`
}

type OrderItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type OrderStatusExtras struct {
	AgentPhone string `json:"agent_phone,omitempty"`
}

type OrderStatus struct {
	Type    OrderStatusType   `json:"type"`
	Reason  string            `json:"reason"`
	Extras  OrderStatusExtras `json:"extras"`
}

type ShippingDetails struct {
	ID           string `json:"_id,omitempty"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	IsDefault    bool   `json:"is_default"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	ClerkToken   string `json:"clerk_token,omitempty"`
	Address      string `json:"address,omitempty"`
}

type Order struct {
	ID              string          `json:"_id"`
	TotalAmount     float64         `json:"total_amount"`
	TotalDiscount   float64         `json:"total_discount"`
	OrderItems      []OrderItem     `json:"order_items"`
	ShippingDetails ShippingDetails `json:"shipping_details"`
	Status          OrderStatus     `json:"status"`
}

type APIRequest struct {
	Type      OperationType `json:"type"`
	Operation string        `json:"operation"`
	Params    interface{}   `json:"params"`
}

type APIResponse struct {
	Data  interface{} `json:"data"`
	Count int         `json:"count,omitempty"`
	Error string      `json:"error,omitempty"`
}

type ProductListParams struct {
	Skip              int    `json:"skip"`
	Take              int    `json:"take"`
	Active            *bool  `json:"active,omitempty"`
	CategoryID        string `json:"category_id,omitempty"`
	IncludeCategories bool   `json:"include_categories,omitempty"`
}

type ProductGetParams struct {
	ID string `json:"id"`
}

type ProductSearchParams struct {
	SearchTerm        string `json:"search_term"`
	Skip              int    `json:"skip"`
	Take              int    `json:"take"`
	IncludeCategories bool   `json:"include_categories,omitempty"`
}

type ProductCreateParams struct {
	Name              string           `json:"name"`
	Brand             string           `json:"brand"`
	Categories        []string         `json:"categories"`
	ProductDescription string          `json:"product_description"`
	MRPPrice          float64          `json:"mrp_price"`
	SellingPrice      float64          `json:"selling_price"`
	Tags              []string         `json:"tags"`
	Medias            []Media          `json:"medias"`
	Features          []string         `json:"features"`
	Active            bool             `json:"active"`
	ProductVariants   []ProductVariant `json:"product_variants"`
}

type ProductUpdateParams struct {
	ID                string           `json:"id"`
	Name              *string          `json:"name,omitempty"`
	Brand             *string          `json:"brand,omitempty"`
	Categories        []string         `json:"categories,omitempty"`
	ProductDescription *string         `json:"product_description,omitempty"`
	MRPPrice          *float64         `json:"mrp_price,omitempty"`
	SellingPrice      *float64         `json:"selling_price,omitempty"`
	Tags              []string         `json:"tags,omitempty"`
	Medias            []Media          `json:"medias,omitempty"`
	Features          []string         `json:"features,omitempty"`
	Active            *bool            `json:"active,omitempty"`
	ProductVariants   []ProductVariant `json:"product_variants,omitempty"`
}

type ProductDeleteParams struct {
	ID string `json:"id"`
}

type CategoryListParams struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type CategoryGetParams struct {
	ID string `json:"id"`
}

type CategoryCreateParams struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Medias      []Media  `json:"medias"`
	Discount    Discount `json:"discount"`
}

type CategoryUpdateParams struct {
	ID          string    `json:"id"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Medias      []Media   `json:"medias,omitempty"`
	Discount    *Discount `json:"discount,omitempty"`
}

type CategoryDeleteParams struct {
	ID string `json:"id"`
}

type OrderListParams struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type OrderGetParams struct {
	ID string `json:"id"`
}

type OrderCreateParams struct {
	ShippingAddress ShippingDetails      `json:"shippingAddress"`
	Items           []OrderItemInput     `json:"items"`
	SpecialMessage  string               `json:"specialMessage,omitempty"`
	Pricing         OrderPricingInput    `json:"pricing"`
	UserEmail       string               `json:"userEmail"`
	Timestamp       string               `json:"timestamp"`
	PaymentMethod   string               `json:"paymentMethod,omitempty"`
}

type OrderItemInput struct {
	ProductID   string            `json:"productId"`
	ProductName string            `json:"productName"`
	Variant     map[string]string `json:"variant"`
	Quantity    int               `json:"quantity"`
	Price       float64           `json:"price"`
	Total       float64           `json:"total"`
}

type OrderPricingInput struct {
	Subtotal float64 `json:"subtotal"`
	Discount float64 `json:"discount"`
	Shipping float64 `json:"shipping"`
	Total    float64 `json:"total"`
}

type OrderUpdateStatusParams struct {
	ID     string      `json:"id"`
	Status OrderStatus `json:"status"`
}

type CartItem struct {
	Product  Product
	Quantity int
	Variant  map[string]string
}

type Cart struct {
	Items []CartItem
}

func (c *Cart) Total() float64 {
	total := 0.0
	for _, item := range c.Items {
		total += item.Product.SellingPrice * float64(item.Quantity)
	}
	return total
}

func (c *Cart) Count() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

func (c *Cart) Add(product Product, quantity int, variant map[string]string) {
	for i := range c.Items {
		if c.Items[i].Product.ID == product.ID && variantsEqual(c.Items[i].Variant, variant) {
			c.Items[i].Quantity += quantity
			return
		}
	}
	c.Items = append(c.Items, CartItem{
		Product:  product,
		Quantity: quantity,
		Variant:  copyVariantMap(variant),
	})
}

func (c *Cart) Remove(productID string, variant map[string]string) {
	for i, item := range c.Items {
		if item.Product.ID == productID && variantsEqual(item.Variant, variant) {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}

func (c *Cart) UpdateQuantity(productID string, variant map[string]string, quantity int) {
	for i := range c.Items {
		if c.Items[i].Product.ID == productID && variantsEqual(c.Items[i].Variant, variant) {
			if quantity <= 0 {
				c.Remove(productID, variant)
			} else {
				c.Items[i].Quantity = quantity
			}
			return
		}
	}
}

func variantsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func copyVariantMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return map[string]string{}
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
