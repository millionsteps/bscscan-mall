package api

const (
	AddUserFail        = 10001
	UserMobileNotExist = 10002
	UserNotExist       = 10003
	UserGetTokenFail   = 10004

	ActivityNotExist = 20001

	ObjectCartItemExist    = 30001
	ObjectCartItemNotExist = 30002
	ObjectCartItemNotEdit  = 30003
	ObjectStatusFail       = 30004
	UserObjectCanSale      = 30005
	UserObjectCanNotSend   = 30006

	ProxyServeHasAudit          = 40001
	ProxyServeNoPayEarnestMoney = 40002
	ServeStatusError            = 40003
	ServeNotExist               = 40004

	ObjectNoSerialNumber = 50001
)

var resultCodeText = map[int]string{
	AddUserFail:                 "添加用户失败",
	UserMobileNotExist:          "请先绑定手机号码！",
	UserNotExist:                "用户不存在！",
	UserGetTokenFail:            "获取用户token失败！",
	ActivityNotExist:            "活动不存在",
	ObjectCartItemExist:         "商品已存在！",
	ObjectCartItemNotExist:      "记录不存在！",
	ObjectCartItemNotEdit:       "已选择版号无法修改",
	ObjectStatusFail:            "收藏失败",
	UserObjectCanSale:           "藏品无法出售",
	UserObjectCanNotSend:        "藏品无法发货",
	ProxyServeHasAudit:          "已审核",
	ProxyServeNoPayEarnestMoney: "服务未支付保证金",
	ServeStatusError:            "服务状态有问题",
	ObjectNoSerialNumber:        "无版号可选",
	ServeNotExist:               "服务不存在",
}

func StatusText(code int) (string, bool) {
	message, ok := resultCodeText[code]
	return message, ok
}
