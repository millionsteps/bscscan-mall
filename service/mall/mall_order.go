package mall

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/bscscan"
	"main.go/model/common"
	"main.go/model/common/enum"
	"main.go/model/mall"
	"main.go/model/mall/response"
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

// SaveBscOrder 保存usdt订单
func (m *MallOrderService) SaveBscOrder(token string, myShoppingCartItems mallRes.BscOrderItemResponse) (err error, orderNo string) {
	var contract bscscan.BscContract
	err = global.GVA_DB.Where("id=?", myShoppingCartItems.ContractId).First(&contract).Error
	if err != nil {
		return errors.New("不存在的合约"), orderNo
	}
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), orderNo
	}
	var goodsIds []int
	goodsIds = append(goodsIds, myShoppingCartItems.GoodsId)
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
	//查出的商品中不存在购物车中的这条关联商品数据，直接返回错误提醒
	if _, ok := newBeeMallGoodsMap[myShoppingCartItems.GoodsId]; !ok {
		return errors.New("购物车数据异常！"), orderNo
	}
	if myShoppingCartItems.GoodsCount > newBeeMallGoodsMap[myShoppingCartItems.GoodsId].StockNum {
		return errors.New("库存不足！"), orderNo
	}

	//删除购物项
	if len(goodsIds) > 0 {
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
		newBeeMallOrder.FromAddress = myShoppingCartItems.FromAddress
		newBeeMallOrder.ToAddress = contract.ToAddress
		newBeeMallOrder.OrderNo = orderNo
		newBeeMallOrder.UserId = userToken.UserId
		//总价
		var goodsInfo manage.MallGoodsInfo
		global.GVA_DB.Where("goods_id =?", myShoppingCartItems.GoodsId).First(&goodsInfo)
		thisPrice := goodsInfo.SellingPrice.Mul(decimal.NewFromInt(int64(myShoppingCartItems.GoodsCount)))
		priceTotal = priceTotal.Add(thisPrice)
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
		//生成所有的订单项快照，并保存至数据库
		//返还三倍的金额
		uksa := priceTotal.Mul(decimal.NewFromInt(3))
		var newBeeMallOrderItem manage.MallOrderItem
		copier.Copy(&newBeeMallOrderItem, &myShoppingCartItems)
		newBeeMallOrderItem.GoodsName = goodsInfo.GoodsName
		newBeeMallOrderItem.GoodsCoverImg = goodsInfo.GoodsCoverImg
		newBeeMallOrderItem.OrderId = newBeeMallOrder.OrderId
		newBeeMallOrderItem.UserId = newBeeMallOrder.UserId
		newBeeMallOrderItem.CreateTime = common.JSONTime{Time: time.Now()}
		newBeeMallOrderItem.DaoFlag = goodsInfo.DaoFlag
		newBeeMallOrderItem.SellingPrice = goodsInfo.SellingPrice
		newBeeMallOrderItem.TotalPrice = uksa
		newBeeMallOrderItem.UsdtFreeze = uksa
		newBeeMallOrderItem.UsdtAble = uksa
		newBeeMallOrderItem.ReleaseFlag = 0
		if err = global.GVA_DB.Save(&newBeeMallOrderItem).Error; err != nil {
			return err, orderNo
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
func (m *MallOrderService) PaySuccessBsc(orderNo string, txHash string) (err error) {
	var mallOrder manage.MallOrder
	global.GVA_LOG.Info("orderNo:" + orderNo + ",txHash:" + txHash)
	err = global.GVA_DB.Where("order_no = ? and is_deleted=0 ", orderNo).First(&mallOrder).Error
	if mallOrder != (manage.MallOrder{}) {
		mallOrder.TxHash = txHash
		mallOrder.UpdateTime = common.JSONTime{time.Now()}
		if err = global.GVA_DB.Save(&mallOrder).Error; err != nil {
			return errors.New("保存订单hash失败")
		}
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
		checkErr, isSuccess := CheckOrder(txHash, mallOrder.FromAddress, mallOrder.ToAddress)
		if checkErr != nil {
			return errors.New("支付校验失败")
		}
		if !isSuccess {
			i := 0
			for i < 3 {
				time.Sleep(500 * time.Millisecond)
				checkErr, isSuccess = CheckOrder(txHash, mallOrder.FromAddress, mallOrder.ToAddress)
				if isSuccess {
					break
				}
				i++
			}
		}

		if !isSuccess {
			return errors.New("支付校验失败")
		}
		orderItem.ReleaseFlag = 1
		orderItem.ReleaseRate = decimal.NewFromFloat32(0.01)
		if err = global.GVA_DB.Where("order_item_id=?", orderItem.OrderItemId).UpdateColumns(&orderItem).Error; err != nil {
			return errors.New("修改订单明细失败")
		}
		mallOrder.OrderStatus = enum.ORDER_PAID.Code()
		mallOrder.PayStatus = 1
		mallOrder.PayTime = common.JSONTime{time.Now()}
		mallOrder.UpdateTime = common.JSONTime{time.Now()}
		if err = global.GVA_DB.Save(&mallOrder).Error; err != nil {
			return errors.New("修改订单失败")
		}
		//计算本人消费的业绩
		var totalUsdt decimal.Decimal
		totalUsdtErr := global.GVA_DB.Model(&manage.MallOrderItem{}).Where("user_id = ? and release_flag != 0", mallOrder.UserId).Select("sum(total_price)").Scan(&totalUsdt).Error
		if err != nil {
			global.GVA_LOG.Error("查询计算本人消费的业绩失败", zap.Error(totalUsdtErr))
		}
		//查询用户账户
		var userAccount bscscan.BscMallUserAccount
		userAccountErr := global.GVA_DB.Where("user_id = ?", mallOrder.UserId).First(&userAccount).Error
		if userAccountErr != nil {
			global.GVA_LOG.Error("用户账户不存在", zap.Error(userAccountErr))
		}
		userAccount.TotalUsdt = totalUsdt
		updateAccountErr := global.GVA_DB.Where("id = ?", userAccount.Id).UpdateColumns(&userAccount).Error
		if updateAccountErr != nil {
			global.GVA_LOG.Error("更新用户账户失败", zap.Error(userAccountErr))
		}
		//判断节点或者用户是否升级
		daoFlag := goods.DaoFlag
		var sourceType int
		var sourceContent string
		sourceContent = "购买产品" + orderItem.GoodsName
		//判断是否是节点商品
		if daoFlag == 1 {
			sourceType = 0
			m.daoGoodsInfo(mallOrder.UserId)
		} else {
			//计算业绩 判断上级是否升级等级
			m.countTotalUsdt(mallOrder.UserId)
			sourceType = 2
		}
		//添加明细
		var detail bscscan.BscMallAccountDetail
		detail.UserId = mallOrder.UserId
		detail.Usdt = orderItem.UsdtFreeze
		detail.SourceType = sourceType
		detail.SourceContent = sourceContent
		detail.UpdateTime = common.JSONTime{time.Now()}
		detail.CreateTime = common.JSONTime{time.Now()}
		detail.Type = 0
		if err = global.GVA_DB.Save(&detail).Error; err != nil {
			return errors.New("保存账户明细失败")
		}
	}
	return
}

func (m *MallOrderService) countWeakSideUsdt(userId int) (err error) {
	var userAccount bscscan.BscMallUserAccount
	err = global.GVA_DB.Where("user_id = ?", userId).First(&userAccount).Error
	if err != nil {
		return errors.New("用户账户不存在")
	}
	//计算 A侧
	var userA mall.MallUser
	err = global.GVA_DB.Where("parent_id = ? and node_type = 'A'", userId).First(&userA).Error
	if err != nil {
		return errors.New("用户不存在")
	}
	usdtA := m.getSubModeUsdt(userA.UserId)
	//A的业绩
	var userAccountA bscscan.BscMallUserAccount
	accountAErr := global.GVA_DB.Where("user_id = ?", userA.UserId).First(&userAccountA).Error
	if accountAErr != nil {
		global.GVA_LOG.Error("查询A账户失败", zap.Error(accountAErr))
	}
	userAccount.TotalUsdtDownA = usdtA
	usdtA = usdtA.Add(userAccountA.TotalUsdt)
	//计算 B侧
	var userB mall.MallUser
	err = global.GVA_DB.Where("parent_id = ? and node_type = 'B'", userId).First(&userB).Error
	if err != nil {
		return errors.New("用户不存在")
	}
	usdtB := m.getSubModeUsdt(userB.UserId)
	//B的业绩
	var userAccountB bscscan.BscMallUserAccount
	accountBErr := global.GVA_DB.Where("user_id = ?", userA.UserId).First(&userAccountB).Error
	if accountBErr != nil {
		global.GVA_LOG.Error("查询B账户失败", zap.Error(accountBErr))
	}
	usdtB = usdtB.Add(userAccountB.TotalUsdt)
	userAccount.TotalUsdtDownB = usdtB
	//如果相等随便取一侧判断等级
	thisUsdt := usdtA
	if usdtA.Cmp(usdtB) == 1 {
		thisUsdt = usdtB
	}
	//伞下所有人的业绩
	totalUsdtDown := usdtA.Add(usdtB)
	userAccount.TotalUsdtDown = totalUsdtDown
	level := getLevelByUsdt(thisUsdt)
	vipLevel := userAccount.VipLevel
	if level > vipLevel {
		userAccount.VipLevel = level
		err = global.GVA_DB.Where("id = ?", userAccount.Id).UpdateColumns(&userAccount).Error
		if err != nil {
			return errors.New("更新用户账户失败")
		}
	}
	return
}

//判断等级计算业绩 自己的业绩不计入
func (m *MallOrderService) countTotalUsdt(userId int) (err error) {
	var ids []int
	var parentIds []int
	err, parentIds = getParentId(userId, ids)
	if err != nil {
		return errors.New("查询所有父级id失败")
	}
	for _, parentId := range parentIds {
		//计算 两侧的业绩 并拿到弱侧业绩
		m.countWeakSideUsdt(parentId)
	}
	return
}

func getParentId(userId int, ids []int) (err error, parentIds []int) {
	var userAccount bscscan.BscMallUserAccount
	err = global.GVA_DB.Where("user_id = ?", userId).First(&userAccount).Error
	if err != nil {
		return errors.New("用户账户不存在"), parentIds
	}

	//遍历所有上级
	parentId := userAccount.ParentId
	if parentId != 0 {
		ids = append(ids, parentId)
		err, parentIds = getParentId(parentId, ids)
		if err != nil {
			global.GVA_LOG.Error("查询父级用户失败", zap.Error(err))
		}
	} else {
		return
	}
	return
}

//计算等级
func getLevelByUsdt(usdt decimal.Decimal) (level int) {
	if usdt.Cmp(decimal.NewFromInt(20000)) >= 0 && usdt.Cmp(decimal.NewFromInt(80000)) == -1 {
		return 1
	} else if usdt.Cmp(decimal.NewFromInt(80000)) >= 0 && usdt.Cmp(decimal.NewFromInt(300000)) == -1 {
		return 2
	} else if usdt.Cmp(decimal.NewFromInt(300000)) >= 0 && usdt.Cmp(decimal.NewFromInt(1000000)) == -1 {
		return 3
	} else if usdt.Cmp(decimal.NewFromInt(1000000)) >= 0 && usdt.Cmp(decimal.NewFromInt(3000000)) == -1 {
		return 4
	} else if usdt.Cmp(decimal.NewFromInt(3000000)) >= 0 {
		return 5
	}
	return
}
func (m *MallOrderService) getSubModeUsdt(userId int) (usdt decimal.Decimal) {
	var accountList []bscscan.BscMallUserAccount
	err := global.GVA_DB.Where("parent_id = ?", userId).Find(&accountList).Error
	if err != nil {
		global.GVA_LOG.Error("查询账户失败", zap.Error(err))
		return
	}
	for _, account := range accountList {
		thisUserId := account.UserId
		totalUsdt := account.TotalUsdt
		usdt = usdt.Add(totalUsdt)
		usdt = usdt.Add(m.getSubModeUsdt(thisUserId))
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
	//usdt 大于5000 就是 可以分红
	user.BonusFlag = 1
	//修改用户类型
	err = global.GVA_DB.Where("user_id = ?", userId).UpdateColumns(&user).Error
	if err != nil {
		return errors.New("保存失败")
	}
	//节点商品计算总usdt
	// var sumUsdt decimal.Decimal
	// err = global.GVA_DB.Table("tb_newbee_mall_order_item").Select("sum(total_price)").Where("user_id=?", userId).Scan(&sumUsdt).Error
	// if sumUsdt.GreaterThanOrEqual(decimal.NewFromInt(5000)) {

	// }
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

// OrderItemList 订单明细列表
func (m *MallOrderService) OrderItemList(pageNumber int, token string) (err error, list []response.NewBeeMallOrderItemVO, total int64) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), list, total
	}
	db := global.GVA_DB.Model(&manage.MallOrderItem{})
	err = db.Where("user_id=? and release_flag != 0", userToken.UserId).Count(&total).Error
	if err != nil {
		return errors.New("查询总数失败"), list, total
	}
	pageSize := 10
	limit := pageSize
	offset := pageSize * (pageNumber - 1)
	var itemList []response.NewBeeMallOrderItemVO
	err = db.Limit(limit).Offset(offset).Order("create_time asc").Find(&itemList).Error
	if err != nil {
		return errors.New("查询失败"), list, total
	}
	return err, itemList, total
}

// BSCProjectOrderListBySearch 搜索项目订单
func (m *MallOrderService) BSCProjectOrderListBySearch(pageNumber int, goodsId int) (err error, list []mallRes.MallOrderResponse, total int64) {
	// 根据搜索条件查询
	var newBeeMallOrders []manage.MallOrder
	db := global.GVA_DB.Table("tb_newbee_mall_order as o")
	db.Joins(" left join tb_newbee_mall_order_item as i on o.order_id=i.order_id")
	err = db.Where("i.goods_id =? and o.order_status = 1", goodsId).Count(&total).Error
	db.Select("o.*")
	//这里前段没有做滚动加载，直接显示全部订单
	limit := 5
	offset := 5 * (pageNumber - 1)
	err = db.Limit(limit).Offset(offset).Order("o.order_id desc").Find(&newBeeMallOrders).Error

	var orderListVOS []mallRes.MallOrderResponse

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

	return err, list, total
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
	limit := 5
	offset := 5 * (pageNumber - 1)
	err = db.Limit(limit).Offset(offset).Order(" order_id desc").Find(&newBeeMallOrders).Error

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
		}
		//已存在不用再次计算
		if bscOrderItemRelease != (bscscan.BscOrderItemRelease{}) {
			continue
		}
		global.GVA_LOG.Info("开始执行释放逻辑")
		//判断节点商品
		daoFlag := orderItem.DaoFlag
		rate := decimal.NewFromInt(0)
		if daoFlag == 0 {
			//普通商品分销逻辑
			rate = releaseGeneralGoodsUsdt(orderItem, timeStr)
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
		fmt.Println("releaseRate:", releaseRate)
		fmt.Println("rate:", rate)
		usdtAble := orderItem.UsdtAble
		bscOrderItemRelease.UsdtBegin = usdtAble
		usdt := orderItem.Usdt
		//本次释放
		var thisUsdt decimal.Decimal
		if rate.Cmp(decimal.Zero) == 0 {
			thisUsdt = usdtFreeze.Mul(releaseRate)
		} else {
			thisUsdt = usdtFreeze.Mul(releaseRate).Mul(rate)
		}
		fmt.Println("usdtFreeze:", usdtFreeze)
		fmt.Println("thisUsdt:", thisUsdt)
		bscOrderItemRelease.ThisUsdt = thisUsdt

		// 添加到账户余额
		userAccount := getUserAccount(orderItem.UserId)
		userUsdt := userAccount.Usdt
		resultUsdt := userUsdt.Add(thisUsdt)
		userAccount.Usdt = resultUsdt
		//保存释放记录
		userAccountSaveErr := global.GVA_DB.Save(&userAccount).Error
		if userAccountSaveErr != nil {
			global.GVA_LOG.Error("更新用户账户失败", zap.Error(userAccountSaveErr))
			continue
		}

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
		releaseSaveErr := global.GVA_DB.Save(&bscOrderItemRelease).Error
		if releaseSaveErr != nil {
			global.GVA_LOG.Error("保存记录失败", zap.Error(err))
			return
		}
	}

	//todo 统一转账
	return
}

