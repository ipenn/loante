package payments

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req"
	"loante/global"
	"loante/tools"
	"sort"
	"strconv"
	"strings"
)

const (
	WhalePayUrl    = "https://api.hgcjbj.com"
	WhaleInUrl  = "/openApi/pay/create"
	WhaleOutUrl = "/payout/create"
)

type WhalePay struct {
	Merchant string `json:"merchant"`
	Key      string `json:"key"`
}
type WhalePayInSuc struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Merchant    string `json:"merchant"`
		OrderId     string `json:"orderId"`
		PlatOrderId string `json:"platOrderId"`
		Sign        string `json:"sign"`
		Url         string `json:"url"`
	} `json:"data,omitempty"`
}

type WhalePayOutSuc struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Amount      string `json:"amount"`
		OrderId     string `json:"orderId"`
		Sign        string `json:"sign"`
		Merchant    string `json:"merchant"`
		PlatOrderId string `json:"platOrderId"`
		Status      string `json:"status"`
	} `json:"data,omitempty"`
}

type WhalePayNotify struct {
	Merchant    string `json:"merchant"`
	PlatOrderId string `json:"platOrderId"`
	OrderId     string `json:"orderId"`
	Amount      string `json:"amount"`
	Status      string `json:"status"` //订单状态 :0已创建 1支付中 2 支付失败 3 支付成功
	Msg         string `json:"msg"`
	Sign        string `json:"sign"`
}

func (t *WhalePay) buildUrl(data req.Param) string {
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
	return strings.Join(query, `&`) + "&key=" + t.Key
}

func (t *WhalePay) init(config string) {
	err := json.Unmarshal([]byte(config), t)
	if err != nil {
		global.Log.Error(err.Error())
	}
}

func (t *WhalePay) sign(data req.Param) string {
	signTemp := t.buildUrl(data)
	fmt.Println(signTemp)
	signData := tools.Md5(signTemp)
	return strings.ToLower(signData)
}

func (t *WhalePay)Verify(config string, ctx *fiber.Ctx) (bool,float64, error) {
	t.init(config)
	pa := req.Param{}
	body := WhalePayNotify{}
	ctx.BodyParser(body)
	m2, _ := tools.StructToMapReflect(&body,"json")
	for index, item := range m2{
		if index != "sign"{
			pa[index] = item
		}
	}
	if t.sign(pa) != body.Sign{
		return false, 0,errors.New("sign fail")
	}
	if body.Status != "3"{
		return false,0 ,errors.New(body.Status)
	}
	amount,_  := strconv.ParseFloat(body.Amount, 10)
	return true, amount,nil
}

func (t *WhalePay)PayIn(config string, pays *Pays) (bool, map[string]interface{}, error) {
	t.init(config)
	data := req.Param{
		"merchant":     t.Merchant,
		"orderId":      pays.OrderId,
		"amount":       pays.Amount,
		"customName":   pays.CustomName,
		"customMobile": pays.CustomMobile,
		"customEmail":  pays.CustomEmail,
		"notifyUrl":    pays.NotifyUrl,
		"callbackUrl":  pays.CallbackUrl,
	}
	data["sign"] = t.sign(data)
	fmt.Println(data["sign"])
	resp, err := req.Post(WhalePayUrl+WhaleInUrl, data)
	if err != nil {
		return false, nil, err
	}
	fmt.Println(resp.String())
	res := WhalePayInSuc{}
	if err := resp.ToJSON(&res); err != nil {
		return false, nil, err
	}
	if res.Code != 1 {
		return false, nil, errors.New(res.Msg)
	}
	pays.PlatOrderId = res.Data.PlatOrderId
	return true, map[string]interface{}{
		"platId":res.Data.PlatOrderId,
		"orderId":res.Data.OrderId,
		"url":res.Data.Url,
	},nil
}

func (t *WhalePay) PayOut(config string, pays *Pays) (bool, error) {
	t.init(config)
	data := req.Param{
		"merchant":     t.Merchant,
		"orderId":      pays.OrderId,
		"amount":       pays.Amount,
		"customName":   pays.CustomName,
		"customMobile": pays.CustomMobile,
		"customEmail":  pays.CustomEmail,
		"bankAccount":  pays.BankAccount,
		"ifscCode":     pays.IfscCode,
		"notifyUrl":    pays.NotifyUrl,
	}
	data["sign"] = t.sign(data)
	resp, err := req.Post(WhalePayUrl+WhaleOutUrl, req.Header{
		"Content-Type": "application/x-www-form-urlencoded",
	}, data)
	if err != nil {
		return false, err
	}
	res := WhalePayOutSuc{}
	if err := resp.ToJSON(&res); err != nil {
		return false, err
	}
	if res.Code != 200 {
		return false, errors.New(res.Msg)
	}
	pays.PlatOrderId = res.Data.PlatOrderId
	return true, nil
}


