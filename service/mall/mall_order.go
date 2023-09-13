package mall

import (
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/bscscan"
	"main.go/model/common"
	"main.go/model/common/enum"
	"main.go/model/mall"
	mallRes "main.go/model/mall/response"
	"main.go/model/manage"
	manageReq "main.go/model/manage/request"
	"main.go/utils"
)

type MallOrderService struct {
}

// SaveOrder 保存订单
func (m *MallOrderService) SaveOrder(token string, userAddress mall.MallUserAddress, myShoppingCartItems []mallRes.CartItemResponse) (err error, orderNo string) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), orderNo
	}
	var itemIdList []int
	var goodsIds []int
	for _, cartItem := range myShoppingCartItems {
		itemIdList = append(itemIdList, cartItem.CartItemId)
		goodsIds = append(goodsIds, cartItem.GoodsId)
	}
	var newBeeMallGoods []manage.MallGoodsInfo
	global.GVA_DB.Where("goods_id in ? ", goodsIds).Find(&newBeeMallGoods)
	//检查是否包含已下架商品
	for _, mallGoods := range newBeeMallGoods {
		if mallGoods.GoodsSellStatus != enum.GOODS_UNDER.Code() {
			return errors.New("已下架，无法生成订单"), orderNo
		}
	}
	newBeeMallGoodsMap := make(map[int]manage.MallGoodsInfo)
	for _, mallGoods := range newBeeMallGoods {
		newBeeMallGoodsMap[mallGoods.GoodsId] = mallGoods
	}
	//判断商品库存
	for _, shoppingCartItemVO := range myShoppingCartItems {
		//查出的商品中不存在购物车中的这条关联商品数据，直接返回错误提醒
		if _, ok := newBeeMallGoodsMap[shoppingCartItemVO.GoodsId]; !ok {
			return errors.New("购物车数据异常！"), orderNo
		}
		if shoppingCartItemVO.GoodsCount > newBeeMallGoodsMap[shoppingCartItemVO.GoodsId].StockNum {
			return errors.New("库存不足！"), orderNo
		}
	}
	//删除购物项
	if len(itemIdList) > 0 && len(goodsIds) > 0 {
		if err = global.GVA_DB.Where("cart_item_id in ?", itemIdList).Updates(mall.MallShoppingCartItem{IsDeleted: 1}).Error; err == nil {
			var stockNumDTOS []manageReq.StockNumDTO
			copier.Copy(&stockNumDTOS, &myShoppingCartItems)
			for _, stockNumDTO := range stockNumDTOS {
				var goodsInfo manage.MallGoodsInfo
				global.GVA_DB.Where("goods_id =?", stockNumDTO.GoodsId).First(&goodsInfo)
				if err = global.GVA_DB.Where("goods_id =? and stock_num>= ? and goods_sell_status = 0", stockNumDTO.GoodsId, stockNumDTO.GoodsCount).Updates(manage.MallGoodsInfo{StockNum: goodsInfo.StockNum - stockNumDTO.GoodsCount}).Error; err != nil {
					return errors.New("库存不足！"), orderNo
				}
			}
			//生成订单号
			orderNo = utils.GenOrderNo()
			priceTotal := decimal.Zero
			//保存订单
			var newBeeMallOrder manage.MallOrder
			newBeeMallOrder.OrderNo = orderNo
			newBeeMallOrder.UserId = userToken.UserId
			//总价
			for _, newBeeMallShoppingCartItemVO := range myShoppingCartItems {
				thisPrice := newBeeMallShoppingCartItemVO.SellingPrice.Mul(decimal.NewFromInt(int64(newBeeMallShoppingCartItemVO.GoodsCount)))
				priceTotal = priceTotal.Add(thisPrice)
			}
			if priceTotal.LessThanOrEqual(decimal.NewFromInt(0)) {
				return errors.New("订单价格异常！"), orderNo
			}
			newBeeMallOrder.CreateTime = common.JSONTime{Time: time.Now()}
			newBeeMallOrder.UpdateTime = common.JSONTime{Time: time.Now()}
			newBeeMallOrder.TotalPrice = priceTotal
			newBeeMallOrder.ExtraInfo = ""
			//生成订单项并保存订单项纪录
			if err = global.GVA_DB.Save(&newBeeMallOrder).Error; err != nil {
				return errors.New("订单入库失败！"), orderNo
			}
			//生成订单收货地址快照，并保存至数据库
			var newBeeMallOrderAddress mall.MallOrderAddress
			copier.Copy(&newBeeMallOrderAddress, &userAddress)
			newBeeMallOrderAddress.OrderId = newBeeMallOrder.OrderId
			//生成所有的订单项快照，并保存至数据库
			var newBeeMallOrderItems []manage.MallOrderItem
			for _, newBeeMallShoppingCartItemVO := range myShoppingCartItems {
				var newBeeMallOrderItem manage.MallOrderItem
				copier.Copy(&newBeeMallOrderItem, &newBeeMallShoppingCartItemVO)
				newBeeMallOrderItem.OrderId = newBeeMallOrder.OrderId
				newBeeMallOrderItem.CreateTime = common.JSONTime{Time: time.Now()}
				newBeeMallOrderItems = append(newBeeMallOrderItems, newBeeMallOrderItem)
			}
			if err = global.GVA_DB.Save(&newBeeMallOrderItems).Error; err != nil {
				return err, orderNo
			}
		}
	}
	return
}

