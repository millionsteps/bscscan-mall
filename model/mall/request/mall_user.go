package request

//用户注册
type RegisterUserParam struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}

//更新用户信息
type UpdateUserInfoParam struct {
	NickName      string `json:"nickName"`
	PasswordMd5   string `json:"passwordMd5"`
	IntroduceSign string `json:"introduceSign"`
	EmailAddress  string `json:"emailAddress"`
}

type UserLoginParam struct {
	LoginName   string `json:"loginName"`
	PasswordMd5 string `json:"passwordMd5"`
}

type UserAddressLoginParam struct {
	BscAddress string `json:"bscAddress"`
	LoginType  int    `json:"loginType"`
	InviteId   int    `json:"inviteId"`
	NodeType   string `json:"nodeType"`
}
