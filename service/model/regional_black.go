package model

import (
	"github.com/uptrace/bun"
	"loante/global"
)

type RegionalBlack struct {
	bun.BaseModel  `bun:"table:regional_black,alias:rb"`
	Id             int    `json:"id"`
	RegionName     string `json:"region_name"`
	IsBlackRegion  int    `json:"is_black_region"`
	ParentRegionId int    `json:"parent_region_id"`
}

func (a *RegionalBlack) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *RegionalBlack) Page(where string, page, limit int) ([]RegionalBlack, int) {
	var d []RegionalBlack
	count, _ := global.C.DB.NewSelect().Model(&d).
		Where(where).Order("rb.id desc").Offset((page - 1) * limit).Limit(limit).
		ScanAndCount(global.C.Ctx)
	return d, count
}

func (a *RegionalBlack) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}
