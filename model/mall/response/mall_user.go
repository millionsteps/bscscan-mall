package response

type MallUserDetailResponse struct {
	NickName      string `json:"nickName"`
	LoginName     string `json:"loginName"`
	IntroduceSign string `json:"introduceSign"`
	UserId        int    `json:"userId"`
	BscAddress    string `json:"bscAddress"`
	LoginType     int    `json:"loginType"`
}
