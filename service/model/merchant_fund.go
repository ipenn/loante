package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type MerchantFund struct {
	bun.BaseModel `bun:"table:merchant_fund,alias:mf"`
	Id            int      `json:"id" bun:",pk"`
	MchId         int      `json:"mch_id"`
	CreateTime    string   `json:"create_time"`
	Amount        float64  `json:"amount"`
	Type          int      `json:"type"`
	FundNo        string   `json:"fund_no"`
	Path          string   `json:"path"`
	Currency      int      `json:"remark"`
	Rate          float64  `json:"rate"`
	Remark        string   `json:"remark"`
	InAccountNo   string   `json:"in_account_no"`
	Merchant      Merchant `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
}

func (a *MerchantFund) Insert() error {
	a.CreateTime = tools.GetFormatTime()
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	//往 balance 添加 cny_balance,usd_balance
	key := "cny_balance = cny_balance + %.2f"
	if a.Currency == 1 {
		key = "usd_balance = usd_balance + %.2f"
	}
	key = fmt.Sprintf(key, a.Amount)
	err = new(Merchant).Set(key, fmt.Sprintf("id = %d", a.MchId))
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}

func (a *MerchantFund) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *MerchantFund) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *MerchantFund) Page(where string, page, limit int) ([]MerchantFund, int) {
	var datas []MerchantFund
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Merchant").Where(where).Order(fmt.Sprintf("mf.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
