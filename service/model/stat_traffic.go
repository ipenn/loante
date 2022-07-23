package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
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
	count, _ := global.C.DB.NewSelect().Model(&StatTraffic{}).
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
		ColumnExpr("SUM(apply_pass)*100/SUM(apply) as apply_pass_rate").
		ColumnExpr("SUM(loan_pass)*100/SUM(apply) as lending_rate").
		ColumnExpr("SUM(loan_pass)*100/SUM(apply_pass) as loan_pass_rate").
		ColumnExpr("SUM(apply)/SUM(apply_user) as avg_apply").
		ColumnExpr("SUM(new_apply)/SUM(new_apply_user) as new_avg_apply").
		ColumnExpr("SUM(old_apply) /SUM(old_apply_user)as old_avg_apply").
		Where(where).GroupExpr(group).Offset((page-1)*limit).
		Relation("Product").Relation("Merchant").Relation("ReferrerConfig").Limit(limit).ScanAndCount(global.C.Ctx, &datas)
	return datas, count
}

type statTrafficValue struct {
	Count     int `json:"count"`
	Amount    int `json:"amount"`
	ProductId int `json:"product_id"`
	MchId     int `json:"mch_id"`
	Type      int `json:"type"`
}

type productId struct {
	ProductId int `json:"product_id"`
}

// GetStatTrafficValue 定时任务,插入报表数据
func GetStatTrafficValue() {
	var product []productId
	global.C.DB.NewSelect().Column("product_id").Model(&Borrow{}).Where("TO_DAYS(NOW())- TO_DAYS( create_time ) = 1").Group("product_id").Scan(global.C.Ctx, &product)
	if len(product) != 0 {
		for _, id := range product {
			var st StatTraffic

			rows2, _ := global.C.DB.QueryContext(global.C.Ctx, fmt.Sprintf("SELECT\n\tCOUNT( id ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t1 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND loan_type = 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t2 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND `status` >= 4 \n\tAND loan_type = 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\tIFNULL( SUM( loan_amount ), 0 ) AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t3 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND `status` >= 5 \n\tAND loan_type = 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tIFNULL( COUNT( id ), 0 ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t4 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND loan_type = 0 \n\tAND product_id = %d \nGROUP BY\n\tuid UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t11 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND loan_type > 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t12 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND `status` >= 4 \n\tAND loan_type > 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\tIFNULL( SUM( loan_amount ), 0 ) AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t13 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND STATUS >= 5 \n\tAND loan_type > 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tIFNULL( COUNT( id ), 0 ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t14 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND loan_type > 0 \n\tAND product_id = %d \nGROUP BY\n\tuid", id, id, id, id, id, id, id, id))
			for rows2.Next() {
				value := new(statTrafficValue)
				err := global.C.DB.ScanRow(global.C.Ctx, rows2, value)
				if err==nil {
					if value.Type == 1 {
						st.NewApply = value.Count
						st.ProductId = value.ProductId
						st.MchId = value.MchId
					}
					if value.Type == 2 {
						st.NewApplyPass = value.Count
					}
					if value.Type == 3 {
						st.NewLoanPass = value.Count
						st.NewLoanAmount = value.Amount
					}
					if value.Type == 4 {
						st.NewApplyUser = value.Count
					}
					if value.Type == 11 {
						st.OldApply = value.Count
					}
					if value.Type == 12 {
						st.OldApplyPass = value.Count
					}
					if value.Type == 13 {
						st.OldLoanPass = value.Count
						st.OldLoanAmount = value.Amount
					}
					if value.Type == 14 {
						st.OldApplyUser = value.Count
					}
					st.Apply = st.NewApply + st.OldApply
					st.ApplyUser = st.NewApplyUser + st.OldApplyUser
					st.ApplyPass = st.NewApplyPass + st.OldApplyPass
					st.LoanPass = st.NewLoanPass + st.OldLoanPass
					st.LoanAmount = st.NewLoanAmount + st.OldLoanAmount
					st.Time = tools.GetFormatDay()
				}
			}
			st.Insert()
		}
	}
}

func (a *StatTraffic) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}
