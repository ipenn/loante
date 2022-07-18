package model

import (
	"loante/global"
	"fmt"
	"github.com/uptrace/bun"
)

type Config struct {
	bun.BaseModel `bun:"table:configs,alias:c"`
	Id	int	`json:"id"`
	Key		string		`json:"key"`
	Name		string		`json:"name"`
	Value		string		`json:"value"`
}

func (a *Config)Insert()  {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *Config)Update(col string,where string)  {
	global.C.DB.NewUpdate().Model(a).Column(col).Where(where).Exec(global.C.Ctx)
}

func (a *Config)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Config)Gets(where string) ([]Config, int) {
	var datas []Config
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Config)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *Config)Del(where string)  {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}
