package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type ReferrerConfig struct {
	bun.BaseModel `bun:"table:referrer_config,alias:rc"`
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Keyworks      string `json:"keyworks"`
	AppToken      string `json:"app_token"`
	Remark        string `json:"remark"`
	Status        int    `json:"status"`
	IsRejectApply int    `json:"is_reject_apply"`
	IsNeedReview  int    `json:"is_need_review"`
}

func (a *ReferrerConfig) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ReferrerConfig) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ReferrerConfig) Set(col string, value interface{}, where string) error {
	_, err := global.C.DB.NewUpdate().Model(a).SetColumn(col, "?", value).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}

func (a *ReferrerConfig) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ReferrerConfig) Page(where string, page, limit int) ([]ReferrerConfig, int) {
	var datas []ReferrerConfig
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