func getUserAccount(userId int) (userAccount bscscan.BscMallUserAccount) {
	userAccountErr := global.GVA_DB.Where("user_id = ?", userId).First(&userAccount).Error
	if userAccountErr != nil {
		global.GVA_LOG.Error("查询用户账户失败", zap.Error(userAccountErr))
	}
	if userAccount == (bscscan.BscMallUserAccount{}) {
		var user mall.MallUser
		userErr := global.GVA_DB.Where("user_id = ?", userId).First(&user).Error
		if userErr != nil {
			global.GVA_LOG.Error("查询用户记录失败", zap.Error(userErr))
			return
		}
		userAccount.UserId = userId
		userAccount.ParentId = user.ParentId
		userAccount.CreateTime = common.JSONTime{Time: time.Now()}
		userAccount.UpdateTime = common.JSONTime{Time: time.Now()}
		userAccount.Usdt = decimal.Zero
		userAccount.TotalUsdt = decimal.Zero
		userAccount.UsdtFreeze = decimal.Zero
		//保存用户账户记录
		userAccountSaveErr := global.GVA_DB.Save(&userAccount).Error
		if userAccountSaveErr != nil {
			global.GVA_LOG.Error("保存记录失败", zap.Error(userAccountSaveErr))
			return
		}
	}
	return
}

