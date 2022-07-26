package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type BorrowUtr struct {
	bun.BaseModel `bun:"table:borrow_utr,alias:bu"`
	Id       int	`json:"id"`
	BorrowId   int	`json:"borrow_id"`
	MchId     int	`json:"mch_id"`
	Type     int	`json:"type"` //来源 0 用户提交 1 催收员提交
	ProductId int	`json:"product_id"`
	Status   int	`json:"status"`
	UtrCode string	`json:"utr_code"`
	UtrPath string	`json:"utr_path"`
	UserId     int	`json:"user_id"`
	UrgeId     int	`json:"urge_id"`
	CreateTime string	`json:"create_time"`
	RejectReason string	`json:"reject_reason"`
	RejectRemark string	`json:"reject_remark"`
	Remark        string	`json:"remark"`
	Merchant 	*MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	Product 	*ProductLittle `json:"product" bun:"rel:belongs-to,join:product_id=id"`
	Borrow 		*BorrowMini 	`json:"borrow" bun:"rel:belongs-to,join:borrow_id=id"`
	User 		*UserLittle 	`json:"user" bun:"rel:belongs-to,join:user_id=id"`
	Urge 		*AdminLittle 	`json:"admin" bun:"rel:belongs-to,join:urge_id=id"`
}

func (a *BorrowUtr)Insert()  {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *BorrowUtr)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *BorrowUtr)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *BorrowUtr)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *BorrowUtr)Page(where string, page, limit int) ([]BorrowUtr, int) {
	var datas []BorrowUtr
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Merchant").Relation("Borrow").Relation("User").Relation("Urge").Where(where).Order(fmt.Sprintf("bu.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *BorrowUtr)Del(where string)  {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}