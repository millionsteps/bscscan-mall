package mall

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main.go/global"
	"main.go/middleware/api"
	"main.go/model/bscscan"
	"main.go/model/common"
	"main.go/model/mall"
	mallReq "main.go/model/mall/request"
	"main.go/model/mall/response"
	mallRes "main.go/model/mall/response"
	"main.go/model/manage"
	"main.go/utils"
)

type MallUserService struct {
}

// RegisterUser 注册用户
func (m *MallUserService) RegisterUser(req mallReq.RegisterUserParam) (err error) {
	if !errors.Is(global.GVA_DB.Where("login_name =?", req.LoginName).First(&mall.MallUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同用户名")
	}

	return global.GVA_DB.Create(&mall.MallUser{
		LoginName:     req.LoginName,
		PasswordMd5:   utils.MD5V([]byte(req.Password)),
		IntroduceSign: "随新所欲，蜂富多彩",
		CreateTime:    common.JSONTime{Time: time.Now()},
	}).Error

}

func (m *MallUserService) UpdateUserInfo(token string, req mallReq.UpdateUserInfoParam) (err error) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户")
	}
	var userInfo mall.MallUser
	err = global.GVA_DB.Where("user_id =?", userToken.UserId).First(&userInfo).Error
	// 若密码为空字符，则表明用户不打算修改密码，使用原密码保存
	if !(req.PasswordMd5 == "") {
		userInfo.PasswordMd5 = utils.MD5V([]byte(req.PasswordMd5))
	}
	if !(req.EmailAddress == "") {
		userInfo.EmailAddress = req.EmailAddress
	}
	if !(req.NickName == "") {
		userInfo.NickName = req.NickName
	}
	if !(req.IntroduceSign == "") {
		userInfo.IntroduceSign = req.IntroduceSign
	}
	err = global.GVA_DB.Where("user_id =?", userToken.UserId).UpdateColumns(&userInfo).Error
	return
}

func (m *MallUserService) GetUserBonusInfo(token string) (err error, userDetail mallRes.MallUserBonusDetailResponse) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), userDetail
	}
	var account bscscan.BscMallUserAccount
	err = global.GVA_DB.Where("user_id =?", userToken.UserId).First(&account).Error
	if err != nil {
		return errors.New("用户账户获取失败"), userDetail
	}
	var userInfo mall.MallUser
	err = global.GVA_DB.Where("user_id =?", userToken.UserId).First(&userInfo).Error
	if err != nil {
		return errors.New("用户信息获取失败"), userDetail
	}
	if err != nil {
		userAccount := bscscan.BscMallUserAccount{
			UserId:        userInfo.UserId,
			ParentId:      userInfo.ParentId,
			VipLevel:      0,
			Usdt:          decimal.NewFromInt(0),
			UsdtFreeze:    decimal.NewFromInt(0),
			TotalUsdt:     decimal.NewFromInt(0),
			TotalUsdtDown: decimal.NewFromInt(0),
			CreateTime:    common.JSONTime{Time: time.Now()},
			UpdateTime:    common.JSONTime{Time: time.Now()},
		}
		err = global.GVA_DB.Create(&userAccount).Error
		if err != nil {
			log.Panic(api.NewException(api.UserNotExist))
			return
		}
		account = userAccount
	}
	userDetail.Usdt = account.Usdt
	//查询用户已提现金额
	var withdrawUsdt decimal.Decimal
	err = global.GVA_DB.Model(&bscscan.BscWithdrawRecord{}).Where("user_id = ?", userToken.UserId).Select("sum(usdt)").Scan(&withdrawUsdt).Error
	if err != nil {
		global.GVA_LOG.Error("查询用户已提现金额失败！", zap.Error(err))
	}
	userDetail.WithdrawUsdt = withdrawUsdt
	//计算节点数量
	var size int
	err = global.GVA_DB.Model(&bscscan.BscMallUserAccount{}).Where("dao_flag = 1").Select("count(1)").Scan(&size).Error
	if err != nil {
		return errors.New("计算节点数量失败！"), userDetail
	}
	userDetail.DaoNum = size
	return
}

