package model

import (
	"github.com/uptrace/bun"
	"loante/global"
)

type ExpendStatistics struct {
	bun.BaseModel  `bun:"table:expend_statistics,alias:es"`
	Id             int             `json:"id"`
	MchId          int             `json:"mch_id" `
	CreateTime     string          `json:"create_time"`
	PutInAmount    float64         `json:"put_in_amount"`
	PutInUamount   float64         `json:"put_in_uamount"`
	PutInCount     int             `json:"put_in_count"`
	PutOutAmount   float64         `json:"put_out_amount"`
	PutOutUamount  float64         `json:"put_out_uamount"`
	PutOutCount    int             `json:"put_out_count"`
	RiskAmount     float64         `json:"risk_amount"`
	RiskUamount    float64         `json:"risk_uamount"`
	RiskCount      int             `json:"risk_count"`
	SmsAmount      float64         `json:"sms_amount" `
	SmsUamount     float64         `json:"sms_uamount" `
	SmsCount       int             `json:"sms_count" `
	ServiceAmount  float64         `json:"service_amount"`
	ServiceUamount float64         `json:"service_uamount"`
	ServiceCount   int             `json:"service_count"`
	Merchant       *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
}

type ExpendStatisticsList struct {
	Id             int             `json:"id"`
	MchId          int             `json:"mch_id" `
	CreateTime     string          `json:"create_time"`
	PutInAmount    float64         `json:"put_in_amount"`
	PutInUamount   float64         `json:"put_in_uamount"`
	PutInCount     int             `json:"put_in_count"`
	PutOutAmount   float64         `json:"put_out_amount"`
	PutOutUamount  float64         `json:"put_out_uamount"`
	PutOutCount    int             `json:"put_out_count"`
	RiskAmount     float64         `json:"risk_amount"`
	RiskUamount    float64         `json:"risk_uamount"`
	RiskCount      int             `json:"risk_count"`
	SmsAmount      float64         `json:"sms_amount" `
	SmsUamount     float64         `json:"sms_uamount" `
	SmsCount       int             `json:"sms_count" `
	ServiceAmount  float64         `json:"service_amount"`
	ServiceUamount float64         `json:"service_uamount"`
	ServiceCount   int             `json:"service_count"`
	AllAmount      float64         `json:"all_amount"`
	AllUAmount     float64         `json:"all_u_amount"`
	Merchant       *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
}

func (a *ExpendStatistics) Page(where string, page, limit int) ([]ExpendStatisticsList, int) {
	var datas []ExpendStatisticsList
	count, _ := global.C.DB.NewSelect().Model(&ExpendStatistics{}).
		Where(where).Relation("Merchant").Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx, &datas)
	if len(datas) != 0 {
		for _, data := range datas {
			data.AllAmount = data.PutInAmount + data.PutOutAmount - data.RiskAmount - data.SmsAmount - data.ServiceAmount
			data.AllUAmount = data.PutInUamount + data.PutOutUamount - data.RiskUamount - data.SmsUamount - data.ServiceUamount
		}
	}
	return datas, count
}



func (a *ExpendStatistics) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}
