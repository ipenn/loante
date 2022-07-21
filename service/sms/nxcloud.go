package sms

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
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
func (x *NxCloud)Send(phone, content string) (bool,error) {
	url := fmt.Sprintf("%s%s", NXURL, "/api/sms/mtsend")
	resp, err := req.Post(url,req.Header{
		"Content-Type":"application/x-www-form-urlencoded",
	}, req.Param{
		"appkey":"",
		"secretkey":"",
		"phone":phone,
		"content":content,
	})
	if err != nil{
		return false,err
	}
	ret := NxCloudSendResp{}
	resp.ToJSON(&ret)
	if ret.Code != "0"{
		return false,errors.New(ret.Result)
	}
	return true, nil
}
