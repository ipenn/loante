package model

import (
	"github.com/uptrace/bun"
	"loante/global"
)

type StatTraffic struct {
	bun.BaseModel    `bun:"table:stat_traffic,alias:st"`
	Id               int                   `json:"id"`
	Time             string                `json:"time"`
	MchId            int                   `json:"mch_id"`
	ProductId        int                   `json:"product_id"`
	ReferrerConfigId int                   `json:"referrer_config_id"`
	Apply            int                   `json:"apply"`
	ApplyPass        int                   `json:"apply_pass"`
	LoanPass         int                   `json:"loan_pass"`
	LoanAmount       int                   `json:"loan_amount"`
	NewApply         int                   `json:"new_apply"`
	NewApplyPass     int                   `json:"new_apply_pass"`
	NewLoanPass      int                   `json:"new_loan_pass"`
	NewLoanAmount    int                   `json:"new_loan_amount"`
	OldApply         int                   `json:"old_apply"`
	OldApplyPass     int                   `json:"old_apply_pass"`
	OldLoanPass      int                   `json:"old_loan_pass"`
	OldLoanAmount    int                   `json:"old_loan_amount"`
	ApplyUser        int                   `json:"apply_user"`
	NewApplyUser     int                   `json:"new_apply_user"`
	OldApplyUser     int                   `json:"old_apply_user"`
	Product          *ProductLittle        `json:"product" bun:"rel:belongs-to,join:product_id=id"`
	Merchant         *MerchantLittle       `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	ReferrerConfig   *ReferrerConfigLittle `json:"referrer_config" bun:"rel:belongs-to,join:referrer_config_id=id"`
}

type StatTrafficList struct {
	Id               int                   `json:"id"`
	Time             string                `json:"time"`
	MchId            int                   `json:"mch_id"`
	ProductId        int                   `json:"product_id"`
	ReferrerConfigId int                   `json:"referrer_config_id"`
	Apply            int                   `json:"apply"`
	ApplyPass        int                   `json:"apply_pass"`
	LoanPass         int                   `json:"loan_pass"`
	LoanAmount       int                   `json:"loan_amount"`
	NewApply         int                   `json:"new_apply"`
	NewApplyPass     int                   `json:"new_apply_pass"`
	NewLoanPass      int                   `json:"new_loan_pass"`
	NewLoanAmount    int                   `json:"new_loan_amount"`
	OldApply         int                   `json:"old_apply"`
	OldApplyPass     int                   `json:"old_apply_pass"`
	OldLoanPass      int                   `json:"old_loan_pass"`
	OldLoanAmount    int                   `json:"old_loan_amount"`
	ApplyPassRate    float64               `json:"apply_pass_rate"`
	LendingRate      float64               `json:"lending_rate"`
	LoanPassRate     float64               `json:"loan_pass_rate"`
	NewApplyPassRate float64               `json:"new_apply_pass_rate"`
	NewLendingRate   float64               `json:"new_lending_rate"`
	NewLoanPassRate  float64               `json:"new_loan_pass_rate"`
	OldApplyPassRate float64               `json:"old_apply_pass_rate"`
	OldLendingRate   float64               `json:"old_lending_rate"`
	OldLoanPassRate  float64               `json:"old_loan_pass_rate"`
	ApplyUser        int                   `json:"apply_user"`
	NewApplyUser     int                   `json:"new_apply_user"`
	OldApplyUser     int                   `json:"old_apply_user"`
	AvgApply         float64               `json:"avg_apply"`
	NewAvgApply      float64               `json:"new_avg_apply"`
	OldAvgApply      float64               `json:"old_avg_apply"`
	Product          *ProductLittle        `json:"product" bun:"rel:belongs-to,join:product_id=id"`
	Merchant         *MerchantLittle       `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	ReferrerConfig   *ReferrerConfigLittle `json:"referrer_config" bun:"rel:belongs-to,join:referrer_config_id=id"`
}

