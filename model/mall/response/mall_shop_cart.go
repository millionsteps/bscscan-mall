package response

import "github.com/shopspring/decimal"

type CartItemResponse struct {
	CartItemId int `json:"cartItemId"`

	GoodsId int `json:"goodsId"`

	GoodsCount int `json:"goodsCount"`

	GoodsName string `json:"goodsName"`

	GoodsCoverImg string `json:"goodsCoverImg"`

	SellingPrice decimal.Decimal `json:"sellingPrice"`
}
