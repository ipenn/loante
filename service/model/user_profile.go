package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type UserProfile struct {
	bun.BaseModel `bun:"table:user_profiles,alias:u"`
	UserId        int    `json:"user_id"`
	HcnId         string `json:"hcn_id"`
	State         int    `json:"state"`
	Path          string `json:"path"`
	IdNo          string `json:"id_no"`
	TrueName      string `json:"true_name"`
	IdPath        string `json:"id_path"`
	AuthState     int    `json:"auth_state"`
	HcnState      int    `json:"hcn_state"`
}

func (a *UserProfile) Insert() {
	global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
}

func (a *UserProfile) Update(col string, where string) {
	sql, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	fmt.Println(sql, err)
}

func (a *UserProfile) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *UserProfile) Gets(where string) ([]UserProfile, int) {
	var datas []UserProfile
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}
