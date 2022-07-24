package cron

import (
	"fmt"
	"loante/global"
	"loante/service/model"
	"loante/tools"
)

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
	global.C.DB.NewSelect().Column("product_id").Model(&model.Borrow{}).Where("TO_DAYS(NOW())- TO_DAYS( create_time ) = 1").Group("product_id").Scan(global.C.Ctx, &product)
	if len(product) != 0 {
		for _, id := range product {
			var st model.StatTraffic

			rows2, _ := global.C.DB.QueryContext(global.C.Ctx, fmt.Sprintf("SELECT\n\tCOUNT( id ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t1 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND loan_type = 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t2 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND `status` >= 4 \n\tAND loan_type = 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\tIFNULL( SUM( loan_amount ), 0 ) AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t3 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND `status` >= 5 \n\tAND loan_type = 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tIFNULL( COUNT( id ), 0 ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t4 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND loan_type = 0 \n\tAND product_id = %d \nGROUP BY\n\tuid UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t11 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND loan_type > 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t12 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND `status` >= 4 \n\tAND loan_type > 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\tIFNULL( SUM( loan_amount ), 0 ) AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t13 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND STATUS >= 5 \n\tAND loan_type > 0 \n\tAND product_id = %d UNION ALL\nSELECT\n\tIFNULL( COUNT( id ), 0 ) AS count,\n\t0 AS amount,\n\tIFNULL( product_id, 0 ) product_id,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t14 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND loan_type > 0 \n\tAND product_id = %d \nGROUP BY\n\tuid", id.ProductId, id.ProductId, id.ProductId, id.ProductId, id.ProductId, id.ProductId, id.ProductId, id.ProductId))
			for rows2.Next() {
				value := new(statTrafficValue)
				err := global.C.DB.ScanRow(global.C.Ctx, rows2, value)
				if err == nil {
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
			//fmt.Println(st)
		}
	}
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
	global.C.DB.NewSelect().Column("mch_id").Model(&model.MerchantFund{}).Where("TO_DAYS(NOW())- TO_DAYS( create_time ) = 1").Group("mch_id").Scan(global.C.Ctx, &mchid)
	if len(mchid) != 0 {
		for _, id := range mchid {
			var esv []ExpendStatisticsValue
			var ev model.ExpendStatistics
			global.C.DB.NewSelect().Model(&model.MerchantFund{}).
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


type GetMerchantStatistics struct {
	Amount float64 `bun:"amount" json:"amount"`
	Count  int     `bun:"count" json:"count"`
	Type   int     `bun:"type" json:"type"`
	MchId  int     `bun:"mch_id" json:"mch_id"`
}

// GetMerchantStatisticsValue 定时任务,插入商户日报
func GetMerchantStatisticsValue() {
	var mchid []MerchantId
	//rows, _ := global.C.DB.QueryContext(global.C.Ctx, `SELECT mch_id FROM merchant_fund WHERE TO_DAYS(NOW())- TO_DAYS( create_time ) = 1 GROUP BY mch_id`)
	global.C.DB.NewSelect().Column("mch_id").Model(&model.MerchantFund{}).Where("TO_DAYS(NOW())- TO_DAYS( create_time ) = 1").Group("mch_id").Scan(global.C.Ctx, &mchid)
	if len(mchid) != 0 {
		for _, id := range mchid {
			context, _ := global.C.DB.QueryContext(global.C.Ctx, "SELECT\n\tSUM( amount ) AS amount,\n\tCOUNT( id ) AS count,\n\tmch_id,\n\t1 AS type \nFROM\n\tmerchant_fund \nWHERE\n\tmch_id = ? \n\tAND TO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND type = 1 \n\tAND currency = 1 UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\tIFNULL( SUM( loan_amount ), 0 ) AS amount,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t3 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 AND mch_id = ? UNION ALL\nSELECT\n\tCOUNT( id ) AS count,\n\tIFNULL( SUM( loan_amount ), 0 ) AS amount,\n\tIFNULL( mch_id, 0 ) mch_id,\n\t13 AS type \nFROM\n\tborrow \nWHERE\n\tTO_DAYS(\n\tNOW())- TO_DAYS( create_time ) = 1 \n\tAND STATUS >= 5 AND mch_id =?", id.MchId, id.MchId, id.MchId)
			defer context.Close()
			var ms model.MerchantStatistics
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
