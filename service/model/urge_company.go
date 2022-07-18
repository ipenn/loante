package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type UrgeCompany struct {
	bun.BaseModel `bun:"table:urge_company,alias:uc"`
	Id            int    `json:"id" bun:",pk"`
	AdminId       int    `json:"admin_id"`
	UserName      string `json:"user_name"`
	MchId         int    `json:"mch_id"`
	CreateTime    string `json:"create_time"`
	CompanyName   string `json:"company_name"`
	Description   string `json:"description"`
	Merchant      *Merchant      `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
}

func (a *UrgeCompany) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *UrgeCompany) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *UrgeCompany) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *UrgeCompany) Page(where string, page, limit int) ([]UrgeCompany, int) {
	var d []UrgeCompany
	count, _ := global.C.DB.NewSelect().Model(&UrgeCompany{}).Relation("Merchant").
		Where(where).Order("uc.id desc").Offset((page-1)*limit).Limit(limit).
		ScanAndCount(global.C.Ctx, &d)
	return d, count
}
