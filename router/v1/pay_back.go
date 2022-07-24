package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/global"
	"loante/service/model"
	"loante/service/payments"
	"loante/service/resp"
)

//接收支付回调
type payBack struct {}

func NewPayBack() *payBack {
	return new(payBack)
}
//PayNotify 收款回调
func (p *payBack)PayNotify(c *fiber.Ctx) error {
	//获取商户订单编号
	global.Log.Info("PayNotify: %s", string(c.Body()))
	MchOrderId := payments.MchOrderId(c)
	global.Log.Info("MchOrderId: %s", MchOrderId)
	if MchOrderId == ""{
		global.Log.Info("未识别出商户订单编号")
		return  resp.Err(c, 1, "未找到账户订单")
	}
	//根据商户编号查询配置
	orders := new(model.Orders)
	orders.One(fmt.Sprintf("payment_request_no = '%s'", MchOrderId))
	if orders.Payment == 0{
		global.Log.Info("未找到支付通道")
		return  resp.Err(c, 1, "未找到支付通道")
	}
	//避免重复回调
	if orders.RepaidStatus == 1{
		global.Log.Info("重复回调")
		return  resp.Err(c, 1, "重复回调")
	}
	payment := new(model.Payment)
	payment.One(fmt.Sprintf("id = '%d'", orders.Payment))
	//查询
	pp := new(model.ProductPayment)
	pp.One(fmt.Sprintf("payment_id = %d and product_id = %d", orders.Payment, orders.ProductId))
	global.Log.Info("参数：%s  %s", payment.Name, pp.Configuration)
	payModel :=  payments.SelectPay(payment.Name)
	ret, amount, err := (*payModel).Verify(pp.Configuration, c)
	global.Log.Info("结果：%v %v", ret, amount)
	if !ret{
		global.Log.Info("结果：err = %v", err.Error())
		return  resp.Err(c, 1, "支付通道签名验证不正确")
	}
	global.Log.Info("处理收款以后的订单逻辑")
	//处理订单的逻辑
	orders.PayAfter(amount)
	return  resp.OK(c, "")
}

//OutNotify 放款回调
func (p *payBack)OutNotify(c *fiber.Ctx) error {
	//获取商户订单编号
	global.Log.Info("OutNotify: %s", string(c.Body()))
	MchOrderId := payments.MchOrderId(c)
	global.Log.Info("MchOrderId: %s", MchOrderId)
	if MchOrderId == ""{
		global.Log.Info("未识别出商户订单编号")
		return  resp.Err(c, 1, "未识别出商户订单编号")
	}
	//根据商户编号查询配置
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("payment_request_no = '%s'", MchOrderId))
	if borrowData.Payment == 0{
		global.Log.Info("未找到支付通道")
		return  resp.Err(c, 1, "未找到支付通道")
	}
	payment := new(model.Payment)
	payment.One(fmt.Sprintf("id = '%d'", borrowData.Payment))
	//查询
	pp := new(model.ProductPayment)
	pp.One(fmt.Sprintf("payment_id = %d and product_id = %d", borrowData.Payment, borrowData.ProductId))
	payModel :=  payments.SelectPay(payment.Name)
	ret, _, err := (*payModel).Verify(pp.Configuration ,c)
	if !ret{
		global.Log.Info("结果：err = %v", err.Error())
		return  resp.Err(c, 1, "支付通道签名验证不正确")
	}
	//处理订单的逻辑
	borrowData.PayAfter()
	return  resp.OK(c, "")
}