func (m *MallUserService) GetUserDetail(token string) (err error, userDetail mallRes.MallUserDetailResponse) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), userDetail
	}
	var userInfo mall.MallUser
	err = global.GVA_DB.Where("user_id =?", userToken.UserId).First(&userInfo).Error
	if err != nil {
		return errors.New("用户信息获取失败"), userDetail
	}
	err = copier.Copy(&userDetail, &userInfo)
	var account bscscan.BscMallUserAccount
	err = global.GVA_DB.Where("user_id =?", userInfo.UserId).First(&account).Error
	if err != nil {
		userAccount := bscscan.BscMallUserAccount{
			UserId:        userInfo.UserId,
			ParentId:      userInfo.ParentId,
			VipLevel:      0,
			Usdt:          decimal.NewFromInt(0),
			UsdtFreeze:    decimal.NewFromInt(0),
			TotalUsdt:     decimal.NewFromInt(0),
			TotalUsdtDown: decimal.NewFromInt(0),
			CreateTime:    common.JSONTime{Time: time.Now()},
			UpdateTime:    common.JSONTime{Time: time.Now()},
		}
		err = global.GVA_DB.Create(&userAccount).Error
		if err != nil {
			log.Panic(api.NewException(api.UserNotExist))
			return
		}
		account = userAccount
	}
	userDetail.VipLevel = account.VipLevel
	userDetail.Usdt = account.Usdt
	userDetail.TotalUsdtDownA = account.TotalUsdtDownA
	userDetail.TotalUsdtDownB = account.TotalUsdtDownB

	//查询卡牌数量
	var cardNum int
	cardNumErr := global.GVA_DB.Model(&manage.MallOrderItem{}).Where("user_id = ? and release_flag != 0", userToken.UserId).Select("sum(goods_count)").Scan(&cardNum).Error
	if err != nil {
		global.GVA_LOG.Error("查询卡牌数量失败", zap.Error(cardNumErr))
	}
	userDetail.CardNum = cardNum
	//查询卡牌价值
	var cardUsdt decimal.Decimal
	cardUsdtErr := global.GVA_DB.Model(&manage.MallOrderItem{}).Where("user_id = ? and release_flag != 0", userToken.UserId).Select("sum(total_price)").Scan(&cardUsdt).Error
	if err != nil {
		global.GVA_LOG.Error("查询卡牌usdt失败", zap.Error(cardUsdtErr))
	}
	userDetail.CardUsdt = cardUsdt
	//查询父级
	parentId := userInfo.ParentId
	if parentId != 0 {
		var parentUser mall.MallUser
		parentUserErr := global.GVA_DB.Where("user_id =?", parentId).First(&parentUser).Error
		if parentUserErr != nil {
			global.GVA_LOG.Error("查询用户父级对象失败", zap.Error(parentUserErr))
		}
		userDetail.ParentBscAddress = parentUser.BscAddress
	}

	//冻结金额 订单冻结金额+余额明细冻结金额
	//查询用户冻结余额
	var usdtFreeze decimal.Decimal
	usdtFreezeErr := global.GVA_DB.Model(&manage.MallOrderItem{}).Where("user_id = ? and release_flag = 1", userToken.UserId).Select("sum(usdt_able)").Scan(&usdtFreeze).Error
	if usdtFreezeErr != nil {
		global.GVA_LOG.Error("查询用户订单明细冻结余额失败", zap.Error(usdtFreezeErr))
	}
	var usdtAccountFreeze decimal.Decimal
	usdtAccountFreezeErr := global.GVA_DB.Model(&bscscan.BscMallAccountDetail{}).Where("user_id = ? and release_flag = 1", userToken.UserId).Select("sum(usdt_able)").Scan(&usdtAccountFreeze).Error
	if usdtAccountFreezeErr != nil {
		global.GVA_LOG.Error("查询用户账户明细冻结余额失败", zap.Error(usdtAccountFreezeErr))
	}
	userDetail.UsdtFreeze = usdtFreeze.Add(usdtAccountFreeze)
	userDetail.BonusFlag = userInfo.BonusFlag
	return
}

func (m *MallUserService) GetUserTeamList(pageNumber int, token string) (err error, list interface{}, total int64) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), list, total
	}
	db := global.GVA_DB.Table("tb_newbee_mall_user as u")
	err = db.Where("u.parent_ids like ?", "%["+strconv.Itoa(userToken.UserId)+"]%").Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("查询团队总数失败", zap.Error(err))
		return errors.New("查询团队总数失败"), list, total
	}
	db.Joins(" left join tb_bsc_mall_user_account a on u.user_id = a.user_id ")
	db.Select("u.bsc_address,u.create_time,a.vip_level,a.total_usdt as usdt,u.user_id")
	limit := 5
	offset := 5 * (pageNumber - 1)
	var userList []response.MallUserDetailResponse
	err = db.Limit(limit).Offset(offset).Order("u.create_time asc").Find(&userList).Error
	if err != nil {
		return errors.New("查询团队失败"), list, total
	}
	return err, userList, total
}

