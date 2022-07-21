package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/payments"
	"loante/service/resp"
)

//接收支付回调
type payBack struct {}

func NewPayBack() *payBack {
	return new(payBack)
}

func (p *payBack)PayNotify(c *fiber.Ctx) error {
	//获取商户订单编号
	MchOrderId := payments.MchOrderId(c)
	if MchOrderId == ""{
		return  resp.Err(c, 1, "未找到账户订单")
	}
	//根据商户编号查询配置
	orders := new(model.Orders)
	orders.One(fmt.Sprintf("payment_request_no = '%s'", MchOrderId))
	if orders.Payment == 0{
		return  resp.Err(c, 1, "未找到支付通道")
	}
	payment := new(model.Payment)
	payment.One(fmt.Sprintf("id = '%d'", orders.Payment))
	//查询
	pp := new(model.ProductPayment)
	pp.One(fmt.Sprintf("payment_id = %d and product_id = %d", orders.Payment, orders.ProductId))
	payModel :=  payments.SelectPay(payment.Name)
	ret, amount, err := (*payModel).Verify(pp.Configuration, c)
	if !ret{
		fmt.Println(err)
		return  resp.Err(c, 1, "支付通道签名验证不正确")
	}
	//处理订单的逻辑
	orders.PayAfter(amount)
	return  resp.OK(c, "")
}

func (p *payBack)OutNotify(c *fiber.Ctx) error {
	//获取商户订单编号
	MchOrderId := payments.MchOrderId(c)
	if MchOrderId == ""{
		return  resp.Err(c, 1, "未找到账户订单")
	}
	//根据商户编号查询配置
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("payment_request_no = '%s'", MchOrderId))
	if borrowData.Payment == 0{
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
		fmt.Println(err)
		return  resp.Err(c, 1, "支付通道签名验证不正确")
	}
	//处理订单的逻辑
	borrowData.PayAfter()
	return  resp.OK(c, "")
}