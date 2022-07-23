package model

import (
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type MerchantStatistics struct {
	bun.BaseModel `bun:"table:merchant_statistics,alias:ms"`
	Id            int             `json:"id"`
	MchId         int             `json:"mch_id"`
	ServiceAmount float64         `json:"service_amount"`
	ServiceCount  int             `json:"service_count"`
	LoanAmount    float64         `json:"loan_amount"`
	LoanCount     int             `json:"loan_count"`
	ApplyCount    int             `json:"apply_count"`
	CreateTime    string          `json:"create_time"`
	Merchant      *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
}
type MerchantStatisticsList struct {
	Id               int             `json:"id"`
	MchId            int             `json:"mch_id"`
	ServiceAmount    float64         `json:"service_amount"`
	ServiceCount     int             `json:"service_count"`
	LoanAmount       float64         `json:"loan_amount"`
	LoanCount        int             `json:"loan_count"`
	ApplyCount       int             `json:"apply_count"`
	CreateTime       string          `json:"create_time"`
	AvgServiceAmount float64         `json:"avg_service_amount"`
	AvgLoan          int             `json:"avg_loan"`
	Merchant         *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
}

func (a *MerchantStatistics) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *MerchantStatistics) Page(where string, page, limit int) ([]MerchantStatisticsList, int) {
	var datas []MerchantStatisticsList
	count, _ := global.C.DB.NewSelect().Model(&MerchantStatistics{}).
		Where(where).Relation("Merchant").Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx, &datas)
	if len(datas) != 0 {
		for _, data := range datas {
			data.AvgServiceAmount = data.ServiceAmount / float64(data.ServiceCount)
			if data.LoanCount == 0 {
				data.AvgLoan = 0
			} else {
				data.AvgLoan = data.LoanCount * 100 / data.ApplyCount
			}
		}
	}
	return datas, count
}

type GetMerchantStatistics struct {
	Amount float64 `bun:"amount" json:"amount"`
	Count  int     `bun:"count" json:"count"`
	Type   int     `bun:"type" json:"type"`
	MchId  int     `bun:"mch_id" json:"mch_id"`
}

// GetMerchantStatisticsValue 定时任务,插入商户费用统计
func GetMerchantStatisticsValue() {
	var mchid []MerchantId
	//rows, _ := global.C.DB.QueryContext(global.C.Ctx, `SELECT mch_id FROM merchant_fund WHERE TO_DAYS(NOW())- TO_DAYS( create_time ) = 1 GROUP BY mch_id`)
	global.C.DB.NewSelect().Column("mch_id").Model(&MerchantFund{}).Where("TO_DAYS(NOW())- TO_DAYS( create_time ) = 1").Group("mch_id").Scan(global.C.Ctx, &mchid)
	if len(mchid) != 0 {
		for _, id := range mchid {
			context, _ := global.C.DB.QueryContext(global.C.Ctx, "SELECT\n\tSUM( amount ) AS amount,\n\tCOUNT( id ) AS count,\n\tmch_id,\n\t1 AS type \nFROM\n\tmerchant_fund \nWHERE\n\tmch_id = ? \n\tAND TO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND type = 1 \n\tAND currency = 1 UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\tIFNULL( SUM( loan_amount ), 0 ) AS amount,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t3 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 AND mch_id = ? UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\tIFNULL( SUM( loan_amount ), 0 ) AS amount,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t13 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND STATUS >= 5 AND mch_id =?", id.MchId, id.MchId, id.MchId)
			defer context.Close()
			var ms MerchantStatistics
			for context.Next() {
				value := new(GetMerchantStatistics)
				err := global.C.DB.ScanRow(global.C.Ctx, context, value)
				if err == nil {
					ms.MchId = value.MchId
					ms.CreateTime = tools.GetFormatDay()
					if value.Type == 1 {
						ms.ServiceAmount = value.Amount
						ms.ServiceCount = value.Count
					}
					if value.Type == 3 {
						ms.ApplyCount = value.Count
					}
					if value.Type == 13 {
						ms.LoanCount = value.Count
						ms.LoanAmount = value.Amount
					}
				}
			}
			ms.Insert()
		}
	}
}
