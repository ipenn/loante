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
	AvgLoan          float64         `json:"avg_loan"`
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
				data.AvgLoan = tools.ToFloat64(data.LoanCount )/ tools.ToFloat64(data.ApplyCount)
			}
		}
	}
	return datas, count
}
