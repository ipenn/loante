package sms

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
	"loante/global"
)

type NxCloud struct {
	AppKey 		string 	`json:"appkey"`
	SecretKey	string 	`json:"secretkey"`
}

type NxCloudSendResp struct {
	Result    string `json:"result"`
	Messageid string `json:"messageid"`
	Code      string `json:"code"`
}

const NXURL = "http://api2.nxcloud.com"

//Send 发送
func (x *NxCloud)Send(phone, content string, smsType int) (bool,error) {
	url := fmt.Sprintf("%s%s", NXURL, "/api/sms/mtsend")
	key := "v3tIfD7v"
	secret := "qBV5wN7L"
	if smsType > 1{
		key = "Jx8Ubn2v"
		secret = "SK0jHWxG"
	}
	resp, err := req.Post(url,req.Header{
		"Content-Type":"application/x-www-form-urlencoded",
	}, req.Param{
		"appkey": key,
		"secretkey": secret,
		"phone":phone,
		"content":content,
	})
	if err != nil{
		return false,err
	}
	global.Log.Info("sms send content: " + resp.String())
	ret := NxCloudSendResp{}
	resp.ToJSON(&ret)
	if ret.Code != "0"{
		return false,errors.New(ret.Result)
	}
	return true, nil
}
