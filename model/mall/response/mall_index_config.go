package response

import "github.com/shopspring/decimal"

type MallIndexConfigGoodsResponse struct {
	GoodsId       int             `json:"goodsId"`
	GoodsName     string          `json:"goodsName"`
	GoodsIntro    string          `json:"goodsIntro"`
	GoodsCoverImg string          `json:"goodsCoverImg"`
	SellingPrice  decimal.Decimal `json:"sellingPrice"`
	Tag           string          `json:"tag"`
}
