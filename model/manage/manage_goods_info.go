// 自动生成模板MallGoodsInfo
package manage

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

// MallGoodsInfo 结构体
// 如果含有time.Time 请自行import time包
type MallGoodsInfo struct {
	GoodsId            int             `json:"goodsId" form:"goodsId" gorm:"primarykey;AUTO_INCREMENT"`
	GoodsName          string          `json:"goodsName" form:"goodsName" gorm:"column:goods_name;comment:商品名;type:varchar(200);"`
	GoodsIntro         string          `json:"goodsIntro" form:"goodsIntro" gorm:"column:goods_intro;comment:商品简介;type:varchar(200);"`
	GoodsCategoryId    int             `json:"goodsCategoryId" form:"goodsCategoryId" gorm:"column:goods_category_id;comment:关联分类id;type:bigint"`
	GoodsCoverImg      string          `json:"goodsCoverImg" form:"goodsCoverImg" gorm:"column:goods_cover_img;comment:商品主图;type:varchar(200);"`
	GoodsCarousel      string          `json:"goodsCarousel" form:"goodsCarousel" gorm:"column:goods_carousel;comment:商品轮播图;type:varchar(500);"`
	GoodsDetailContent string          `json:"goodsDetailContent" form:"goodsDetailContent" gorm:"column:goods_detail_content;comment:商品详情;type:text;"`
	OriginalPrice      decimal.Decimal `json:"originalPrice" form:"originalPrice" gorm:"column:original_price;comment:商品价格;type:decimal"`
	SellingPrice       decimal.Decimal `json:"sellingPrice" form:"sellingPrice" gorm:"column:selling_price;comment:商品实际售价;type:decimal"`
	StockNum           int             `json:"stockNum" form:"stockNum" gorm:"column:stock_num;comment:商品库存数量;type:int"`
	Tag                string          `json:"tag" form:"tag" gorm:"column:tag;comment:商品标签;type:varchar(20);"`
	GoodsSellStatus    int             `json:"goodsSellStatus" form:"goodsSellStatus" gorm:"column:goods_sell_status;comment:商品上架状态 1-下架 0-上架;type:tinyint"`
	DaoFlag            int             `json:"daoFlag" form:"daoFlag" gorm:"column:dao_flag;comment:是否是节点商品 0否 1是 2nft;type:tinyint"`
	ProjectBeginTime   common.JSONTime `json:"projectBeginTime" form:"projectBeginTime" gorm:"column:project_begin_time;comment:项目开始时间;type:datetime"`
	ProjectEndTime     common.JSONTime `json:"projectEndTime" form:"projectEndTime" gorm:"column:project_end_time;comment:项目结束时间;type:datetime"`
	ProjectNum         int             `json:"projectNum" form:"projectNum" gorm:"column:project_num;comment:目标数;type:int"`
	CreateUser         int             `json:"createUser" form:"createUser" gorm:"column:create_user;comment:添加者主键id;type:int"`
	CreateTime         common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:商品添加时间;type:datetime"`
	UpdateUser         int             `json:"updateUser" form:"updateUser" gorm:"column:update_user;comment:修改者主键id;type:int"`
	UpdateTime         common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:商品修改时间;type:datetime"`
}

// TableName MallGoodsInfo 表名
func (MallGoodsInfo) TableName() string {
	return "tb_newbee_mall_goods_info"
}
