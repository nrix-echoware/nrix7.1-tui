package types

type OperationType string

const (
	OperationTypeQuery    OperationType = "query"
	OperationTypeMutation OperationType = "mutation"
)

type OrderStatusType string

const (
	OrderStatusAccepted        OrderStatusType = "accepted"
	OrderStatusRejected        OrderStatusType = "rejected"
	OrderStatusRejectedByUser  OrderStatusType = "rejected_by_user"
	OrderStatusDelivered       OrderStatusType = "delivered"
	OrderStatusOutForDelivery  OrderStatusType = "out_for_delivery"
	OrderStatusAgent           OrderStatusType = "agent"
	OrderStatusAgentChanged    OrderStatusType = "agent_changed"
	OrderStatusInHub           OrderStatusType = "in_hub"
)

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeDirect     DiscountType = "direct"
)

type Screen int

const (
	ScreenHome Screen = iota
	ScreenSearch
	ScreenProduct
	ScreenCart
	ScreenAddress
	ScreenCheckout
	ScreenOrderSuccess
)