// PaySuccess 支付订单
func (m *MallOrderService) PaySuccess(orderNo string, payType int) (err error) {
	var mallOrder manage.MallOrder
	err = global.GVA_DB.Where("order_no = ? and is_deleted=0 ", orderNo).First(&mallOrder).Error
	if mallOrder != (manage.MallOrder{}) {
		if mallOrder.OrderStatus != 0 {
			return errors.New("订单状态异常！")
		}
		mallOrder.OrderStatus = enum.ORDER_PAID.Code()
		mallOrder.PayType = payType
		mallOrder.PayStatus = 1
		mallOrder.PayTime = common.JSONTime{time.Now()}
		mallOrder.UpdateTime = common.JSONTime{time.Now()}
		err = global.GVA_DB.Save(&mallOrder).Error
	}
	return
}

// PaySuccessBsc 支付成功订单
func (m *MallOrderService) PaySuccessBsc(orderNo string, payType int) (err error) {
	var mallOrder manage.MallOrder
	err = global.GVA_DB.Where("order_no = ? and is_deleted=0 ", orderNo).First(&mallOrder).Error
	if mallOrder != (manage.MallOrder{}) {
		orderId := mallOrder.OrderId
		var orderItem manage.MallOrderItem
		itemErr := global.GVA_DB.Where("order_id = ? ", orderId).First(&orderItem).Error
		if itemErr != nil {
			return errors.New("订单明细不存在")
		}
		goodsId := orderItem.GoodsId
		var goods manage.MallGoodsInfo
		goodsErr := global.GVA_DB.Where("goods_id = ? ", goodsId).First(&goods).Error
		if goodsErr != nil {
			return errors.New("商品不存在")
		}
		daoFlag := goods.DaoFlag
		//判断是否是节点商品
		if daoFlag == 1 {
			m.daoGoodsInfo(mallOrder.UserId)
		}
	}
	return
}
func (m *MallOrderService) daoGoodsInfo(userId int) (err error) {
	//判断用户是否可以分红
	var user mall.MallUser
	err = global.GVA_DB.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return errors.New("用户不存在")
	}
	bonusFlag := user.BonusFlag
	if bonusFlag == 1 {
		return errors.New("用户已经是可以分红角色")
	}
	//节点商品计算总usdt
	var sumUsdt decimal.Decimal
	err = global.GVA_DB.Table("tb_newbee_mall_order_item").Select("sum(total_price)").Where("user_id=?", userId).Scan(&sumUsdt).Error
	if sumUsdt.GreaterThanOrEqual(decimal.NewFromInt(5000)) {
		//usdt 大于5000 就是 可以分红
		user.BonusFlag = 1
		//修改用户类型
		err = global.GVA_DB.Where("user_id = ?", userId).UpdateColumns(&user).Error
		if err != nil {
			return errors.New("保存失败")
		}
	}
	return
}

