package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type UrgeGroupLittle struct {
	bun.BaseModel `bun:"table:urge_group,alias:ug"`
	Id            int          `json:"id" bun:",pk"`
	GroupName     string       `json:"group_name"`
}
type UrgeGroup struct {
	bun.BaseModel `bun:"table:urge_group,alias:ug"`
	UrgeGroupLittle
	CompanyId     int          `json:"company_id"`
	MchId         int          `json:"mch_id"`
	AdminId       int          `json:"admin_id"`
	Status        int          `json:"status"`
	UrgeCompany   *UrgeCompany `json:"urge_company" bun:"rel:belongs-to,join:company_id=id"`
	Merchant      *Merchant    `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	Admin         *Admin       `json:"admin" bun:"rel:belongs-to,join:admin_id=id"`
}

func (a *UrgeGroup) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *UrgeGroup) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *UrgeGroup) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *UrgeGroup) Page(where string, page, limit int) ([]UrgeGroup, int) {
	var d []UrgeGroup
	count, _ := global.C.DB.NewSelect().Model(&d).Relation("UrgeCompany").
		Relation("Merchant").Relation("Admin").
		Where(where).Order("ug.id desc").Offset((page - 1) * limit).Limit(limit).
		ScanAndCount(global.C.Ctx)
	return d, count
}
