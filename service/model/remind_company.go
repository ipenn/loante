package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type RemindCompany struct {
	bun.BaseModel `bun:"table:remind_company,alias:rc"`
	Id            int    `json:"id" bun:",pk"`
	AdminId       int    `json:"admin_id"`
	UserName      string `json:"user_name"`
	MchId         int    `json:"mch_id"`
	CreateTime    string `json:"create_time"`
	CompanyName   string `json:"company_name"`
	Description   string `json:"description"`
	Merchant      *Merchant      `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
}

func (a *RemindCompany) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *RemindCompany) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *RemindCompany) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *RemindCompany) Page(where string, page, limit int) ([]RemindCompany, int) {
	var d []RemindCompany
	count, _ := global.C.DB.NewSelect().Model(&RemindCompany{}).Relation("Merchant").
		Where(where).Order("rc.id desc").Offset((page-1)*limit).Limit(limit).
		ScanAndCount(global.C.Ctx, &d)
	return d, count
}
