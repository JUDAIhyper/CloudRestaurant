package param

//手机号+短信验证码登录的解析参数
type SmsLoginParam struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