func (m *MallUserService) GetAccountDetailList(pageNumber int, token string) (err error, list interface{}, total int64) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), list, total
	}
	db := global.GVA_DB.Model(&bscscan.BscMallAccountDetail{})
	err = db.Where("user_id = ?", userToken.UserId).Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("查询账户明细总数失败", zap.Error(err))
		return errors.New("查询账户明细总数失败"), list, total
	}
	limit := 5
	offset := 5 * (pageNumber - 1)
	var detailList []bscscan.BscMallAccountDetail
	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&detailList).Error
	if err != nil {
		return errors.New("查询账户明细失败"), list, total
	}
	return err, detailList, total
}

func (m *MallUserService) UserLogin(params mallReq.UserLoginParam) (err error, user mall.MallUser, userToken mall.MallUserToken) {
	err = global.GVA_DB.Where("login_name=? AND password_md5=?", params.LoginName, params.PasswordMd5).First(&user).Error
	if user != (mall.MallUser{}) {
		token := getNewToken(time.Now().UnixNano()/1e6, int(user.UserId))
		global.GVA_DB.Where("user_id", user.UserId).First(&token)
		nowDate := time.Now()
		// 48小时过期
		expireTime, _ := time.ParseDuration("48h")
		expireDate := nowDate.Add(expireTime)
		// 没有token新增，有token 则更新
		if userToken == (mall.MallUserToken{}) {
			userToken.UserId = user.UserId
			userToken.Token = token
			userToken.UpdateTime = nowDate
			userToken.ExpireTime = expireDate
			if err = global.GVA_DB.Save(&userToken).Error; err != nil {
				return
			}
		} else {
			userToken.Token = token
			userToken.UpdateTime = nowDate
			userToken.ExpireTime = expireDate
			if err = global.GVA_DB.Save(&userToken).Error; err != nil {
				return
			}
		}
	}
	return err, user, userToken
}

func (m *MallUserService) UserAddressLogin(params mallReq.UserAddressLoginParam) (err error, user mall.MallUser, userToken mall.MallUserToken) {
	var db *gorm.DB
	db = global.GVA_DB
	err = db.Where("bsc_address=? AND is_deleted=? And login_type=? ", params.BscAddress, 0, params.LoginType).First(&user).Error
	inviteId := params.InviteId
	nodeType := params.NodeType
	fmt.Println("inviteId:", inviteId)
	fmt.Println("nodeType:", nodeType)
	if user != (mall.MallUser{}) {
		//用户已存在
		errGetToken, token := getToken(user.UserId)
		if errGetToken != nil {
			log.Panic(api.NewException(api.UserGetTokenFail))
		}
		userToken = token
		thisParentId := user.ParentId
		fmt.Println("user.UserId:", user.UserId)
		fmt.Println("thisParentId:", thisParentId)
		if thisParentId == 0 && inviteId != 0 && nodeType != "" && inviteId != user.UserId {
			//判断是否循环绑定上级
			flag := m.checkParentId(user.UserId, inviteId)
			if !flag {
				return
			}
			parentId := 0
			parentIds := ""
			if inviteId != 0 && nodeType != "" {
				parentId = m.getParentId(inviteId, nodeType)
				var thisUser mall.MallUser
				thisUserErr := global.GVA_DB.Where("user_id = ?", parentId).First(&thisUser).Error
				if thisUserErr != nil {
					global.GVA_LOG.Error("查询节点记录失败", zap.Error(thisUserErr))
				}
				if thisUser.ParentIds != "" {
					parentIds = thisUser.ParentIds + ",[" + strconv.Itoa(parentId) + "]"
				} else {
					parentIds = "[" + strconv.Itoa(parentId) + "]"
				}
			}
			user.InviteId = inviteId
			user.ParentId = parentId
			user.ParentIds = parentIds
			err = global.GVA_DB.Where("user_id = ?", user.UserId).UpdateColumns(&user).Error
			if err != nil {
				log.Panic(api.NewException(api.AddUserFail))
				return
			}
		}
	} else {
		parentId := 0
		parentIds := ""
		if inviteId != 0 && nodeType != "" {
			parentId = m.getParentId(inviteId, nodeType)
			var thisUser mall.MallUser
			thisUserErr := global.GVA_DB.Where("user_id = ?", parentId).First(&thisUser).Error
			if thisUserErr != nil {
				global.GVA_LOG.Error("查询节点记录失败", zap.Error(thisUserErr))
			}
			if thisUser.ParentIds != "" {
				parentIds = thisUser.ParentIds + ",[" + strconv.Itoa(parentId) + "]"
			} else {
				parentIds = "[" + strconv.Itoa(parentId) + "]"
			}
		}
		// 保存用户数据
		tx := global.GVA_DB.Begin()
		userNew := mall.MallUser{
			BscAddress: params.BscAddress,
			InviteId:   params.InviteId,
			ParentId:   parentId,
			ParentIds:  parentIds,
			NodeType:   params.NodeType,
			LoginType:  params.LoginType,
			CreateTime: common.JSONTime{Time: time.Now()},
		}

		err = global.GVA_DB.Create(&userNew).Error
		user = userNew
		if err != nil {
			tx.Rollback()
			log.Panic(api.NewException(api.AddUserFail))
			return
		}
		userAccount := bscscan.BscMallUserAccount{
			UserId:        userNew.UserId,
			ParentId:      parentId,
			VipLevel:      0,
			Usdt:          decimal.NewFromInt(0),
			UsdtFreeze:    decimal.NewFromInt(0),
			TotalUsdt:     decimal.NewFromInt(0),
			TotalUsdtDown: decimal.NewFromInt(0),
			CreateTime:    common.JSONTime{Time: time.Now()},
			UpdateTime:    common.JSONTime{Time: time.Now()},
		}
		err = global.GVA_DB.Create(&userAccount).Error
		if err != nil {
			tx.Rollback()
			log.Panic(api.NewException(api.AddUserFail))
			return
		}
		tx.Commit()
		err, token := getToken(userNew.UserId)
		if err != nil {
			log.Panic(api.NewException(api.UserGetTokenFail))
		}
		userToken = token
	}
	return err, user, userToken
}

