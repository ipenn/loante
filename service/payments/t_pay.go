package payments

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"loante/global"
	"loante/tools"
	"sort"
	"strings"
)

const (
	TPayUrl = "https://www.tpays.in"
	TPayInUrl = "/openApi/pay/createOrder"
	TPayOutUrl = "/openApi/payout/createOrder"
)

type TPays struct {
	Merchant string `json:"merchant"`
	Key string `json:"key"`
}

type TPayNotify struct {
	Amount      string `json:"amount"`
	Merchant    string `json:"merchant"`
	Msg         string `json:"msg"`
	OrderId     string `json:"orderId"`
	PlatOrderId string `json:"platOrderId"`
	Sign        string `json:"sign"`
	Status      string `json:"status"`
}

type TPayInSuc struct {
	Code int `json:"code"`
	Data struct {
		Merchant    string `json:"merchant"`
		OrderId     string `json:"orderId"`
		PlatOrderId string `json:"platOrderId"`
		Sign        string `json:"sign"`
		Url         string `json:"url"` //要返回给客户的
	} `json:"data"`
	Success bool `json:"success"`
}

type TPayOutSuc struct {
	Code int `json:"code"`
	Data struct {
		Amount      string `json:"amount"`
		Merchant    string `json:"merchant"`
		Msg         string `json:"msg"`
		OrderId     string `json:"orderId"`
		PlatOrderId string `json:"platOrderId"`
		Sign        string `json:"sign"`
		Status      string `json:"status"`
	} `json:"data"`
	Success bool `json:"success"`
}

type TPayErr struct {
	Code          int    `json:"code"`
	ErrorMessages string `json:"errorMessages"`
	Success       bool   `json:"success"`
}


func (t *TPays)buildUrl(data req.Param) string {
	var (
		keys  []string
		query []string
	)
	for index, _ := range data {
		keys = append(keys, index)
	}
	sort.Strings(keys)
	for _, k := range keys {
		query = append(query, fmt.Sprintf("%s=%v", k, data[k]))
	}
	return strings.Join(query, `&`)  + "&key=" + t.Key
}
func (t *TPays)init(config string)  {
	err := json.Unmarshal([]byte(config), t)
	if err != nil{
		global.Log.Error(err.Error())
	}
}
func (t *TPays)sign(data req.Param) string {
	signTemp := t.buildUrl(data)
	fmt.Println(signTemp)
	signData := tools.Md5(signTemp)
	return strings.ToLower(signData)
}

func (t *TPays)VerifySign(data map[string]interface{}, sign string) (bool, error) {
	pa := req.Param{}
	for index, item := range data{
		pa[index] = item
	}
	if t.sign(pa) != sign{
		return false,errors.New("sign fail")
	}
	return true,nil
}

func (t *TPays)PayIn(config string, pays Pays) (bool, interface{}, error) {
	t.init(config)
	data := req.Param{
		"merchant":t.Merchant,
		"orderId":pays.OrderId,
		"amount":pays.Amount,
		"customName":pays.CustomName,
		"customMobile":pays.CustomMobile,
		"customEmail":pays.CustomEmail,
		"notifyUrl":pays.NotifyUrl,
		"callbackUrl":pays.CallbackUrl,
	}
	data["sign"] = t.sign(data)
	fmt.Println(data["sign"])
	resp, err := req.Post(TPayUrl + TPayInUrl, data)
	if err != nil{
		return false, nil,err
	}
	fmt.Println(resp.String())
	res := TPayInSuc{}
	if err := resp.ToJSON(&res); err!= nil{
		return false, nil,err
	}
	if res.Code != 200{
		res2 := TPayErr{}
		resp.ToJSON(&res2)
		return false, nil, errors.New(res2.ErrorMessages)
	}
	return true, res,nil
}

func (t *TPays)PayOut(config string, pays Pays) (bool, error) {
	t.init(config)
	data := req.Param{
		"merchant":t.Merchant,
		"orderId":pays.OrderId,
		"amount":pays.Amount,
		"customName":pays.CustomName,
		"customMobile":pays.CustomMobile,
		"customEmail":pays.CustomEmail,
		"mode":"IMPS",
		"bankAccount":pays.BankAccount,
		"ifscCode":pays.IfscCode,
		"upiAccount":pays.IfscCode,
		"notifyUrl":pays.NotifyUrl,
	}
	data["sign"] = t.sign(data)
	resp, err := req.Post(TPayUrl + TPayOutUrl, req.Header{
		"Content-Type":"application/x-www-form-urlencoded",
	}, data)
	if err != nil{
		return false,err
	}
	res := TPayOutSuc{}
	if err := resp.ToJSON(&res); err!= nil{
		return false,err
	}
	if res.Code != 200{
		res2 := TPayErr{}
		resp.ToJSON(&res2)
		return false, errors.New(res2.ErrorMessages)
	}
	return true,nil
}