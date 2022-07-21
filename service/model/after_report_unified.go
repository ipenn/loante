package model

import (
	"github.com/uptrace/bun"
	"loante/global"
)

type AfterReportUnified struct {
	bun.BaseModel         `bun:"table:after_report_unified,alias:aru"`
	Id                    int            `json:"id"`
	Time                  int            `json:"time"`
	ProductId             int            `json:"product_id"`
	Type                  int            `json:"type"`
	DelayAmount           int            `json:"delay_amount"`
	DelayCount            int            `json:"delay_count"`
	DueCount              int            `json:"due_count"`
	DueAmount             int            `json:"due_amount"`
	InCollectionCount     int            `json:"in_collection_count"`
	InCollectionAmount    int            `json:"in_collection_amount"`
	NoRepayCount          int            `json:"no_repay_count"`
	NoRepayAmount         int            `json:"no_repay_amount"`
	OverDueOneDayCount    int            `json:"over_due_one_day_count"`
	OverDueOneDayAmount   int            `json:"over_due_one_day_amount"`
	OverDueThreeDayCount  int            `json:"over_due_three_day_count"`
	OverDueThreeDayAmount int            `json:"over_due_three_day_amount"`
	OverDueSevenDayCount  int            `json:"over_due_seven_day_count"`
	OverDueSevenDayAmount int            `json:"over_due_seven_day_amount"`
	Product               *ProductLittle `json:"product" bun:"rel:belongs-to,join:product_id=id"`
}

type AfterReportUnifiedList struct {
	Id                        int            `json:"id"`
	Time                      int            `json:"time"`
	ProductId                 int            `json:"product_id"`
	Type                      int            `json:"type"`
	DelayAmount               int            `json:"delay_amount"`
	DelayCount                int            `json:"delay_count"`
	DueCount                  int            `json:"due_count"`
	DueAmount                 int            `json:"due_amount"`
	InCollectionCount         int            `json:"in_collection_count"`
	InCollectionAmount        int            `json:"in_collection_amount"`
	NoRepayCount              int            `json:"no_repay_count"`
	NoRepayAmount             int            `json:"no_repay_amount"`
	OverDueOneDayCount        int            `json:"over_due_one_day_count"`
	OverDueOneDayAmount       int            `json:"over_due_one_day_amount"`
	OverDueThreeDayCount      int            `json:"over_due_three_day_count"`
	OverDueThreeDayAmount     int            `json:"over_due_three_day_amount"`
	OverDueSevenDayCount      int            `json:"over_due_seven_day_count"`
	OverDueSevenDayAmount     int            `json:"over_due_seven_day_amount"`
	DelayCountRate            float64        `json:"delay_count_rate"`
	DelayAmountRate           float64        `json:"delay_amount_rate"`
	InCollectionCountRate     float64        `json:"in_collection_count_rate"`
	InCollectionAmountRate    float64        `json:"in_collection_amount_rate"`
	NoRepayCountRate          float64        `json:"no_repay_count_rate"`
	NoRepayAmountRate         float64        `json:"no_repay_amount_rate"`
	OverDueOneDayCountRate    float64        `json:"over_due_one_day_count_rate"`
	OverDueOneDayAmountRate   float64        `json:"over_due_one_day_amount_rate"`
	OverDueThreeDayCountRate  float64        `json:"over_due_three_day_count_rate"`
	OverDueThreeDayAmountRate float64        `json:"over_due_three_day_amount_rate"`
	OverDueSevenDayCountRate  float64        `json:"over_due_seven_day_count_rate"`
	OverDueSevenDayAmountRate float64        `json:"over_due_seven_day_amount_rate"`
	Product                   *ProductLittle `json:"product" bun:"rel:belongs-to,join:product_id=id"`
}

func (a *AfterReportUnified) Page(where, group string, page, limit int) ([]AfterReportUnifiedList, int) {
	var datas []AfterReportUnifiedList
	count, _ := global.C.DB.NewSelect().Model(&AfterReportUnified{}).Relation("Product").
		ColumnExpr("SUM(delay_amount) as delay_amount").
		ColumnExpr("SUM(delay_count) as delay_count").
		ColumnExpr("SUM(due_count) as due_count").
		ColumnExpr("SUM(due_amount) as due_amount").
		ColumnExpr("SUM(in_collection_count) as in_collection_count").
		ColumnExpr("SUM(in_collection_amount) as in_collection_amount").
		ColumnExpr("SUM(no_repay_count) as no_repay_count").
		ColumnExpr("SUM(no_repay_amount) as no_repay_amount").
		ColumnExpr("SUM(over_due_one_day_count) as over_due_one_day_count").
		ColumnExpr("SUM(over_due_one_day_amount) as over_due_one_day_amount").
		ColumnExpr("SUM(over_due_three_day_count) as over_due_three_day_count").
		ColumnExpr("SUM(over_due_three_day_amount) as over_due_three_day_amount").
		ColumnExpr("SUM(over_due_seven_day_count) as over_due_seven_day_count").
		
		ColumnExpr("SUM(delay_count)*100/SUM(due_count) as delay_count_rate").
		ColumnExpr("SUM(delay_amount)*100/SUM(due_amount) as delay_amount_rate").
		ColumnExpr("SUM(in_collection_count)*100/SUM(due_count) as in_collection_count_rate").
		ColumnExpr("SUM(in_collection_amount)*100/SUM(due_amount) as in_collection_amount_rate").
		ColumnExpr("SUM(no_repay_count)*100/SUM(due_count) as no_repay_count_rate").
		ColumnExpr("SUM(no_repay_amount)*100/SUM(due_amount) as no_repay_amount_rate").
		ColumnExpr("SUM(over_due_one_day_count)*100/SUM(due_count) as over_due_one_day_count_rate").
		ColumnExpr("SUM(over_due_one_day_amount)*100/SUM(due_amount) as over_due_one_day_amount_rate").
		ColumnExpr("SUM(over_due_three_day_count)*100/SUM(due_count) as over_due_three_day_count_rate").
		ColumnExpr("SUM(over_due_three_day_amount)*100/SUM(due_amount) as over_due_three_day_amount_rate").
		ColumnExpr("SUM(over_due_seven_day_count)*100/SUM(due_count) as over_due_seven_day_count_rate").
		ColumnExpr("SUM(over_due_seven_day_amount)*100/SUM(due_amount) as over_due_seven_day_amount_rate").
		Where(where).GroupExpr(group).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx, &datas)
	return datas, count
}
