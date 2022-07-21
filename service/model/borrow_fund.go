package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type BorrowFund struct {
	bun.BaseModel   `bun:"table:borrow_fund,alias:bf"`
	Id              int		`json:"id" bun:",pk"`
	BorrowId        int		`json:"borrow_id"`
	UserId          int		`json:"user_id"`
	PaymentId          int		`json:"payment_id"`
	BeRepaidAmount  float64		`json:"be_repaid_amount"`
	Amount          float64		`json:"amount"`
	PayAmounted      float64		`json:"pay_amounted"`
	RemainingAmount float64		`json:"remaining_amount"`
	CreateTime      string		`json:"create_time"`
	Remark          string		`json:"remark"`
	Type    		int           `json:"type"`
	OrderNo 		string        `json:"order_no"`
	Borrow  		*BorrowLittle `json:"borrow" bun:"rel:belongs-to,join:borrow_id=id"`
}

func (a *BorrowFund)Insert()  {
	a.CreateTime = tools.GetFormatTime()
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *BorrowFund)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *BorrowFund)Gets(where string) ([]BorrowFund, int) {
	var datas []BorrowFund
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *BorrowFund)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *BorrowFund)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *BorrowFund)Page(where string, page, limit int) ([]BorrowFund, int) {
	var datas []BorrowFund
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Borrow").Where(where).Order(fmt.Sprintf("b.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *BorrowFund)Del(where string)  {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}