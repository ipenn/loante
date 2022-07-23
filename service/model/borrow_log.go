package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type BorrowLog struct {
	bun.BaseModel `bun:"table:borrow_log,alias:bl"`
	Id          string	`json:"id"`
	BorrowId   string	`json:"borrow_id"`
	CreateTime string	`json:"create_time"`
	Amount         string	`json:"amount"`
	BeRepaidAmount string	`json:"be_repaid_amount"`
	Remark         string	`json:"remark"`
}

func (a *BorrowLog) Insert() {
	a.CreateTime = tools.GetFormatTime()
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *BorrowLog) Page(where string, page, limit int) ([]BorrowLog, int) {
	var datas []BorrowLog
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