// FinishOrder 完结订单
func (m *MallOrderService) FinishOrder(token string, orderNo string) (err error) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户")
	}
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	if mallOrder.UserId != userToken.UserId {
		return errors.New("禁止该操作！")
	}
	mallOrder.OrderStatus = enum.ORDER_SUCCESS.Code()
	mallOrder.UpdateTime = common.JSONTime{time.Now()}
	err = global.GVA_DB.Save(&mallOrder).Error
	return
}

// CancelOrder 关闭订单
func (m *MallOrderService) CancelOrder(token string, orderNo string) (err error) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户")
	}
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	if mallOrder.UserId != userToken.UserId {
		return errors.New("禁止该操作！")
	}
	if utils.NumsInList(mallOrder.OrderStatus, []int{enum.ORDER_SUCCESS.Code(),
		enum.ORDER_CLOSED_BY_MALLUSER.Code(), enum.ORDER_CLOSED_BY_EXPIRED.Code(), enum.ORDER_CLOSED_BY_JUDGE.Code()}) {
		return errors.New("订单状态异常！")
	}
	mallOrder.OrderStatus = enum.ORDER_CLOSED_BY_MALLUSER.Code()
	mallOrder.UpdateTime = common.JSONTime{time.Now()}
	err = global.GVA_DB.Save(&mallOrder).Error
	return
}

// GetOrderDetailByOrderNo 获取订单详情
func (m *MallOrderService) GetOrderDetailByOrderNo(token string, orderNo string) (err error, orderDetail mallRes.MallOrderDetailVO) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), orderDetail
	}
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！"), orderDetail
	}
	if mallOrder.UserId != userToken.UserId {
		return errors.New("禁止该操作！"), orderDetail
	}
	var orderItems []manage.MallOrderItem
	err = global.GVA_DB.Where("order_id = ?", mallOrder.OrderId).Find(&orderItems).Error
	if len(orderItems) <= 0 {
		return errors.New("订单项不存在！"), orderDetail
	}

	var newBeeMallOrderItemVOS []mallRes.NewBeeMallOrderItemVO
	copier.Copy(&newBeeMallOrderItemVOS, &orderItems)
	copier.Copy(&orderDetail, &mallOrder)
	// 订单状态前端显示为中文
	_, OrderStatusStr := enum.GetNewBeeMallOrderStatusEnumByStatus(orderDetail.OrderStatus)
	_, payTapStr := enum.GetNewBeeMallOrderStatusEnumByStatus(orderDetail.PayType)
	orderDetail.OrderStatusString = OrderStatusStr
	orderDetail.PayTypeString = payTapStr
	orderDetail.NewBeeMallOrderItemVOS = newBeeMallOrderItemVOS

	return
}

// MallOrderListBySearch 搜索订单
func (m *MallOrderService) MallOrderListBySearch(token string, pageNumber int, status string) (err error, list []mallRes.MallOrderResponse, total int64) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), list, total
	}
	// 根据搜索条件查询
	var newBeeMallOrders []manage.MallOrder
	db := global.GVA_DB.Model(&newBeeMallOrders)

	if status != "" {
		db.Where("order_status = ?", status)
	}
	err = db.Where("user_id =? and is_deleted=0 ", userToken.UserId).Count(&total).Error
	//这里前段没有做滚动加载，直接显示全部订单
	//limit := 5
	offset := 5 * (pageNumber - 1)
	err = db.Offset(offset).Order(" order_id desc").Find(&newBeeMallOrders).Error

	var orderListVOS []mallRes.MallOrderResponse
	if total > 0 {
		//数据转换 将实体类转成vo
		copier.Copy(&orderListVOS, &newBeeMallOrders)
		//设置订单状态中文显示值
		for _, newBeeMallOrderListVO := range orderListVOS {
			_, statusStr := enum.GetNewBeeMallOrderStatusEnumByStatus(newBeeMallOrderListVO.OrderStatus)
			newBeeMallOrderListVO.OrderStatusString = statusStr
		}
		// 返回订单id
		var orderIds []int
		for _, order := range newBeeMallOrders {
			orderIds = append(orderIds, order.OrderId)
		}
		//获取OrderItem
		var orderItems []manage.MallOrderItem
		if len(orderIds) > 0 {
			global.GVA_DB.Where("order_id in ?", orderIds).Find(&orderItems)
			itemByOrderIdMap := make(map[int][]manage.MallOrderItem)
			for _, orderItem := range orderItems {
				itemByOrderIdMap[orderItem.OrderId] = []manage.MallOrderItem{}
			}
			for k, v := range itemByOrderIdMap {
				for _, orderItem := range orderItems {
					if k == orderItem.OrderId {
						v = append(v, orderItem)
					}
					itemByOrderIdMap[k] = v
				}
			}
			//封装每个订单列表对象的订单项数据
			for _, newBeeMallOrderListVO := range orderListVOS {
				if _, ok := itemByOrderIdMap[newBeeMallOrderListVO.OrderId]; ok {
					orderItemListTemp := itemByOrderIdMap[newBeeMallOrderListVO.OrderId]
					var newBeeMallOrderItemVOS []mallRes.NewBeeMallOrderItemVO
					copier.Copy(&newBeeMallOrderItemVOS, &orderItemListTemp)
					newBeeMallOrderListVO.NewBeeMallOrderItemVOS = newBeeMallOrderItemVOS
					_, OrderStatusStr := enum.GetNewBeeMallOrderStatusEnumByStatus(newBeeMallOrderListVO.OrderStatus)
					newBeeMallOrderListVO.OrderStatusString = OrderStatusStr
					list = append(list, newBeeMallOrderListVO)
				}
			}
		}
	}
	return err, list, total
}

