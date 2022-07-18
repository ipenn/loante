package model

import (
	"loante/global"
	"fmt"
	"github.com/uptrace/bun"
)

type Coupon struct {
	bun.BaseModel `bun:"table:coupons,alias:u"`
	Id	int	`json:"id"`
	UserId	int	`json:"user_id"`
	CouponNo	string	`json:"coupon_no"`
	CreateTime	string	`json:"create_time"`
	EndTime	string	`json:"end_time"`
	FromUserId	int	`json:"from_user_id"`
	Type	string	`json:"type"`
	Status	int	`json:"status"`
	UseTime	string	`json:"use_time"`
}

func (a *Coupon)Insert()  {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *Coupon)Update(where string)  {
	global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
}

func (a *Coupon)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Coupon)Gets(where string) ([]Coupon, int) {
	var datas []Coupon
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Coupon)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *Coupon)Del(where string)  {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}


