package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type BorrowVisitDetail struct {
	bun.BaseModel `bun:"table:borrow_visit_detail,alias:bvd"`
	Id         	int		`json:"id" bun:",pk"`
	Type         	int		`json:"type"`
	MchId     	int		`json:"mch_id"`
	ProductId 	int		`json:"product_id"`
	BorrowId 	int		`json:"borrow_id"`
	UserId          int		`json:"user_id"`
	UrgeId        int		`json:"urge_id"`
	UrgeCompanyId int		`json:"urge_company_id"`
	UrgeGroupId   int		`json:"urge_group_id"`
	ContactName  string		`json:"contact_name"`
	ContactPhone string		`json:"contact_phone"`
	Relationship int		`json:"relationship"`
	CreateTime   string		`json:"create_time"`
	Tag          int		`json:"tag"`
	Wish         int		`json:"wish"`
	Remark                  string		`json:"remark"`
	PromisedRepaymentAmount string		`json:"promised_repayment_amount"`
	PromisedRepaymentTime string		`json:"promised_repayment_time"`
	NextVisitTime         string		`json:"next_visit_time"`

	Product  		*ProductLittle 	`json:"product" bun:"rel:belongs-to,join:product_id=id"`
	Merchant  		*MerchantLittle 	`json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	User  			*UserLittle 	`json:"user" bun:"rel:belongs-to,join:user_id=id"`
	Borrow  		*BorrowMini 	`json:"borrow" bun:"rel:belongs-to,join:borrow_id=id"`

	UrgeCompany 	*UrgeCompanyLittle 	`json:"urge_company" bun:"rel:belongs-to,join:urge_company_id=id"`
	UrgeGroup 		*UrgeGroupLittle 	`json:"urge_group" bun:"rel:belongs-to,join:urge_group_id=id"`
	UrgeUser 		*AdminLittle 	`json:"urge_user" bun:"rel:belongs-to,join:urge_id=id"`

	RemindCompany 	*RemindCompanyLittle 	`json:"remind_company" bun:"rel:belongs-to,join:urge_company_id=id"`
	RemindGroup 	*RemindGroupLittle 		`json:"remind_group" bun:"rel:belongs-to,join:urge_group_id=id"`
	RemindUser 		*AdminLittle 			`json:"remind_user" bun:"rel:belongs-to,join:urge_id=id"`
}

func (a *BorrowVisitDetail)Insert()  {
	global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
}

func (a *BorrowVisitDetail)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *BorrowVisitDetail)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *BorrowVisitDetail)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *BorrowVisitDetail)UrgePage(where string, page, limit int) ([]BorrowVisitDetail, int) {
	var datas []BorrowVisitDetail
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Merchant").Relation("Product").Relation("User").Relation("Borrow").Relation("UrgeCompany").Relation("UrgeGroup").Relation("UrgeUser").
		Where(where).Order(fmt.Sprintf("bvd.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
func (a *BorrowVisitDetail)RemindPage(where string, page, limit int) ([]BorrowVisitDetail, int) {
	var datas []BorrowVisitDetail
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Merchant").Relation("Product").Relation("User").Relation("Borrow").Relation("RemindCompany").Relation("RemindGroup").Relation("RemindUser").
		Where(where).Order(fmt.Sprintf("bvd.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}