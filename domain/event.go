package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderType int8

const (
	SpotOrderType_Default OrderType = 0
	SpotOrderType_Market  OrderType = 1
	SpotOrderType_Limit   OrderType = 2
)

type OrderState uint32

const (
	OrderStateDefault  OrderState = 0
	OrderStateCreated  OrderState = 1
	OrderStateOpen     OrderState = 2
	OrderStateCanceled OrderState = 3
	OrderStateFilled   OrderState = 4
)

type Side int8

const (
	Side_Default Side = 0
	Side_Buy     Side = 1
	Side_Sell    Side = 2
)

type TimeInForce int8

const (
	TIF_Default TimeInForce = 0
	TIF_GTC     TimeInForce = 1
)

type CreatedOrderEvent struct {
	ID            string
	Sequence      uint64
	ClientOrderID string
	Market        string
	Base          string
	Quote         string
	Type          OrderType
	Price         decimal.Decimal
	Size          decimal.Decimal
	Side          Side
	TimeInForce   TimeInForce
	TakerFeeRate  decimal.Decimal
	MakerFeeRate  decimal.Decimal
	UserID        uint64
	Source        string
	CreatedAt     time.Time
}