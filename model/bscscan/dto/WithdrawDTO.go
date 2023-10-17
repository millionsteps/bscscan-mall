package dto

import "github.com/shopspring/decimal"

type WithdrawDTO struct {
	Usdt    decimal.Decimal `json:"usdt"`
	Address string          `json:"address"`
}
