package mall

import "main.go/service"

type MallGroup struct {
	MallIndexApi
	MallGoodsInfoApi
	MallGoodsCategoryApi
	MallUserApi
	MallUserAddressApi
	MallShopCartApi
	MallOrderApi
	BscApi
}

var mallCarouselService = service.ServiceGroupApp.MallServiceGroup.MallCarouselService
var mallGoodsInfoService = service.ServiceGroupApp.MallServiceGroup.MallGoodsInfoService
var mallIndexConfigService = service.ServiceGroupApp.MallServiceGroup.MallIndexInfoService
var mallGoodsCategoryService = service.ServiceGroupApp.MallServiceGroup.MallGoodsCategoryService
var mallUserService = service.ServiceGroupApp.MallServiceGroup.MallUserService
var mallUserTokenService = service.ServiceGroupApp.MallServiceGroup.MallUserTokenService
var mallUserAddressService = service.ServiceGroupApp.MallServiceGroup.MallUserAddressService
var mallShopCartService = service.ServiceGroupApp.MallServiceGroup.MallShopCartService
var mallOrderService = service.ServiceGroupApp.MallServiceGroup.MallOrderService
var bscService = service.ServiceGroupApp.MallServiceGroup.BscService
