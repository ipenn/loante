package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type Score struct {
	bun.BaseModel `bun:"table:score,alias:s"`
	Id            int     `json:"id"`
	UserId        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	CreateTime    string  `json:"create_time"`
	FromUserId    int     `json:"from_user_id"`
	Comment       string  `json:"comment"`
	Ticket        int     `json:"ticket"`
	Status        int     `json:"status"`
}

func (a *Score) Insert() {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *Score) Update(col string, where string) {
	global.C.DB.NewUpdate().Model(a).Column(col).Where(where).Exec(global.C.Ctx)
}

func (a *Score) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Score) Gets(where string) ([]Score, int) {
	var datas []Score
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Score) Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *Score) Page(where string, page, limit int) ([]Score, int) {
	var datas []Score
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Score) Balance(where string) float64 {
	var equity float64
	global.C.DB.NewSelect().Model(a).Column("sum(amount) as amount").Where(where).Scan(global.C.Ctx, &equity)
	return equity
}

func (a *Score) Del(where string) {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}