func (a *StatTraffic) Page(where, group string, page, limit int) ([]StatTrafficList, int) {
	var datas []StatTrafficList
	count := 0
	if group == "" {
		count, _ = global.C.DB.NewSelect().Model(&StatTraffic{}).
			Column("id", "mch_id", "product_id", "time").
			ColumnExpr("SUM( apply ) AS apply").
			ColumnExpr("SUM( apply_pass ) AS apply_pass").
			ColumnExpr("SUM( loan_pass ) AS loan_pass").
			ColumnExpr("SUM( loan_amount ) AS loan_amount").
			ColumnExpr("SUM( new_apply ) AS new_apply").
			ColumnExpr("SUM( new_apply_pass ) AS new_apply_pass").
			ColumnExpr("SUM( new_loan_pass ) AS new_loan_pass").
			ColumnExpr("SUM( new_loan_amount ) AS new_loan_amount").
			ColumnExpr("SUM( old_apply ) AS old_apply").
			ColumnExpr("SUM( old_apply_pass ) AS old_apply_pass").
			ColumnExpr("SUM( old_loan_pass ) AS old_loan_pass").
			ColumnExpr("SUM( old_loan_amount ) AS old_loan_amount").
			ColumnExpr("SUM(apply_pass)/SUM(apply) as apply_pass_rate").
			ColumnExpr("SUM(loan_pass)/SUM(apply) as lending_rate").
			ColumnExpr("SUM(loan_pass)/SUM(apply_pass) as loan_pass_rate").
			ColumnExpr("SUM(new_apply_pass)/SUM(new_apply) as new_apply_pass_rate").
			ColumnExpr("SUM(new_loan_pass)/SUM(new_apply) as new_lending_rate").
			ColumnExpr("SUM(new_loan_pass)/SUM(new_apply_pass) as new_loan_pass_rate").
			ColumnExpr("SUM(old_apply_pass)/SUM(old_apply) as old_apply_pass_rate").
			ColumnExpr("SUM(old_loan_pass)/SUM(old_apply) as old_lending_rate").
			ColumnExpr("SUM(old_loan_pass)/SUM(old_apply_pass) as old_loan_pass_rate").
			ColumnExpr("SUM(apply)/SUM(apply_user) as avg_apply").
			ColumnExpr("SUM(new_apply)/SUM(new_apply_user) as new_avg_apply").
			ColumnExpr("SUM(old_apply) /SUM(old_apply_user)as old_avg_apply").
			Where(where).Offset((page-1)*limit).
			Relation("Product").Relation("Merchant").Relation("ReferrerConfig").Limit(limit).ScanAndCount(global.C.Ctx, &datas)
	} else {
		count, _ = global.C.DB.NewSelect().Model(&StatTraffic{}).
			Column("id", "mch_id", "product_id", "time").
			ColumnExpr("SUM( apply ) AS apply").
			ColumnExpr("SUM( apply_pass ) AS apply_pass").
			ColumnExpr("SUM( loan_pass ) AS loan_pass").
			ColumnExpr("SUM( loan_amount ) AS loan_amount").
			ColumnExpr("SUM( new_apply ) AS new_apply").
			ColumnExpr("SUM( new_apply_pass ) AS new_apply_pass").
			ColumnExpr("SUM( new_loan_pass ) AS new_loan_pass").
			ColumnExpr("SUM( new_loan_amount ) AS new_loan_amount").
			ColumnExpr("SUM( old_apply ) AS old_apply").
			ColumnExpr("SUM( old_apply_pass ) AS old_apply_pass").
			ColumnExpr("SUM( old_loan_pass ) AS old_loan_pass").
			ColumnExpr("SUM( old_loan_amount ) AS old_loan_amount").
			ColumnExpr("SUM(apply_pass)/SUM(apply) as apply_pass_rate").
			ColumnExpr("SUM(loan_pass)/SUM(apply) as lending_rate").
			ColumnExpr("SUM(loan_pass)/SUM(apply_pass) as loan_pass_rate").
			ColumnExpr("SUM(new_apply_pass)/SUM(new_apply) as new_apply_pass_rate").
			ColumnExpr("SUM(new_loan_pass)/SUM(new_apply) as new_lending_rate").
			ColumnExpr("SUM(new_loan_pass)/SUM(new_apply_pass) as new_loan_pass_rate").
			ColumnExpr("SUM(old_apply_pass)/SUM(old_apply) as old_apply_pass_rate").
			ColumnExpr("SUM(old_loan_pass)/SUM(old_apply) as old_lending_rate").
			ColumnExpr("SUM(old_loan_pass)/SUM(old_apply_pass) as old_loan_pass_rate").
			ColumnExpr("SUM(apply)/SUM(apply_user) as avg_apply").
			ColumnExpr("SUM(new_apply)/SUM(new_apply_user) as new_avg_apply").
			ColumnExpr("SUM(old_apply) /SUM(old_apply_user)as old_avg_apply").
			Where(where).GroupExpr(group).Offset((page-1)*limit).
			Relation("Product").Relation("Merchant").Relation("ReferrerConfig").Limit(limit).ScanAndCount(global.C.Ctx, &datas)
	}
	return datas, count
}


func (a *StatTraffic) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}