func (m *MallUserService) checkParentId(thisUserId int, inviteId int) (flag bool) {
	//查询邀请的用户
	var inviteUser mall.MallUser
	flag = true
	inviteUserErr := global.GVA_DB.Where("user_id = ?", inviteId).First(&inviteUser).Error
	if inviteUserErr != nil {
		global.GVA_LOG.Error("查询邀请人记录失败", zap.Error(inviteUserErr))
		flag = false
		return
	}
	parentId := inviteUser.ParentId
	if parentId == thisUserId {
		flag = false
		return
	}
	thisInviteId := inviteUser.InviteId
	if thisInviteId == thisUserId {
		flag = false
		return
	}
	var bf bytes.Buffer
	bf.WriteString("[")
	bf.WriteString(strconv.Itoa(thisUserId))
	bf.WriteString("]")
	parentIds := inviteUser.ParentIds
	if strings.Contains(parentIds, bf.String()) {
		flag = false
		return
	}
	return
}

func (m *MallUserService) getParentId(inviteId int, nodeType string) (parentId int) {
	if inviteId != 0 && nodeType != "" {
		//查询当前用户的父节点是哪个
		var subUser mall.MallUser
		subUserErr := global.GVA_DB.Where("parent_id = ? and node_type = ?", inviteId, nodeType).First(&subUser).Error
		if subUserErr != nil {
			global.GVA_LOG.Error("查询子节点记录失败", zap.Error(subUserErr))
		}
		//当前子节点用户为空
		if subUser == (mall.MallUser{}) {
			parentId = inviteId
		} else {
			thisParentId := m.getParentId(subUser.UserId, nodeType)
			if thisParentId != 0 {
				parentId = thisParentId
				return
			}
		}
	}
	return
}

func getToken(userId int) (err error, userToken mall.MallUserToken) {
	token := getNewToken(time.Now().UnixNano()/1e6, int(userId))
	global.GVA_DB.Where("user_id", userId).First(&token)
	nowDate := time.Now()
	// 48小时过期
	expireTime, _ := time.ParseDuration("48h")
	expireDate := nowDate.Add(expireTime)
	// 没有token新增，有token 则更新
	if userToken == (mall.MallUserToken{}) {
		userToken.UserId = userId
		userToken.Token = token
		userToken.UpdateTime = nowDate
		userToken.ExpireTime = expireDate
		if err = global.GVA_DB.Save(&userToken).Error; err != nil {
			return
		}
	} else {
		userToken.Token = token
		userToken.UpdateTime = nowDate
		userToken.ExpireTime = expireDate
		if err = global.GVA_DB.Save(&userToken).Error; err != nil {
			return
		}
	}
	return err, userToken
}

func getNewToken(timeInt int64, userId int) (token string) {
	var build strings.Builder
	build.WriteString(strconv.FormatInt(timeInt, 10))
	build.WriteString(strconv.Itoa(userId))
	build.WriteString(utils.GenValidateCode(6))
	return utils.MD5V([]byte(build.String()))
}
