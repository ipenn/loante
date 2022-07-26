package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"strings"
)

type StatUrge struct {
	bun.BaseModel `bun:"table:stat_urge,alias:su"`

	Id       int	`json:"id"`
	Type        int	`json:"type"`
	StatTime        string	`json:"stat_time"`
	MchId         int	`json:"mch_id"`
	UrgeCompanyId int	`json:"urge_company_id"`
	UrgeGroupId  int	`json:"urge_group_id"`
	UrgeId      int	`json:"urge_id"`
	BorrowCount       int	`json:"borrow_count"`
	BorrowNew        int	`json:"borrow_new"`
	UrgeClosedCount    int	`json:"urge_closed_count"`
	UrgeClosedRate   float64	`json:"urge_closed_rate"`
	UrgeClosedAmount   float64	`json:"urge_closed_amount"`
	VisitCount       int	`json:"visit_count"`
	VisitDetailCount int	`json:"visit_detail_count"`

	Merchant 	*MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	UrgeCompany 	*UrgeCompanyLittle `json:"urge_company" bun:"rel:belongs-to,join:urge_company_id=id"`
	UrgeGroup 	*UrgeGroupLittle `json:"urge_group" bun:"rel:belongs-to,join:urge_group_id=id"`
	Urge 		*AdminLittle 	`json:"urge" bun:"rel:belongs-to,join:urge_id=id"`
}

func (a *StatUrge)Set(col string, value interface{},  where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).SetColumn(col, "?", value).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *StatUrge) Insert() {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *StatUrge)Page(where string, page, limit int) ([]StatUrge, int) {
	var datas []StatUrge
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Merchant").Relation("UrgeCompany").Relation("UrgeGroup").Relation("Urge").Where(where).Order(fmt.Sprintf("su.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}


func (a *StatUrge)GroupPage(where, group string, page, limit int) ([]StatUrge, int) {
	var datas []StatUrge
	db := global.C.DB.NewSelect().Model(&datas).
		Column("stat_time","type","mch_id","urge_company_id","urge_group_id","urge_id").
		ColumnExpr("sum(borrow_count) as borrow_count").
		ColumnExpr("sum(borrow_new) as borrow_new").
		ColumnExpr("sum(urge_closed_count) as urge_closed_count").
		ColumnExpr("sum(urge_closed_rate) as urge_closed_rate").
		ColumnExpr("sum(urge_closed_amount) as urge_closed_amount").
		ColumnExpr("sum(visit_count) as visit_count").
		ColumnExpr("sum(visit_detail_count) as visit_detail_count").
		Relation("Merchant").Relation("UrgeCompany").Relation("UrgeGroup").Relation("Urge").Where(where)
	if len(group) > 0{
		gs := strings.Split(group, ",")
		db = db.Group(gs...)
	}
	count, _ := db.Order(fmt.Sprintf("su.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}