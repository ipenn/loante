package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type BorrowVisit struct {
	bun.BaseModel `bun:"table:borrow_visit,alias:bv"`
	BorrowId         int		`json:"borrow_id" bun:",pk"`
	RemindId         int		`json:"remind_id"`
	UrgeId           int		`json:"urge_id"`
	RemindCompanyId  int		`json:"remind_company_id"`
	RemindGroupId  int		`json:"remind_group_id"`
	UrgeCompanyId    int		`json:"urge_company_id"`
	UrgeGroupId    int		`json:"urge_group_id"`
	MchId            int    	`json:"mch_id"`
	UserId         int    `json:"user_id"`
	ProductId        int    	`json:"product_id"`
	RemindAssignTime string 	`json:"remind_assign_time"`
	RemindLastTime   string		`json:"remind_last_time"`
	UrgeAssignTime   string		`json:"urge_assign_time"`
	Wish      		int		`json:"wish"`
	Tag      		int		`json:"tag"`
	UrgeLastTime     string		`json:"urge_last_time"`
	Borrow  		*BorrowLittle `json:"borrow" bun:"rel:belongs-to,join:borrow_id=id"`
	RemindCompany 	*RemindCompanyLittle `json:"remind_company" bun:"rel:belongs-to,join:remind_company_id=id"`
	RemindGroup 	*RemindGroupLittle `json:"remind_group" bun:"rel:belongs-to,join:remind_group_id=id"`
	RemindUser 		*AdminLittle `json:"remind_user" bun:"rel:belongs-to,join:remind_id=id"`
	UrgeCompany 	*UrgeCompanyLittle `json:"urge_company" bun:"rel:belongs-to,join:urge_company_id=id"`
	UrgeGroup 		*UrgeGroupLittle `json:"urge_group" bun:"rel:belongs-to,join:urge_group_id=id"`
	UrgeUser 		*AdminLittle `json:"urge_user" bun:"rel:belongs-to,join:urge_id=id"`
	User 			*UserLittle `json:"user" bun:"rel:belongs-to,join:user_id=id"`
}

func (a *BorrowVisit)Insert()  {
	global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
}

func (a *BorrowVisit)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *BorrowVisit)Gets(where string) ([]BorrowVisit, int) {
	var datas []BorrowVisit
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *BorrowVisit)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *BorrowVisit)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *BorrowVisit)RemindPage(where string, page, limit int) ([]BorrowVisit, int) {
	var datas []BorrowVisit
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Borrow").Relation("User").Relation("RemindCompany").Relation("RemindGroup").Relation("RemindUser").Where(where).Order(fmt.Sprintf("bv.borrow_id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
func (a *BorrowVisit)UrgePage(where string, page, limit int) ([]BorrowVisit, int) {
	var datas []BorrowVisit
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Borrow").Relation("User").Relation("UrgeCompany").Relation("UrgeGroup").Relation("UrgeUser").Where(where).Order(fmt.Sprintf("bv.borrow_id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *BorrowVisit)GroupPage(where, group string, page, limit int) ([]BorrowVisit, int) {
	var datas []BorrowVisit
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Borrow").Relation("UrgeCompany").Relation("UrgeUser").Where(where).Order(fmt.Sprintf("bv.borrow_id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *BorrowVisit)Del(where string)  {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}