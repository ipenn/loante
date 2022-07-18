package test

import (
	"fmt"
	"loante/service/payments"
	"testing"
)

func TestTPayIn(t *testing.T)  {
	//准备数据
	paysData := payments.Pays{
		OrderId:"nxsanxsa",
		Amount:6000,
		CustomName:"vlaxsa",
		CustomMobile:"1478542134",
		CustomEmail:"66897@qq.com",
		BankAccount:"12345676",
		IfscCode:"1222",
		UpiAccount:"1212121",
		NotifyUrl:"http://127.0.0.1/notify",
		CallbackUrl:"http://127.0.0.1/notify",
	}
	pay := payments.TPays{}
	ret, res, err := pay.PayIn("{\"merchant\":\"C1657530632421\",\"key\":\"85wDfJRBx6YNvPyw0Rdl3N00n5REN3ylLtkSjS\"}", paysData)
	fmt.Println(ret)
	fmt.Println(res)
	fmt.Println(err)
}
func TestHxPayIn(t *testing.T)  {
	//准备数据
	paysData := payments.Pays{
		OrderId:"nxsanxsa",
		Amount:6000.00,
		CustomName:"vlaxsa",
		CustomMobile:"1478542134",
		CustomEmail:"66897@qq.com",
		BankAccount:"12345676",
		IfscCode:"1222",
		Remark:"1222",
		UpiAccount:"1212121",
		NotifyUrl:"http://127.0.0.1/notify",
		CallbackUrl:"http://127.0.0.1/notify",
	}
	pay := payments.HXPays{}
	ret, res, err := pay.PayIn("{\"merchant\":\"Loante02\",\"key\":\"oaDO6tAZdfeExRnlGjrJ\"}", paysData)
	fmt.Println(ret)
	fmt.Println(res)
	fmt.Println(err)
}