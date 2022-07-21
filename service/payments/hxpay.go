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
	"strings"
)

const (
	HXPayUrl    = "https://hxpayment.top"
	HXPayInUrl  = "/payment/collection"
	HXPayOutUrl = "/payment/payout"
)

type HXPays struct {
	Merchant string `json:"merchant"`
	Key      string `json:"key"`
}

type HXPayInSuc struct {
	PlatformOrderCode string `json:"platformOrderCode"`
	PaymentUrl        string `json:"paymentUrl"`
}

type HXPayOutSuc struct {
	PlatformOrderCode string `json:"platformOrderCode"`
}

type HXPayErr struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Detail  string `json:"detail"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type HXPayNotify struct {
	MerchantLogin      string `json:"merchantLogin"`
	OrderCode    string `json:"orderCode"`
	MerchantCode         string `json:"merchantCode"`
	Status     string `json:"status"`
	OrderAmount int `json:"orderAmount,omitempty"`
	PaidAmount  float64 `json:"paidAmount,omitempty"`
	Sign      string `json:"sign"`
}

func (hx *HXPays) buildUrl(data req.Param) string {
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
	return strings.Join(query, `&`) + "&key=" + hx.Key
}

func (hx *HXPays) init(config string) {
	err := json.Unmarshal([]byte(config), hx)
	if err != nil {
		global.Log.Error(err.Error())
	}
}

func (hx *HXPays) sign(data req.Param) string {
	signTemp := hx.buildUrl(data)
	fmt.Println(signTemp)
	signData := tools.Md5(signTemp)
	return strings.ToLower(signData)
}

func (hx *HXPays)Verify(config string, ctx *fiber.Ctx) (bool,float64, error) {
	hx.init(config)
	pa := req.Param{}
	body := HXPayNotify{}
	ctx.BodyParser(body)
	m2, _ := tools.StructToMapReflect(&body,"json")
	for index, item := range m2{
		if index == "sign"{
			pa[index] = item
		}
	}
	if hx.sign(pa) != body.Sign{
		return false, 0, errors.New("sign fail")
	}
	if body.Status != "SUCCESS"{
		return false, 0, errors.New(body.Status)
	}
	return true, body.PaidAmount, nil
}

func (hx *HXPays)PayIn(config string, pays *Pays) (bool, map[string]interface{}, error) {
	hx.init(config)
	data := req.Param{
		"merchantLogin": hx.Merchant,
		"orderCode":     pays.OrderId,
		"amount":        fmt.Sprintf("%.2f", pays.Amount),
		"name":          pays.CustomName,
		"phone":         pays.CustomMobile,
		"email":         pays.CustomEmail,
		"remark":        pays.Remark,
	}
	data["sign"] = hx.sign(data)
	fmt.Println(data["sign"])
	bdata, err := json.Marshal(&data)
	resp, err := req.Post(HXPayUrl+HXPayInUrl, req.Header{
		"Content-Type": "application/json",
	}, bdata)
	if err != nil {
		return false, nil, err
	}
	fmt.Println(resp.String())
	res := HXPayInSuc{}
	if err := resp.ToJSON(&res); err != nil {
		return false, nil, err
	}
	if res.PaymentUrl == "" {
		res2 := HXPayErr{}
		resp.ToJSON(&res2)
		return false, nil, errors.New(res2.Detail)
	}
	pays.PlatOrderId = res.PlatformOrderCode
	return true, map[string]interface{}{
		"platId":res.PlatformOrderCode, //平仓的
		"orderId":pays.OrderId,
		"url":res.PaymentUrl,
	},nil
}

func (hx *HXPays) PayOut(config string, pays *Pays) (bool, error) {
	hx.init(config)
	data := req.Param{
		"merchantLogin": hx.Merchant,
		"orderCode":     pays.OrderId,
		"amount":        pays.Amount,
		"name":          pays.CustomName,
		"account":       pays.BankAccount,
		"ifsc":          pays.IfscCode,
		"remark":        pays.Remark,
		"notifyUrl":     pays.NotifyUrl,
	}
	data["sign"] = hx.sign(data)
	bdata, err := json.Marshal(&data)
	resp, err := req.Post(TPayUrl+TPayOutUrl, req.Header{
		"Content-Type": "application/json",
	}, bdata)
	if err != nil {
		return false, err
	}
	res := HXPayOutSuc{}
	if err := resp.ToJSON(&res); err != nil {
		return false, err
	}
	if res.PlatformOrderCode == "" {
		res2 := HXPayErr{}
		resp.ToJSON(&res2)
		return false, errors.New(res2.Detail)
	}
	pays.PlatOrderId = res.PlatformOrderCode
	return true, nil
}
