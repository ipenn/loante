package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type MerchantLittle struct {
	bun.BaseModel `bun:"table:merchant,alias:m"`
	Id            int    `json:"id" bun:",pk"`
	Name          string `json:"name"`
}
type Merchant struct {
	bun.BaseModel `bun:"table:merchant,alias:m"`
	MerchantLittle
	Type          int     `json:"type"`
	CnyBalance    float64 `json:"cny_balance"`
	UsdBalance    float64 `json:"usd_balance"`
	CnyCredit     float64 `json:"cny_credit"`
	UsdCredit     float64 `json:"usd_credit"`
	ContactName   string  `json:"contact_name"`
	ContactMobile string  `json:"contact_mobile"`
	ContactEmail  string  `json:"contact_email"`
	CreateTime    string  `json:"create_time"`
	UpdateTime    string  `json:"update_time"`
	Status        int     `json:"status"`
}

func (a *Merchant) Insert() {
	a.CreateTime = tools.GetFormatTime()
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Merchant) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Merchant) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Merchant) Page(where string, page, limit int) ([]Merchant, int) {
	var datas []Merchant
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Merchant) SetColumn(col string, value interface{}, where string) error {
	_, err := global.C.DB.NewUpdate().Model(a).SetColumn(col, "?", value).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}

func (a *Merchant) Set(col string, where string) error {
	_, err := global.C.DB.NewUpdate().Model(a).Set(col).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}

//AddService 服务消费  <count 件数> <t 类型 = 1 短信 2=风控服务费 >
func (a *Merchant)AddService(count, t int)  {
	//获取费用标准
	sPrice := ServicePrice{}
	sPrice.One(fmt.Sprintf("service_type = %d", t))
	amount := sPrice.Price * float64(count)
	if sPrice.DeductType == 2{ //美元
		a.UsdBalance -= amount
	}else{ //人名币
		a.CnyBalance -= amount
	}
	a.Update(fmt.Sprintf("id = %d", a.Id))
	fund := new(MerchantFund) //mch_id,create_time,amount,type,fund_no,path,remark,currency,in_account_no,rate
	fund.MchId = a.Id
	fund.CreateTime = tools.GetFormatTime()
	fund.Amount = amount * -1
	if t == 1{
		fund.Type = 4 //短信扣费
	}else if t == 2{
		fund.Type = 5 //风控服务费
	}
	fund.FundNo = fmt.Sprintf("%s-%d", tools.InviteCode(8), a.Id)
	fund.Remark = ""
	fund.Currency = sPrice.DeductType
	fund.Insert()
}