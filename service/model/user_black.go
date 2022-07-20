package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type UserBlack struct {
	bun.BaseModel `bun:"table:user_black,alias:ub"`
	Id            int    `json:"id"`
	Content       string `json:"content"`
	Description   string `json:"description"`
	Type          int    `json:"type"`
	CreateTime    string `json:"create_time"`
}

func (a *UserBlack) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *UserBlack) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *UserBlack) Page(where string, page, limit int) ([]UserBlack, int) {
	var d []UserBlack
	count, _ := global.C.DB.NewSelect().Model(&d).
		Where(where).Order("ub.id desc").Offset((page - 1) * limit).Limit(limit).
		ScanAndCount(global.C.Ctx)
	return d, count
}

func (a *UserBlack) Del(where string) {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}
