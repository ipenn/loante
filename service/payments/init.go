package payments

import (
	"github.com/gofiber/fiber/v2"
)

type Payments interface {
	PayIn(config string, pays *Pays) (bool, map[string]interface{}, error) //收款
	PayOut(config string, pays *Pays) (bool, error) //放款
	Verify(config string, ctx *fiber.Ctx) (bool,float64, error) //验证签名代收
}

type Pays struct {
	OrderId string //商户订单号
	Amount float64	//金额
	CustomName string	//客户姓名
	CustomMobile string //客户电话
	CustomEmail string	//客户email地址
	BankAccount string	//收款人银行账号
	IfscCode string		//收款人IFSC CODE
	Remark string		//备注
	NotifyUrl string //异步通知回调地址
	CallbackUrl string //页面回跳地址
	PlatOrderId string //平台相应的编号
}
type ReturnPay struct {
	OrderId string
	Amount float64
	Ret bool
}

func SelectPay(name string) *Payments {
	var pp = Payments(nil)
	if name == "T-Pays"{
		pp = new(TPays)
	}else if name == "HXPay"{
		pp = new(HXPays)
	}else if name == "WhalePay"{
		pp = new(WhalePay)
	}
	return &pp
}

func MchOrderId(ctx *fiber.Ctx) string {
	body := TPayNotify{}
	ctx.BodyParser(body)
	if len(body.OrderId) > 0{
		return body.OrderId
	}
	body2 := HXPayNotify{}
	ctx.BodyParser(body2)
	if len(body2.MerchantCode) > 0{
		return body2.MerchantCode
	}
	return ""
}

//func VerifySign(ctx *fiber.Ctx) ReturnPay {
//	//解析 TPays
//	var pp = Payments(nil)
//	body := TPayNotify{}
//	ctx.BodyParser(body)
//	m2 := map[string]interface{}{}
//	sign := ""
//	retPay := ReturnPay{}
//	if len(body.PlatOrderId) > 0{
//		pp = new(TPays)
//		sign = body.Sign
//		body.Sign = ""
//		m2, _ = tools.StructToMapReflect(&body,"json")
//		if ret, _ := pp.VerifySign(m2, sign);ret{
//			retPay.Amount,_ = strconv.ParseFloat(body.Amount, 10)
//			retPay.OrderId = body.OrderId
//			retPay.Ret = true
//		}
//		return retPay
//	}
//	body2 := HXPayNotify{}
//	ctx.BodyParser(body2)
//	if len(body2.Sign) > 0{
//		pp = new(TPays)
//		sign = body2.Sign
//		body2.Sign = ""
//		m2, _ = tools.StructToMapReflect(&body,"json")
//		if ret, _ := pp.VerifySign(m2, sign); ret {
//			retPay.OrderId = body2.MerchantCode
//			retPay.Amount = body2.PaidAmount
//			retPay.Ret = true
//		}
//		return retPay
//	}
//	return retPay
//}