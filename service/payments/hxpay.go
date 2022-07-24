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
	KeyIn      string `json:"key_in"`
	KeyOut      string `json:"key_out"`
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
	ErrorMessages string `json:"errorMessages"`
}

type HXPayNotify struct {
	MerchantLogin      string `json:"merchantLogin"`
	OrderCode    string `json:"orderCode"`
	MerchantCode         string `json:"merchantCode"`
	Status     string `json:"status"`
	OrderAmount string `json:"orderAmount"`
	PaidAmount  string `json:"paidAmount"`
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
	return strings.Join(query, `&`)
}

func (hx *HXPays) init(config string) {
	err := json.Unmarshal([]byte(config), hx)
	if err != nil {
		global.Log.Error(err.Error())
	}
}

func (hx *HXPays) sign(data req.Param, key string) string {
	signTemp := hx.buildUrl(data)+ "&key=" + key
	fmt.Println(signTemp)
	signData := tools.Md5(signTemp)
	return strings.ToLower(signData)
}

func (hx *HXPays)Verify2(config string, data string) (bool,float64, error) {
	hx.init(config)
	pa := req.Param{}
	body := HXPayNotify{}
	json.Unmarshal([]byte(data), &body)
	m2, _ := tools.StructToMapReflect(&body,"json")
	for index, item := range m2{
		if index != "sign"{
			pa[index] = item
		}
	}
	if hx.sign(pa, hx.KeyIn) == body.Sign && body.Status == "SUCCESS"{
		return true, tools.ToFloat64(body.PaidAmount), nil
	}

	if hx.sign(pa, hx.KeyOut) == body.Sign && body.Status == "SUCCESS"{
		return true,  tools.ToFloat64(body.PaidAmount), nil
	}
	if body.Status != "SUCCESS"{
		return false, 0, errors.New(body.Status)
	}
	return false, 0, errors.New("sign fail")
}
func (hx *HXPays)Verify(config string, ctx *fiber.Ctx) (bool,float64, error) {
	hx.init(config)
	pa := req.Param{}
	body := new(HXPayNotify)
	ctx.BodyParser(body)
	//json.Unmarshal([]byte("{\"merchantLogin\":\"Loante01\",\"orderCode\":\"C20220722131454777480\",\"merchantCode\":\"QWCFlh-1-1-2\",\"status\":\"SUCCESS\",\"orderAmount\":\"420\",\"paidAmount\":\"420\",\"sign\":\"8ac00bde28c43c043c28950fe29b0533\"}"), &body)
	m2, _ := tools.StructToMapReflect(body,"json")
	for index, item := range m2{
		if index != "sign"{
			pa[index] = item
		}
	}
	if hx.sign(pa, hx.KeyIn) == body.Sign && body.Status == "SUCCESS"{
		return true, tools.ToFloat64(body.PaidAmount), nil
	}
	if hx.sign(pa, hx.KeyOut) == body.Sign && body.Status == "SUCCESS"{
		return true, tools.ToFloat64(body.PaidAmount), nil
	}
	if body.Status != "SUCCESS"{
		return false, 0, errors.New(body.Status)
	}
	return false, 0, errors.New("sign fail")
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
	data["sign"] = hx.sign(data, hx.KeyIn)
	fmt.Println(data["sign"])
	bdata, err := json.Marshal(&data)
	resp, err := req.Post(HXPayUrl+HXPayInUrl, req.Header{
		"Content-Type": "application/json",
	}, bdata)
	if err != nil {
		return false, nil, err
	}
	global.Log.Info(resp.String())
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
		//"amount":        0,
		"name":          pays.CustomName,
		"account":       pays.BankAccount,
		"ifsc":          pays.IfscCode,
		"remark":        pays.Remark,
	}
	data["sign"] = hx.sign(data, hx.KeyOut)
	bdata, err := json.Marshal(&data)
	resp, err := req.Post(HXPayUrl+HXPayOutUrl, req.Header{
		"Content-Type": "application/json",
	}, bdata)
	if err != nil {
		return false, err
	}
	global.Log.Info(resp.String())
	res := HXPayOutSuc{}
	if err := resp.ToJSON(&res); err != nil {
		return false, err
	}
	if res.PlatformOrderCode == "" {
		res2 := HXPayErr{}
		resp.ToJSON(&res2)
		return false, errors.New(res2.ErrorMessages + res2.Detail)
	}
	pays.PlatOrderId = res.PlatformOrderCode
	return true, nil
}
