package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type RemindGroup struct {
	bun.BaseModel `bun:"table:remind_group,alias:rg"`
	Id            int             `json:"id" bun:",pk"`
	CompanyId     int             `json:"company_id"`
	MchId         int             `json:"mch_id"`
	AdminId       int             `json:"admin_id"`
	GroupName     string          `json:"group_name"`
	Status        int             `json:"status"`
	RemindCompany *RemindCompany  `json:"remind_company" bun:"rel:belongs-to,join:company_id=id"`
	Merchant      *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	Admin         *Admin          `json:"admin" bun:"rel:belongs-to,join:admin_id=id"`
}

func (a *RemindGroup) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *RemindGroup) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *RemindGroup) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *RemindGroup) Page(where string, page, limit int) ([]RemindGroup, int) {
	var d []RemindGroup
	count, _ := global.C.DB.NewSelect().Model(&d).Relation("RemindCompany").Relation("Merchant").
		Relation("Admin").Where(where).Order("rg.id desc").Offset((page - 1) * limit).Limit(limit).
		ScanAndCount(global.C.Ctx)
	return d, count
}
