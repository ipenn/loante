package test

import (
	"fmt"
	"loante/service/payments"
	"testing"
)

func TestTPayIn(t *testing.T) {
	//准备数据
	paysData := payments.Pays{
		OrderId:      "nxsanxsa",
		Amount:       6000,
		CustomName:   "vlaxsa",
		CustomMobile: "1478542134",
		CustomEmail:  "66897@qq.com",
		BankAccount:  "12345676",
		IfscCode:     "1222",
		NotifyUrl:    "http://127.0.0.1/notify",
		CallbackUrl:  "http://127.0.0.1/notify",
	}
	pay := payments.TPays{}
	ret, res, err := pay.PayIn("{\"merchant\":\"C1657530632421\",\"key\":\"85wDfJRBx6YNvPyw0Rdl3N00n5REN3ylLtkSjS\"}", &paysData)
	fmt.Println(ret)
	fmt.Println(res)
	fmt.Println(err)
}
func TestHxPayIn(t *testing.T) {
	//准备数据
	paysData := payments.Pays{
		OrderId:      "nxsanxsa",
		Amount:       6000.00,
		CustomName:   "vlaxsa",
		CustomMobile: "1478542134",
		CustomEmail:  "66897@qq.com",
		BankAccount:  "12345676",
		IfscCode:     "1222",
		Remark:       "1222",
		NotifyUrl:    "http://127.0.0.1/notify",
		CallbackUrl:  "http://127.0.0.1/notify",
	}
	pay := payments.HXPays{}
	ret, res, err := pay.PayIn("{\"merchant\":\"Loante02\",\"key\":\"oaDO6tAZdfeExRnlGjrJ\"}", &paysData)
	fmt.Println(ret)
	fmt.Println(res)
	fmt.Println(err)
}

func TestHxPayInBack(t *testing.T) {
	//准备数据
	pay := payments.HXPays{}
	v1, v2, v3 := pay.Verify2("{\"merchant\":\"Loante01\",\"key_in\":\"oaDO6tAZdfeExRnlGjrJ\",\"key_out\":\"AN1lByrZ3p9x8wUReq6I\",\"remark\":\"qwe\"}","{\"merchantLogin\":\"Loante01\",\"orderCode\":\"C20220722131454777480\",\"merchantCode\":\"QWCFlh-1-1-2\",\"status\":\"SUCCESS\",\"orderAmount\":\"420\",\"paidAmount\":\"420\",\"sign\":\"8ac00bde28c43c043c28950fe29b0533\"}")
	fmt.Println("start...")
	fmt.Println(v1)
	fmt.Println(v2)
	fmt.Println(v3)
	fmt.Println("end...")
}
