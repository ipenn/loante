package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
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

type MerchantId struct {
	MchId int `json:"id"`
}

type ExpendStatisticsValue struct {
	Id       int     `json:"id"`
	MchId    int     `json:"mch_id"`
	Count    int     `json:"count"`
	Amount   float64 `json:"amount"`
	Type     int     `json:"type"`
	Currency int     `json:"currency"`
}

// GetExpendStatisticsValue 定时任务,插入商户费用统计
func GetExpendStatisticsValue() {
	var mchid []MerchantId
	//global.C.DB.QueryContext(global.C.Ctx, `SELECT mch_id FROM merchant_fund WHERE TO_DAYS(NOW())- TO_DAYS( create_time ) = 1 GROUP BY mch_id`)
	global.C.DB.NewSelect().Column("mch_id").Model(&MerchantFund{}).Where("TO_DAYS(NOW())- TO_DAYS( create_time ) = 1").Group("mch_id").Scan(global.C.Ctx, &mchid)
	if len(mchid) != 0 {
		for _, id := range mchid {
			var esv []ExpendStatisticsValue
			var ev ExpendStatistics
			global.C.DB.NewSelect().Model(&MerchantFund{}).
				Column("id", "mch_id", "type", "currency").
				ColumnExpr("SUM( amount ) as amount").
				ColumnExpr("COUNT(id) as count").
				Where("amount != 0").
				Where(" TO_DAYS(NOW())- TO_DAYS( create_time ) = 1").
				Where(fmt.Sprintf("mch_id=%d", id.MchId)).
				GroupExpr("type").GroupExpr("currency").GroupExpr("amount").
				Scan(global.C.Ctx, &esv)
			for _, value := range esv {
				ev.MchId = value.MchId
				ev.CreateTime = tools.GetFormatDay()
				if value.Type == 0 && value.Currency == 1 {
					ev.PutInAmount = value.Amount
				}
				if value.Type == 0 && value.Currency == 2 {
					ev.PutInUamount = value.Amount
				}
				if value.Type == 0 {
					ev.PutInCount += value.Count
				}
				if value.Type == 1 && value.Currency == 1 {
					ev.PutOutAmount = value.Amount
				}
				if value.Type == 1 && value.Currency == 2 {
					ev.PutOutUamount = value.Amount
				}
				if value.Type == 1 {
					ev.PutOutCount += value.Count
				}
				if value.Type == 5 && value.Currency == 1 {
					ev.RiskAmount = value.Amount
				}
				if value.Type == 5 && value.Currency == 2 {
					ev.RiskUamount = value.Amount
				}
				if value.Type == 5 {
					ev.RiskCount += value.Count
				}
				if value.Type == 4 && value.Currency == 1 {
					ev.SmsAmount = value.Amount
				}
				if value.Type == 4 && value.Currency == 2 {
					ev.SmsUamount = value.Amount
				}
				if value.Type == 4 {
					ev.SmsCount += value.Count
				}
				if value.Type == 3 && value.Currency == 1 {
					ev.ServiceAmount = value.Amount
				}
				if value.Type == 3 && value.Currency == 2 {
					ev.ServiceUamount = value.Amount
				}
				if value.Type == 3 {
					ev.ServiceCount += value.Count
				}
			}
			ev.Insert()
		}
	}
}

func (a *ExpendStatistics) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}