// ReleaseUsdt 每天释放1%
func ReleaseUsdt() {
	global.GVA_LOG.Info("ReleaseUsdt--每天释放1%")
	var orderItems []manage.MallOrderItem
	itemErr := global.GVA_DB.Where("release_flag = 1").Find(&orderItems).Error
	if itemErr != nil {
		global.GVA_LOG.Error("查询可释放订单失败", zap.Error(itemErr))
		return
	}
	timeStr := time.Now().Format("2006-01-02")

	for _, orderItem := range orderItems {
		//查询当前订单明细是否已释放
		var bscOrderItemRelease bscscan.BscOrderItemRelease
		releaseErr := global.GVA_DB.Where("user_id=? and order_item_id =? and relesae_date=?", orderItem.UserId, orderItem.OrderItemId, timeStr).First(&bscOrderItemRelease).Error
		if releaseErr != nil {
			global.GVA_LOG.Error("查询释放订单记录失败", zap.Error(releaseErr))
			continue
		}
		//已存在不用再次计算
		if bscOrderItemRelease != (bscscan.BscOrderItemRelease{}) {
			continue
		}
		bscOrderItemRelease.UserId = orderItem.UserId
		bscOrderItemRelease.OrderId = orderItem.OrderId
		bscOrderItemRelease.OrderItemId = orderItem.OrderItemId
		bscOrderItemRelease.ReleaseState = 0
		bscOrderItemRelease.RelesaeDate = timeStr
		bscOrderItemRelease.ReleaseRate = orderItem.ReleaseRate
		bscOrderItemRelease.UsdtFreeze = orderItem.UsdtFreeze

		usdtFreeze := orderItem.UsdtFreeze
		releaseRate := orderItem.ReleaseRate
		usdtAble := orderItem.UsdtAble
		bscOrderItemRelease.UsdtBegin = usdtAble
		usdt := orderItem.Usdt
		//本次释放
		thisUsdt := usdtFreeze.Mul(releaseRate).Mul(decimal.NewFromFloat(0.01))
		bscOrderItemRelease.ThisUsdt = thisUsdt
		//todo 添加到账户余额
		//本次剩余
		thisUsdtAble := usdtAble.Sub(thisUsdt)
		bscOrderItemRelease.UsdtEnd = thisUsdtAble
		orderItem.UsdtAble = thisUsdtAble
		//累计已释放
		totalUsdt := usdt.Add(thisUsdt)
		orderItem.Usdt = totalUsdt
		err := global.GVA_DB.Where("order_item_id = ?", orderItem.OrderItemId).UpdateColumns(&orderItem).Error
		if err != nil {
			global.GVA_LOG.Error("更新订单失败", zap.Error(err))
			return
		}
		//保存释放记录
		releaseSaveErr := global.GVA_DB.Save(bscOrderItemRelease).Error
		if releaseSaveErr != nil {
			global.GVA_LOG.Error("保存记录失败", zap.Error(err))
			return
		}
	}

	//todo 统一转账
	return
}