//普通商品分红
func releaseGeneralGoodsUsdt(orderItem manage.MallOrderItem, timeStr string) (rate decimal.Decimal) {
	userId := orderItem.UserId
	var parents map[int]int
	err, parentsResult := selectAllLevel(userId, 0, parents)
	if err != nil {
		return
	}
	usedRate := 0
	for i := 0; i < len(parentsResult); i++ {
		//当前等级
		level := i + 1
		//当前等级的用户id
		userId = parentsResult[level]
		if userId == 0 {
			continue
		}
		//当前等级可拿到比例
		thisRate := level*10 - usedRate
		//当前等级已占用的比例
		usedRate = level * 10
		//每人分销抽取比例
		usdtFreeze := orderItem.UsdtFreeze
		releaseRate := orderItem.ReleaseRate
		thisRateDec := decimal.NewFromFloat32(0.01).Mul(decimal.NewFromInt(int64(thisRate)))
		//最终拿到的分销比例
		resultRate := releaseRate.Mul(decimal.NewFromFloat(0.01)).Mul(thisRateDec)
		//计算本次释放给该等级用户的u
		thisUsdt := usdtFreeze.Mul(resultRate)

		var bscOrderItemRelease bscscan.BscOrderItemRelease
		bscOrderItemRelease.UserId = userId
		bscOrderItemRelease.OrderId = orderItem.OrderId
		bscOrderItemRelease.OrderItemId = orderItem.OrderItemId
		bscOrderItemRelease.ReleaseState = 0
		bscOrderItemRelease.RelesaeDate = timeStr
		bscOrderItemRelease.ReleaseRate = resultRate
		bscOrderItemRelease.UsdtFreeze = orderItem.UsdtFreeze
		bscOrderItemRelease.ThisUsdt = thisUsdt

		//保存释放记录
		releaseSaveErr := global.GVA_DB.Save(bscOrderItemRelease).Error
		if releaseSaveErr != nil {
			global.GVA_LOG.Error("保存记录失败", zap.Error(err))
			return
		}
	}
	//返回已用比例
	return decimal.NewFromInt(int64(100 - usedRate)).Mul(decimal.NewFromFloat(0.01))
}

//遍历所有上级拿到可参与分销的会员
func selectAllLevel(userId int, level int, parents map[int]int) (err error, parentsResult map[int]int) {
	var userAccount bscscan.BscMallUserAccount
	err = global.GVA_DB.Where("user_id = ?", userId).First(&userAccount).Error
	if err != nil {
		return errors.New("用户账户不存在"), parentsResult
	}

	//遍历所有上级
	vipLevel := userAccount.VipLevel
	if vipLevel != 0 && level < vipLevel {
		parents[vipLevel] = userAccount.UserId
		if vipLevel+1 > 5 {
			parentsResult = parents
			return
		}
		if userAccount.ParentId == 0 {
			parentsResult = parents
			return
		}
		err, parentsResult = selectAllLevel(userAccount.ParentId, vipLevel+1, parents)
	} else {
		if userAccount.ParentId != 0 {
			err, parentsResult = selectAllLevel(userAccount.ParentId, level, parents)
		} else {
			parentsResult = parents
			return
		}
	}
	return
}
