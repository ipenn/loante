package payments

type Payments interface {
	PayIn(config string, pays Pays) (bool, interface{}, error) //收款
	PayOut(config string, pays Pays) (bool, error) //放款
	VerifySign(data map[string]interface{}, sign string) (bool, error) //验证签名
}

type Pays struct {
	OrderId string //商户订单号
	Amount float64	//金额
	CustomName string	//客户姓名
	CustomMobile string //客户电话
	CustomEmail string	//客户email地址
	BankAccount string	//收款人银行账号
	IfscCode string		//收款人IFSC CODE
	UpiAccount string		//收款人UPI账户
	Remark string		//备注
	NotifyUrl string //异步通知回调地址
	CallbackUrl string //页面回跳地址
}