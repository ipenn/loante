package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type WhitePhone struct {
	bun.BaseModel `bun:"table:white_phone,alias:wb"`
	Id            int    `json:"id"`
	Phone         string `json:"phone"`
	Description   string `json:"description"`
	CreateTime    string `json:"create_time"`
}

func (a *WhitePhone) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *WhitePhone) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *WhitePhone) Page(where string, page, limit int) ([]WhitePhone, int) {
	var d []WhitePhone
	count, _ := global.C.DB.NewSelect().Model(&d).
		Where(where).Order("wb.id desc").Offset((page - 1) * limit).Limit(limit).
		ScanAndCount(global.C.Ctx)
	return d, count
}

func (a *WhitePhone) Del(where string) {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}
