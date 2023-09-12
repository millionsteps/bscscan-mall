package response

import "github.com/shopspring/decimal"

type GoodsSearchResponse struct {
	GoodsId       int             `json:"goodsId"`
	GoodsName     string          `json:"goodsName"`
	GoodsIntro    string          `json:"goodsIntro"`
	GoodsCoverImg string          `json:"goodsCoverImg"`
	SellingPrice  decimal.Decimal `json:"sellingPrice"`
}

type GoodsInfoDetailResponse struct {
	GoodsId            int             `json:"goodsId"`
	GoodsName          string          `json:"goodsName"`
	GoodsIntro         string          `json:"goodsIntro"`
	GoodsCoverImg      string          `json:"goodsCoverImg"`
	SellingPrice       decimal.Decimal `json:"sellingPrice"`
	GoodsDetailContent string          `json:"goodsDetailContent"  `
	OriginalPrice      decimal.Decimal `json:"originalPrice" `
	Tag                string          `json:"tag" form:"tag" `
	GoodsCarouselList  []string        `json:"goodsCarouselList" `
}
