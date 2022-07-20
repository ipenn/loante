package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type SystemSetting struct {
	bun.BaseModel `bun:"table:system_setting,alias:ss"`
	Id            int    `json:"id"`
	UpdateTime    string `json:"update_time"`
	ParamKey      string `json:"param_key"`
	ParamValue    string `json:"param_value"`
	ParamType     string `json:"param_type"`
	Remark        string `json:"remark"`
}

func (a *SystemSetting) Page(page, limit int) ([]SystemSetting, int) {
	var d []SystemSetting
	count, _ := global.C.DB.NewSelect().Model(&d).Order("ss.id desc").Offset((page - 1) * limit).Limit(limit).
		ScanAndCount(global.C.Ctx)
	return d, count
}

func (a *SystemSetting) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *SystemSetting) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}
